package main

import (
    "log"
    "github.com/gorilla/websocket"
    "github.com/pion/webrtc/v3"
	"sync"
    "time"
    "os"
    "fmt"
    "github.com/buger/jsonparser"
    
)
var once sync.Once

type PeerSession struct {
    PeerConnection *webrtc.PeerConnection
    DCReady        chan *webrtc.DataChannel
	WsMu           sync.Mutex
}

func PeerSetup(conn *websocket.Conn ,csvFile *os.File) (*PeerSession, error) {
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

            msgType, err := jsonparser.GetString(msg.Data, "type")
            if err != nil || msgType != "benchmark_ack" {
                return
            }
            // 2. Grab the numbers directly from the byte slice
            // No Unmarshal, no reflection, no allocations.
            seq, _ := jsonparser.GetInt(msg.Data, "seq")
            ts, _ := jsonparser.GetInt(msg.Data, "ts")
            latency, _ := jsonparser.GetInt(msg.Data, "latency")
            x, _ := jsonparser.GetInt(msg.Data, "x")
            y, _ := jsonparser.GetInt(msg.Data, "y")

            now := time.Now().UnixMicro()
            totalDrift := now - ts

            // Append row in the file
            fmt.Fprintf(csvFile, "%f,%d,%d,%d,%.2f,%.2f\n",
            float64(now)/1000000.0,
            seq,
            x,
            y,
            float64(latency)/1000.0,
            float64(totalDrift)/1000.0,
            )        
        })
        dc.OnClose(func() {
            log.Println("Data channel closed")
        })
    })

    // Consume video packets to keep connection alive
    pc.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
        
        buf := make([]byte, 1500)
        _, _, err := track.Read(buf)
        if err != nil {
            return
        }
        once.Do(func() {
            log.Println("First packet received..")
            close(videoStarted)
        })
        for {
            if _, _, err := track.Read(buf); err != nil {
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
        return
    })

    return session, nil
}
