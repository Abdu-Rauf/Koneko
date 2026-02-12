package main

import (
	"github.com/Abdu-Rauf/koneko/utils"
	"encoding/json"
	"log"
	"net/http"
	"sync"	
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
	"context"
	"time"
	"github.com/docker/docker/client"
)


func CheckDomain(r *http.Request) bool {
	// allow only koneko's domain
	return true
}

var upgrader = websocket.Upgrader{
	CheckOrigin: CheckDomain,
}


type DataChannelInputs struct {
	Type	string	`json:"type"`
	X	int 	`json:"x,omitempty"`
	Y	int 	`json:"y,omitempty"`
	Key	string 	`json:"key,omitempty"`
	Click	int `json:"click,omitempty"`
}

type Server struct {
	DockerCli *client.Client
}

func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {

	var closeOnce sync.Once
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		http.Error(w, "Failed to upgrade to websocket", http.StatusInternalServerError)
		return
	} else {
		log.Println("WebSocket connection established")
	}

	// Signal Channels for closing ws and peer
	connectionClosed := make(chan struct{})
	dataChannelReady := make(chan struct{})

	// Necessary for Cleanup
	var container *Container
	var cancelStream context.CancelFunc

	var browserMsg struct{
		Type string `json:"type"`
		Browser string `json:"browser"`
	}
	err = conn.ReadJSON(&browserMsg)
	if err!=nil{
		log.Println("Error reading browser selected")
		return
	}
	log.Println("Browser selected:" ,browserMsg.Browser)

	videoTrack, err := webrtc.NewTrackLocalStaticSample(
		webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeH264},
		"video",
		"koneko-stream",
	)
	if err != nil {
		log.Println("Failed to create video track:", err)
		return
	}

	go func(){

		log.Println("Starting Container")
		containerID,image,streamURL, err := utils.StartContainer(r.Context(),s.DockerCli,browserMsg.Browser)

		if err != nil {
			log.Println("Failed to start container:", err)
			return
		}
		container = &Container{
			ID : containerID,
			Address:streamURL,
			Image: image,
		}
		// Connect server to the Conatiner 
		err = container.Connect()
		if err!=nil{
			log.Println("Error connecting to the container",err)
			return
		}
		log.Println("Container is Ready",container)

		// Create context (holy)
		ctx, cancel := context.WithCancel(context.Background())
		cancelStream = cancel

		// start sending video streams to the client
		go func(){
			log.Println("Starting stream")
			err := container.StreamVid(ctx,videoTrack)
			if err!=nil{
				log.Println("Video Streaming stopped",err)
			}
		}()
	}()


	// Handle WebRTC signaling here
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		log.Println("Error creating PeerConnection:", err)
		return
	}
	
	rtpSender, err := peerConnection.AddTrack(videoTrack) 
	if err != nil {
		log.Println("Failed to add track:", err)
		return
	}
	
	// is this necessary to add??
	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
				log.Println("RTCP reader stopped:", rtcpErr)
				return
			}
		}
	}()


	dc ,err := peerConnection.CreateDataChannel("inputs", nil)
	if err!=nil{
		log.Println("Error craeting DataChannel",err)
		return
	}
	log.Println("Server created data channel")

	dc.OnOpen(func() {
		log.Println("Data channel opened through server")
	
	})
	
	dc.OnMessage(func(msg webrtc.DataChannelMessage) {

		select{
		case <- connectionClosed:
			log.Println("Container no Longer Connected")
			return
		default:
			
		}
		
		// Define structure for incoming data
		var data DataChannelInputs
		
		// Unmarshal JSON string into struct
		err := json.Unmarshal(msg.Data, &data)
		if err != nil {
			log.Println("Error parsing message:", err)
			return
		}
		
		// Launch the container based on selected browser image
		switch data.Type {
		case "dc_ready":
			log.Println("Data Channel established")
			close(dataChannelReady)
		case "mouse_move":
			if container!=nil{
				go func(){
					start := time.Now()
					err = container.MouseMove(data.X,data.Y)
					if err!=nil{
						log.Println("Error moving mouse")
					}
					conn.WriteJSON(map[string]interface{}{
						"type":       "benchmark_rep",
						"action":     "mouse_move",
						"latency_us": time.Since(start).Microseconds(),
					})
				}()
			}
		case "key_press":
			if container!=nil{
				go func(){
					start:=time.Now()
					err = container.KeyPress(data.Key)
					if err!=nil{
						log.Println("Error moving mouse")
					}
					conn.WriteJSON(map[string]interface{}{
						"type":       "benchmark_rep",
						"action":     "key_press",
						"latency_us": time.Since(start).Microseconds(),
					})
				}()
			}
			
		case "mouse_click":
			if container!=nil{
				go func(){
					start:=time.Now()
					err = container.MouseClick(data.Click)
					if err!=nil{
						log.Println("Error moving mouse")
					}
					conn.WriteJSON(map[string]interface{}{
						"type":       "benchmark-rep",
						"action":     "mouse_click",
						"latency_us": time.Since(start).Microseconds(),
					})

				}()
			}
		default:
			log.Println("Unkown message type",data.Type)


		}
	})
	
	dc.OnClose(func() {
		log.Println("Data channel closed")
		closeOnce.Do(func() {
			close(connectionClosed) 
		})
	})
	
	dc.OnError(func(err error) {
		log.Println("Data channel error:", err)
	})

	offer , err := peerConnection.CreateOffer(nil)
	if err!=nil{
		log.Println("Error creating offer",err)
		return
	}
	err = peerConnection.SetLocalDescription(offer)
	if err!=nil{
		log.Println("Failed to set local description")
		return
	}
	conn.WriteJSON(utils.SignalInfo{
		Type:"offer",
		Sdp: offer.SDP,
	})
    

	// Send ICE candidates to client
	peerConnection.OnICECandidate(func(cd *webrtc.ICECandidate) {
		if cd != nil {
			log.Println("Sending ICE candidate to client:", cd.ToJSON())
			conn.WriteJSON(utils.SignalInfo{
				Type:      "ice-candidate",
				Candidate: cd.ToJSON(), 
			})
		}
	})

	// Set up watcher
	peerConnection.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		log.Printf("State: %s\n", state.String())
		// If connection dies, send signal
		if state == webrtc.PeerConnectionStateFailed || 
			state == webrtc.PeerConnectionStateClosed ||
			state == webrtc.PeerConnectionStateDisconnected {

				closeOnce.Do(func(){
					close(connectionClosed) // Send dataChannel closed signal
				})  
		}
	})

	// Read messages from websocket
	go utils.Signalling(conn,peerConnection)

	// Wait for dc to be ready then close the ws
	<-dataChannelReady
	log.Println("Closing Websocket, data channel is ready")
	time.Sleep(200 * time.Millisecond)
	conn.Close()
	// check if container is ready , if it is ready then start streaming 

	// Wait till connection is closed by client
	<-connectionClosed
	log.Println("Connection closed , clean up")

	// Cancel video streaming
	if cancelStream!=nil{
		log.Println("Closing Video Stream")
		cancelStream()
		
	}

	// save session info 
	if container!=nil{
		log.Println("Save session history")
	}
	if container != nil {
		if container.Conn != nil {
			log.Println("Closing the tcp connection")
			container.Conn.Close()
		}
		utils.RemoveContainer(context.Background(),s.DockerCli,container.ID)
	}

	// Close Peerconnection
	peerConnection.Close()

}


func main() {
	cli,err := client.NewClientWithOpts(client.FromEnv,client.WithAPIVersionNegotiation())
	if err!=nil{
		log.Println("Failed To Create Docker Cli")
	}
	defer cli.Close()
	myServer := &Server{
		DockerCli:cli,
	}
	
	fs := http.FileServer(http.Dir("../client"))
    http.Handle("/", fs)

	http.HandleFunc("/ws", myServer.wsHandler)
	http.ListenAndServe(":8080", nil)
}
