package main

import (
	"log"
)

func ForwardUserInputs(data *DataChannelInputs, dataChannelReady chan struct{}, container *Container){

	// Launch the container based on selected browser image
	switch data.Type {
	case "dc_ready":
		log.Println("Data Channel established")
		close(dataChannelReady)
	case "mouse_move":
		if container!=nil{
			err := container.MouseMove(data.X,data.Y)
			if err!=nil{
				log.Println("Error moving mouse",err)
				return 
			}

		}
	case "key_press":
		if container!=nil{
			err := container.KeyPress(data.Key)
			if err!=nil{
				log.Println("Error moving mouse",err)
				return 
			}

		}
		
	case "mouse_click":
		if container!=nil{
			err := container.MouseClick(data.Click)
			if err!=nil{
				log.Println("Error moving mouse",err)
				return 
			}

		}
	default:
		log.Println("Unkown message type",data.Type)
		return

	}
}