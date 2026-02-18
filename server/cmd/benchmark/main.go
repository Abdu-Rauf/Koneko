package main

import (
	"log"	
	"github.com/gorilla/websocket"
	"time"
    "os"
)

var videoStarted = make(chan struct{})
func main(){
    
    // create csv file writer
    f, err := os.Create("agent_socket.csv")
    if err != nil {
        log.Fatal("Could not create CSV file")
    }
    defer f.Close()

    // Write the Header
    f.WriteString("timestamp,seq,x,y,latency_ms,drift_ms\n")

    // connect to server and select browser
    conn,_,err:= websocket.DefaultDialer.Dial("ws://localhost:8080/ws",nil)
    if err!=nil{
        log.Println("Error connecting to Koneko")
    }
    
    conn.WriteJSON(map[string]string{
        "type":"browser_select",
        "browser":"chrome",
    })

    // webrtc setup
    pc,err := PeerSetup(conn,f)
    if err!=nil{
        log.Println("Error Creating Peer Connection",err)
    }
    go Signalling(conn,pc)
    dc := <-pc.DCReady
    <-videoStarted

    // wait for server to start streaming
    log.Println("Soaking for 5 seconds...")
    time.Sleep(5 * time.Second)
    SimulateMouse(dc)

    // Cleanup
    log.Println("Closing data channel...")
    dc.Close()
    time.Sleep(500 * time.Millisecond) 

    log.Println("Closing peer connection...")
    pc.PeerConnection.Close()
    time.Sleep(500 * time.Millisecond)

}
