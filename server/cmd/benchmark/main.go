package main

import (
	"log"	
	"github.com/gorilla/websocket"
	"time"
	// "github.com/docker/docker/client"
)


func main(){

    conn,_,err:= websocket.DefaultDialer.Dial("ws://localhost:8080/ws",nil)
    if err!=nil{
        log.Println("Error connecting to Koneko")
    }
    
    conn.WriteJSON(map[string]string{
        "type":"browser_select",
        "browser":"chrome",
    })

    pc,err := PeerSetup(conn)
    if err!=nil{
        log.Println("Error Creating Peer Connection",err)
    }
    go Signalling(conn,pc)
    dc := <-pc.DCReady
    log.Println("DataChannel Ready")
    _ = dc

    log.Println("Closing data channel...")
    dc.Close()
    time.Sleep(500 * time.Millisecond) 

    log.Println("Closing peer connection...")
    pc.PeerConnection.Close()
    time.Sleep(500 * time.Millisecond)

    // go func (){

    //     _,msg,err:= conn.ReadMessage()

    //     var response struct {
    //         Type   string `json:"type"`
    //         Action  string `json:"action"`
    //         LatencyUs float64 `json:"latency_us"`
    //     }

    //     json.Unmarshal(msg,&response)


        
    // }



}