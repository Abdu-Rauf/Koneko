package main

import (
    "encoding/json"
    "log"

    "github.com/gorilla/websocket"
    "github.com/pion/webrtc/v3"
)

type SignalMessage struct {
    Type      string                  `json:"type"`
    Sdp       string                  `json:"sdp"`
    Candidate webrtc.ICECandidateInit `json:"candidate"`
}

func Signalling(conn *websocket.Conn, pcSession *PeerSession) {
    for {
        _, rawMsg, err := conn.ReadMessage()  
        if err != nil {
            log.Println("Error reading message from ws:", err)
            return  
        }

        var msg SignalMessage  
        json.Unmarshal(rawMsg, &msg)  

        switch msg.Type {
        case "offer":
            log.Println("Received offer from server")
            pcSession.PeerConnection.SetRemoteDescription(webrtc.SessionDescription{
                Type: webrtc.SDPTypeOffer,
                SDP:  msg.Sdp,
            })
            answer, _ := pcSession.PeerConnection.CreateAnswer(nil)
            pcSession.PeerConnection.SetLocalDescription(answer)
			pcSession.WsMu.Lock()
            conn.WriteJSON(SignalMessage{  
                Type: "answer",
                Sdp:  answer.SDP,  
            })
			pcSession.WsMu.Unlock()

        case "ice-candidate":
            pcSession.PeerConnection.AddICECandidate(msg.Candidate)
        }
    }
}