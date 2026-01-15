package main

import (
	// "encoding/json"
	// "log"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
	// "github.com/pion/webrtc/v3"
)

func CheckDomain(r *http.Request) bool {
	// allow only koneko's domain
	return true
}

var upgrader = websocket.Upgrader{
	CheckOrigin: CheckDomain,
}

type Message struct {
	Type      string                  `json:"type"`
	Sdp       string                  `json:"sdp"`
	Candidate webrtc.ICECandidateInit `json:"candidate"` // why pointer causes issues?
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		http.Error(w, "Failed to upgrade to websocket", http.StatusInternalServerError)
	} else {
		log.Println("WebSocket connection established")
	}
	// Handle WebRTC signaling here

	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		log.Println("Error creating PeerConnection:", err)
		return
	}

	// Channel to signal when peer connection closes
	connectionClosed := make(chan struct{})

	// Data Channel

    peerConnection.OnDataChannel(func(dc *webrtc.DataChannel) {
        log.Println("Data channel opened:", dc.Label())
        
        dc.OnOpen(func() {
            log.Println("Data channel ready to use")
        })
        
        dc.OnMessage(func(msg webrtc.DataChannelMessage) {
            log.Println("Received message:", string(msg.Data))
        })
        
        dc.OnClose(func() {
            log.Println("Data channel closed")
        })
        
        dc.OnError(func(err error) {
            log.Println("Data channel error:", err)
        })
    })
	// Send ICE candidates to client

	peerConnection.OnICECandidate(func(cd *webrtc.ICECandidate) {
		if cd != nil {
			log.Println("Sending ICE candidate to client:", cd.ToJSON())
			conn.WriteJSON(Message{
				Type:      "ice-candidate",
				Candidate: cd.ToJSON(), /// check the format of cd before and afete ToJSON
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

				close(connectionClosed)  // ← Send dataChannel closed signal
		}
	})

	// Read messages from websocket
	go func(){
		defer conn.Close()
		defer log.Println("closing websocket")

		for {
			var msg Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Println("Read error:", err)
				break
			}
			log.Println("Received message:", msg.Type)
			if msg.Type == "offer" {
				log.Println(msg.Sdp)
				// create an offer with the session description info
				offer := webrtc.SessionDescription{
					Type: webrtc.SDPTypeOffer,
					SDP:  msg.Sdp,
				}

				// set the remote description according to the recieved offer
				err = peerConnection.SetRemoteDescription(offer)
				if err != nil {
					log.Println("Error setting remote description:", err)
					continue
				}
				// create an answer according to the remote description
				answer, err := peerConnection.CreateAnswer(nil)
				if err != nil {
					log.Println("Error creating answer")
					continue
				}
				// set the local description with the created answer
				err = peerConnection.SetLocalDescription(answer)
				if err != nil {
					log.Println("Error setting the local description", err)
					continue
				}
				// send the answer to the js client
				log.Println(`sending answer to the client`,answer)
				err = conn.WriteJSON(Message{Type: "answer", Sdp: answer.SDP})
				if err != nil {
					log.Println("Error sending msg to client", err)
					continue
				}
			}
			if msg.Type == "ice-candidate" {
				log.Println(msg.Candidate)
				err = peerConnection.AddICECandidate(msg.Candidate)
				if err != nil {
					log.Println("Error adding the ice candidates")
				}
			}
		}
	}()
	<-connectionClosed

	log.Println(`cleaning up`)
	peerConnection.Close()

}


func main() {
	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe(":8080", nil)
}
