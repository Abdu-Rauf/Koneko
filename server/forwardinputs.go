package main

import (
    "fmt"
    "log"
    "time"
    "github.com/pion/webrtc/v3" 
)

func ForwardUserInputs(data *DataChannelInputs, dataChannelReady chan struct{}, container *Container,dc *webrtc.DataChannel) {
	// Send Signal to close ws
    if data.Type == "dc_ready" {
        log.Println("Data Channel established and synchronized")
 
        select {
        case <-dataChannelReady:
        default:
            close(dataChannelReady)
        }
        return
    }

	// If setup failed, container is nil even if channel is closed, used defer(containerready)
    if container == nil {
        return 
    }

    // Handle Agent Forwarding
	start := time.Now()
	
    if err := container.ForwardAgent(data); err != nil {
        log.Println("Forwarding Error:", err)
    }
	msg := fmt.Sprintf(`{"type":"benchmark_ack","latency":%d,"seq":%d,"ts":%d}`,
        time.Since(start).Microseconds(),
        data.Seq,
        data.Ts,
        )
    dc.SendText(msg)
}