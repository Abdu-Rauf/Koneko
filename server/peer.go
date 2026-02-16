package main

import (
	"log"
	"sync"	
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

func PeerSetup(
    conn *websocket.Conn,
    videoTrack *webrtc.TrackLocalStaticSample,
    connectionClosed chan struct{},
    closeOnce *sync.Once,
) (*PeerSession,error){

	// Create New Peer Session
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		log.Println("Error creating PeerConnection:", err)
		return nil,err
	}

	// Add the Videotrack to recieve VidPackets
	rtpSender, err := peerConnection.AddTrack(videoTrack) 
	if err != nil {
		log.Println("Failed to add track:", err)
		return nil,err
	}

	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
				log.Println("RTCP reader stopped:", rtcpErr)
				return 
			}
		}
	}()

	// Send ICE candidates to client
	peerConnection.OnICECandidate(func(cd *webrtc.ICECandidate) {
		if cd != nil {
			log.Println("Sending ICE candidate to client:", cd.ToJSON())
			conn.WriteJSON(SignalInfo{
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

	// Create DataChannel For Recieving 
	dc ,err := peerConnection.CreateDataChannel("inputs", nil)
	if err!=nil{
		log.Println("Error craeting DataChannel",err)
		return nil,err
	}
	log.Println("Server created data channel")

	return &PeerSession{
        PC: peerConnection,
        DC: dc,
    }, nil

}