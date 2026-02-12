package main

import (
	// "context"
	// "encoding/csv"
	// "encoding/json"
	// "flag"
	// "fmt"
	"log"
	// "os"
	// "sync"
	"time"

	// "github.com/docker/docker/api/types"
	// "github.com/docker/docker/client"
	"github.com/gorilla/websocket"
	// "github.com/shirou/gopsutil/v3/cpu"
	// "github.com/shirou/gopsutil/v3/mem"
)

// type SystemStats struct {
// 	mu           sync.RWMutex
// 	ServerCPU    float64
// 	ServerMem    float64
// 	ContainerCPU float64
// 	ContainerMem float64
// }

type BenchAction struct {
    Type  string `json:"type"`
    X     int    `json:"x,omitempty"`
    Y     int    `json:"y,omitempty"`
    Key   string `json:"key,omitempty"`
    Click int    `json:"click,omitempty"`
}

func main(){

    // file,err := os.Create(outputFile)
    // if err!=nil{
    //     log.Println("Error Creating File")
    // }
    // defer file.Close()
	// writer := csv.NewWriter(file)
	// defer writer.Flush()

	// writer.Write([]string{
	// 	"timestamp", "operation", "latency_us",
	// 	"server_cpu_pct", "server_mem_gb",
	// 	"container_cpu_pct", "container_mem_mb",
	// })

    conn,_,err:= websocket.DefaultDialer.Dial("ws://localhost:8080/ws",nil)
    if err!=nil{
        log.Println("Error connecting to Koneko")
    }
    conn.WriteJSON(map[string]string{
        "type":"browser_select",
        "browser":"chrome",
    })
    time.Sleep(10*time.Second)
    defer conn.Close()

    // get container and call monitorstats

    go func (){

        _,msg,err:= conn.ReadMessage()

        var response struct {
            Type   string `json:"type"`
            Action  string `json:"action"`
            LatencyUs float64 `json:"latency_us"`
        }

        json.Unmarshal(msg,&response)

        // lock and store it in record

        
    }
    // Scenario A: Mouse Hover (Fast)
	for i := 0; i < 50; i++ {
		sendAction(conn, "mouse_move")
		time.Sleep(50 * time.Millisecond)
	}

	// Scenario B: Open Tabs (Heavy)
    // Scenario C: Open Containers (Multiple)


	log.Println("Done! Check", *outputFile)

    func sendAction(conn *websocket.Conn, actionType string) {
        conn.WriteJSON(&BenchAction{
            Type: actionType,
            X:    100,
            Y:    100,
        })
    }
}