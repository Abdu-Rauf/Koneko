package main

import (
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

type Server struct {
	DockerCli *client.Client
}

func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {

	var closeOnce sync.Once

	// Connect for signalling
	conn,err:= SignalSocket(w,r)
	if err!=nil{
		log.Println("Error in establishing ws connection")
		return
	}
	log.Println("Websocket connected")

	// Signal Channels for closing ws and peer
	connectionClosed := make(chan struct{})
	dataChannelReady := make(chan struct{})
	// containerReady := make(chan struct{})

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


	// Add video Track before creating Peer
	videoTrack, err := webrtc.NewTrackLocalStaticSample(
		webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeH264},
		"video",
		"koneko-stream",
	)
	if err != nil {
		log.Println("Failed to create video track:", err)
		return
	}

	// Start Container Setup In a Routine
	go func(){
		c,cs,err := ContainerSetup(videoTrack,browserMsg.Browser,s.DockerCli)
		if err!=nil{
			log.Println("Error Setting Up Container")
			conn.Close()
			return
		}
		container = c
		cancelStream = cs
	}()


	// Setup Peer Before Signalling and Create DataChannel
	peer, err := PeerSetup(conn,videoTrack,connectionClosed,&closeOnce)
	if err!=nil{
		log.Println("Error setting Up peers")
		conn.Close()
		return
	}
	
	// Attach DataChannel Listeners For User Input Forwarding
	AttachDcListeners(
		peer.DC,
		&container,
		dataChannelReady,
		connectionClosed,
		&closeOnce,
	)
	

	// Read messages from websocket
	go Signalling(conn,peer.PC)

	// Wait for dc to be ready then close the ws
	<-dataChannelReady
	log.Println("Closing Websocket, data channel is ready")
	time.Sleep(200 * time.Millisecond)
	conn.Close()

	// Wait till connection is closed by client
	<-connectionClosed
	log.Println("Connection closed , clean up")

	// Cleanup
	Cleanup(cancelStream,container,s.DockerCli,peer.PC)

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
