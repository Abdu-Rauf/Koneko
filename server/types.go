package main

import "github.com/pion/webrtc/v3"

type DataChannelInputs struct {
    Type  string `json:"type"`
    X     int    `json:"x,omitempty"`
    Y     int    `json:"y,omitempty"`
    Key   string `json:"key,omitempty"`
    Click int    `json:"click,omitempty"`
    Seq    int     `json:"seq,omitempty"`
    Ts    int     `json:"ts,omitempty"`
}

type PeerSession struct {
    PC *webrtc.PeerConnection
    DC *webrtc.DataChannel
}
