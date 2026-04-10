package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/asticode/go-astits"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
)

type Container struct {
	ID             string
	Image          string
	Address        string
	AgentAddress   string
	Conn           net.Conn
	InputConn      net.Conn
	agentEncoder   *json.Encoder
	MousePredictor *Predictor
	logFile        *os.File
	csvWriter      *bufio.Writer
}

// Method For logging Cursor Inputs
func (c *Container) StartLogging(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}
	c.logFile = file
	c.csvWriter = bufio.NewWriter(file)

	// Write the CSV header immediately
	_, err = c.csvWriter.WriteString("Timestamp,X,Y\n")
	if err != nil {
		return fmt.Errorf("failed to write csv header: %w", err)
	}

	log.Println("CSV Data collection started:", filename)
	return nil
}

// Method For Connecting to The VideoStream
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

// Method For Connecting to the Agent
func (c *Container) AgentConnect() error {

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

// Method To Forward VideoStreams to the client
func (c *Container) StreamVid(ctx context.Context, videoTrack *webrtc.TrackLocalStaticSample) error {
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
		if packetCount%100 == 0 {
			log.Println("Number of packets sent", packetCount)
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

// Method To Forward UserInputs To The Agent
func (c *Container) ForwardAgent(data *DataChannelInputs) error {

	if c.agentEncoder == nil {
		return fmt.Errorf("agent encoder not found, is the agent connected?")
	}
	if data.Type == "mouse_move" {
		// If the logger is initialized, write the data to the buffer
		if c.csvWriter != nil {
			timestamp := time.Now().UnixMilli()
			// Using Fprintf to format the string cleanly without log prefixes
			fmt.Fprintf(c.csvWriter, "%d,%d,%d\n", timestamp, data.X, data.Y)
		}
	}
	// convert to json bytes & write directly to the tcp stream
	return c.agentEncoder.Encode(data)
}
