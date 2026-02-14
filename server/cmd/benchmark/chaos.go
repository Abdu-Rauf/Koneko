package main

import (

	"github.com/pion/webrtc/v3"
	"log"
	"time"
	"fmt"
	"math"
)
func SimulateMouse(dc *webrtc.DataChannel){
    log.Println("starting mouse simulation")

    // send 60 inputs/sec
    ticker := time.NewTicker(16*time.Millisecond)
    timeout := time.After(30*time.Second)

    seq := 0

    // Setup cordinates for circle 
    centerX, centerY := 960.0, 540.0
    radius := 100.0
    angle := 0.0

    for {
        select{
        case<-timeout:
            return
        case<-ticker.C:
            seq++
            angle+=0.1

            // Calculate Circle Coordinates
            x := centerX + radius*math.Cos(angle)
            y := centerY + radius*math.Sin(angle)

            ts := time.Now().UnixMicro()

            // Construct Payload
            msg := fmt.Sprintf(`{"type":"mouse_move","x":%d,"y":%d,"seq":%d,"ts":%d}`, 
                int(x), int(y), seq, ts)
            
            dc.SendText(msg)

        }
    }
}