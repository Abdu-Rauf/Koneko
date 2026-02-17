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

	if pc!=nil{
		pc.Close()
	}

	// save session info 
	if container!=nil{
		log.Println("Save session history")
	}
	if container != nil {
		if container.Conn != nil  {
			log.Println("Closing the videoStream connections")
			container.Conn.Close()
		}
		if container.InputConn !=nil{
			log.Println("Closing the Input Connection")
			container.InputConn.Close()
		}
		RemoveContainer(context.Background(),dockerCli,container.ID)
		log.Println("Completed Cleanup")
	}


}