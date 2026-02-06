package utils

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
		err := conn.ReadJSON(&msg)
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
			// set the remote description according to the recieved offer
			// err = peerConnection.SetRemoteDescription(offer)
			// if err != nil {
			// 	log.Println("Error setting remote description:", err)
			// 	continue
			// }
			// // create an answer according to the remote description
			// answer, err := peerConnection.CreateAnswer(nil)
			// if err != nil {
			// 	log.Println("Error creating answer")
			// 	continue
			// }
			// // set the local description with the created answer
			// err = peerConnection.SetLocalDescription(answer)
			// if err != nil {
			// 	log.Println("Error setting the local description", err)
			// 	continue
			// }
			// // send the answer to the js client
			// log.Println(`sending answer to the client`,answer)
			// err = conn.WriteJSON(SignalInfo{Type: "answer", Sdp: answer.SDP})
			// if err != nil {
			// 	log.Println("Error sending msg to client", err)
			// 	continue
			// }
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