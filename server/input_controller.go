package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "os/exec"
    "strconv"
    "time"
    "github.com/asticode/go-astits"
    "github.com/pion/webrtc/v3"
    "github.com/pion/webrtc/v3/pkg/media"   
    "context"
)

type Container struct {
	ID string
	Image string
	Address string
	Conn net.Conn
}


func (c *Container) Connect() error {
    log.Println("Attempting to connect to container:", c.Address)
    
    maxRetries := 20
    for i := 0; i < maxRetries; i++ {
        conn, err := net.DialTimeout("tcp", c.Address, 1*time.Second)
        if err == nil {
            c.Conn = conn
            log.Println("Connected to container:", c.ID)
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

func (c *Container ) MouseMove(x,y int) error{
	cmd := exec.Command("docker","exec",c.ID,"xdotool","mousemove","--sync",strconv.Itoa(x),strconv.Itoa(y))
	return cmd.Run()
}

func (c *Container ) KeyPress(key string) error{
	cmd := exec.Command("docker","exec",c.ID,"xdotool","key",key)
	return cmd.Run()
}
func (c *Container) MouseClick(button int) error {
    var buttstr string
    
    switch button {
    case 0:
        buttstr = "1"  // Left click
    case 1:
        buttstr = "2"  // Middle click
    case 2:
        buttstr = "3"  // Right click
    default:
        log.Printf("Unknown mouse button: %d, defaulting to left click", button)
        buttstr = "1"
    }
    
    cmd := exec.Command("docker", "exec", c.ID, "xdotool", "click", buttstr)
    return cmd.Run()
}
