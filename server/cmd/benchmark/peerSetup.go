package main

import (
    "log"

    "github.com/gorilla/websocket"
    "github.com/pion/webrtc/v3"
	"sync"
)

type PeerSession struct {
    PeerConnection *webrtc.PeerConnection
    DCReady        chan *webrtc.DataChannel
	WsMu           sync.Mutex
}

func PeerSetup(conn *websocket.Conn) (*PeerSession, error) {
    pc, err := webrtc.NewPeerConnection(webrtc.Configuration{})
    if err != nil {
        return nil, err
    }

    session := &PeerSession{
        PeerConnection: pc,
        DCReady:        make(chan *webrtc.DataChannel, 1),
    }

    // Handle data channel from server
    pc.OnDataChannel(func(dc *webrtc.DataChannel) {
        dc.OnOpen(func() {
            log.Println("Data channel opened")
            dc.Send([]byte(`{"type":"dc_ready"}`))
            session.DCReady <- dc
        })

        dc.OnMessage(func(msg webrtc.DataChannelMessage) {
            log.Println("Received:", string(msg.Data))
        })
    })

    // Consume video packets to keep connection alive
    pc.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
        log.Println("Receiving video track:", track.Codec().MimeType)
        buf := make([]byte, 1500)
        for {
            _, _, err := track.Read(buf)
            if err != nil {
                return
            }
        }
    })

    // Send ICE candidates to server
    pc.OnICECandidate(func(c *webrtc.ICECandidate) {
        if c != nil {
			session.WsMu.Lock()
            conn.WriteJSON(map[string]interface{}{
                "type":      "ice-candidate",
                "candidate": c.ToJSON(),
            })
			session.WsMu.Unlock()
        }
    })

    pc.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
        log.Println("Peer connection state:", state.String())
    })

    return session, nil
}
