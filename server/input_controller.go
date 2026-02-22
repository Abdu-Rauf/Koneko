package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "time"
    "github.com/asticode/go-astits"
    "github.com/pion/webrtc/v3"
    "github.com/pion/webrtc/v3/pkg/media"   
    "context"
    "encoding/json"
)

type Container struct {
	ID string
	Image string
	Address string
    AgentAddress string
	Conn net.Conn
    InputConn net.Conn
    agentEncoder *json.Encoder
    MousePredictor *Predictor
}


func (c *Container) Connect() error {
    log.Println("Attempting to connect to Vidstream:", c.Address)
    
    maxRetries := 20
    for i := 0; i < maxRetries; i++ {
        conn, err := net.DialTimeout("tcp", c.Address, 1*time.Second)
        if err == nil {
            c.Conn = conn
            log.Println("Connected to vidstream:", c.ID)
            return nil 
        }
        
        log.Printf("Connection attempt %d/%d failed, retrying...", i+1, maxRetries)
        time.Sleep(500 * time.Millisecond)
    }
    
    return fmt.Errorf("failed to connect after %d attempts", maxRetries)
}

func (c *Container) AgentConnect() error{

    log.Println("Attempting to connect to Agent:", c.AgentAddress)
    
    maxRetries := 20
    for i := 0; i < maxRetries; i++ {
        conn, err := net.DialTimeout("tcp", c.AgentAddress, 1*time.Second)
        if err == nil {
            // save the connection and encoder
            c.InputConn = conn
            c.agentEncoder = json.NewEncoder(conn)
            log.Println("Connected to Agent:", c.ID)
            return nil 
        }
        
        log.Printf("Connection attempt %d/%d failed, retrying...", i+1, maxRetries)
        time.Sleep(500 * time.Millisecond)
    }
    
    return fmt.Errorf("failed to connect after %d attempts", maxRetries)
}

func (c *Container) StreamVid(ctx context.Context, videoTrack *webrtc.TrackLocalStaticSample) error{
    bufferedReader := bufio.NewReaderSize(c.Conn, 188*1024)
    demuxer := astits.NewDemuxer(ctx, bufferedReader)
    
    packetCount := 0
    for {
        select {
        case <-ctx.Done():
            log.Println("Context cancelled")
            return ctx.Err()
        default:
        }
        
        data, err := demuxer.NextData()
        if err != nil {
            return fmt.Errorf("Error reading nal units: %w", err)
        }
        if data.PES == nil || data.PES.Data == nil {
            continue
        }
        packetCount++
        if packetCount%100 == 0{
            log.Println("Number of packets sent",packetCount)
        }

        err = videoTrack.WriteSample(media.Sample{
            Data:     data.PES.Data,
            Duration: time.Millisecond * 33,
        })
        if err != nil {
            return fmt.Errorf("Failed to write media sample: %w", err)
        }
    }
}

func (c *Container) ForwardAgent(data *DataChannelInputs) error {

    if c.agentEncoder == nil {
        return fmt.Errorf("agent encoder not found, is the agent connected?")
    }
    if data.Type == "mouse_move" && c.MousePredictor != nil {
        px, py := c.MousePredictor.Predict(float64(data.X), float64(data.Y))
        log.Println("using predictor")
        // Update the struct before sending
        data.X = int(px)
        data.Y = int(py)
    }
    // converts the struct to JSON and writes it to the conn
    return c.agentEncoder.Encode(data)
}
