package main

import (
	"log"
)

func ForwardUserInputs(data *DataChannelInputs, dataChannelReady chan struct{}, container *Container) {
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

    if err := container.ForwardAgent(data); err != nil {
        log.Println("Forwarding Error:", err)
    }
}