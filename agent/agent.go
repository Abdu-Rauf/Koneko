	package main

	import (
		"encoding/json"
		"log"
		"net"
		"os/exec"
		"sync/atomic"
		"strconv"
		"time"
	)

	var keyMap = map[string]string{
		" ":          "space",
		"Enter":      "Return",
		"Backspace":  "BackSpace",
		"ArrowUp":    "Up",
		"ArrowDown":  "Down",
		"ArrowLeft":  "Left",
		"ArrowRight": "Right",
		"Escape":     "Escape",
		"Delete":     "Delete",
		"Tab":        "Tab",
		"Home":       "Home",
		"End":        "End",
		"PageUp":     "Prior",
		"PageDown":   "Next",
		"Shift":      "Shift_L",
		"Control":    "Control_L",
		"Alt":        "Alt_L",
		"Meta":       "Super_L",
	}

	type Input struct {
		Type  string `json:"type"`
		X     int    `json:"x,omitempty"`
		Y     int    `json:"y,omitempty"`
		Key   string `json:"key,omitempty"`
		Click int    `json:"click,omitempty"`
	}

	var execCount uint64

	func main() {

		go func() {
			ticker := time.NewTicker(1 * time.Second)
			for range ticker.C {
				// Safely read the current count and reset it to 0
				count := atomic.SwapUint64(&execCount, 0)
				if count > 0 { // Only clutter the logs if there was actual input
					log.Printf("=> xdotool spawned %d times in the last second", count)
				}
			}
		}()

		// Listen on Port 5050
		listener, err := net.Listen("tcp", ":5050")
		if err != nil {
			log.Fatal("Failed to bind port 5050:", err)
		}
		defer listener.Close()
		log.Println("Koneko Input Agent listening on :5050")

		for {
			// Accept new connection from the main server
			conn, err := listener.Accept()
			if err != nil {
				log.Println("Connection error:", err)
				continue
			}
			
			// Handle connection in a goroutine
			go handleConnection(conn)
		}
	}

	func handleConnection(conn net.Conn) {
		defer conn.Close()
		decoder := json.NewDecoder(conn)

		for {
			var input Input
			if err := decoder.Decode(&input); err != nil {
				log.Println("Connection closed or decode error")
				return
			}

			// Execute xdotool command based on input
			var cmd *exec.Cmd

			switch input.Type {
			case "mouse_move":
				cmd = exec.Command("xdotool", "mousemove", strconv.Itoa(input.X), strconv.Itoa(input.Y))
			case "mouse_click":
				buttonmap:= map[int]string{
					0:"1",
					1:"2",
					2:"3",
				}
				btn,ok := buttonmap[input.Click]
				if !ok{
					btn = "1"
				}
				cmd = exec.Command("xdotool", "click", btn)
			case "key_press":
				mappedKey := input.Key
				if val, ok := keyMap[input.Key]; ok {
					mappedKey = val
				}
				cmd = exec.Command("xdotool", "key", mappedKey)
			// case "scroll": 
			//      // Add scroll logic if needed
			default:
				continue
			}

			if cmd != nil {
				atomic.AddUint64(&execCount, 1)
				if err := cmd.Run(); err != nil {
					log.Println("xdotool error:", err)
				}
			}
		}
	}