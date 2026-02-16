package main

import (
    "log"
    
    "github.com/gorilla/websocket"
    "github.com/pion/webrtc/v3"
)

type SignalInfo struct {
	Type      string                  `json:"type"`
	Sdp       string                  `json:"sdp"`
	Candidate webrtc.ICECandidateInit `json:"candidate"` // why pointer causes issues?
}

func Signalling(conn *websocket.Conn,peerConnection *webrtc.PeerConnection){
	for {

		var msg SignalInfo
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
		conn.WriteJSON(SignalInfo{
			Type:"offer",
			Sdp: offer.SDP,
		})

		err = conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Println("Received message:", msg.Type)
		if msg.Type == "answer" {
			log.Println(msg.Sdp)
			// create an offer with the session description info
			answer := webrtc.SessionDescription{
				Type: webrtc.SDPTypeAnswer,
				SDP:  msg.Sdp,
	
			}
			err = peerConnection.SetRemoteDescription(answer)
			if err!=nil{
				log.Println("Error setting remote desc",err)
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
}