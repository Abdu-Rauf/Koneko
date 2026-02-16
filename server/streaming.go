package main

import (
	"log"
	"github.com/pion/webrtc/v3"
	"context"
)

func StartStreaming(container *Container , ctx context.Context, videoTrack *webrtc.TrackLocalStaticSample){
	log.Println("Starting stream")
	err := container.StreamVid(ctx,videoTrack)
	if err!=nil{
		log.Println("Video Streaming stopped",err)
		return
	}	
}