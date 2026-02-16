package main

import (
    "encoding/json"
    "log"
    "sync"
    "github.com/pion/webrtc/v3"
)


func AttachDcListeners(
    dc *webrtc.DataChannel,
    container **Container,
    dataChannelReady chan struct{},
    connectionClosed chan struct{},
    closeOnce *sync.Once,
) {

    dc.OnOpen(func() {
        log.Println("Data channel opened through server")
    })

    dc.OnMessage(func(msg webrtc.DataChannelMessage) {
        select {
        case <-connectionClosed:
            log.Println("Connection already closed")
            return
        default:
        }

        var data DataChannelInputs
        if err := json.Unmarshal(msg.Data, &data); err != nil {
            log.Println("Error parsing message:", err)
            return
        }

        go ForwardUserInputs(&data, dataChannelReady, *container, dc)
    })

    dc.OnClose(func() {
        log.Println("Data channel closed")
        closeOnce.Do(func() {
            close(connectionClosed)
        })
    })

    dc.OnError(func(err error) {
        log.Println("Data channel error:", err)
    })
}

