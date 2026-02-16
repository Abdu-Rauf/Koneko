package main

import (
	"log"
	"time"
	"fmt"
	"github.com/pion/webrtc/v3"
)

func ForwardUserInputs(data *DataChannelInputs, dataChannelReady chan struct{}, container *Container, dc *webrtc.DataChannel){

	// Launch the container based on selected browser image
	switch data.Type {
	case "dc_ready":
		log.Println("Data Channel established")
		close(dataChannelReady)
	case "mouse_move":
		if container!=nil{
			start := time.Now()
			err := container.MouseMove(data.X,data.Y)
			if err!=nil{
				log.Println("Error moving mouse")
			}
			msg := fmt.Sprintf(`{"type":"benchmark_ack","latency":%d,"seq":%d,"ts":%d}`,
				time.Since(start).Microseconds(),
				data.Seq,
				data.Ts,
			)
			dc.SendText(msg)

		}
	case "key_press":
		if container!=nil{
			start := time.Now()
			err := container.MouseMove(data.X,data.Y)
			if err!=nil{
				log.Println("Error moving mouse")
			}
			msg := fmt.Sprintf(`{"type":"benchmark_ack","latency":%d,"seq":%d,"ts":%d}`,
				time.Since(start).Microseconds(),
				data.Seq,
				data.Ts,
			)
			dc.SendText(msg)
		}
		
	case "mouse_click":
		if container!=nil{
			start := time.Now()
			err := container.MouseMove(data.X,data.Y)
			if err!=nil{
				log.Println("Error moving mouse")
			}
			msg := fmt.Sprintf(`{"type":"benchmark_ack","latency":%d,"seq":%d,"ts":%d}`,
				time.Since(start).Microseconds(),
				data.Seq,
				data.Ts,
			)
			dc.SendText(msg)

		}
	default:
		log.Println("Unkown message type",data.Type)
		return

	}
}