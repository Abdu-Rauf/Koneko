package main

import(
	"context"
	"log"
	"github.com/pion/webrtc/v3"
	"github.com/docker/docker/client"

)


func Cleanup(
    cancelStream context.CancelFunc,
    container *Container,
    dockerCli *client.Client,
    pc *webrtc.PeerConnection,
) {
	// Cancel video streaming
	if cancelStream!=nil{
		log.Println("Closing Video Stream")
		cancelStream()
	}

	// save session info 
	if container!=nil{
		log.Println("Save session history")
	}
	if container != nil {
		if container.Conn != nil {
			log.Println("Closing the tcp connection")
			container.Conn.Close()
		}
		RemoveContainer(context.Background(),dockerCli,container.ID)
	}

	// Close Peerconnection
	pc.Close()

}