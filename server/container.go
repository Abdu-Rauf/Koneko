package main

import (
	"log"
	"github.com/pion/webrtc/v3"
	"context"
	"github.com/docker/docker/client"
	"fmt"
	containerTypes "github.com/docker/docker/api/types/container"
	"time"
)

func StartContainer(ctx context.Context,cli *client.Client, browser string) (string,string,string,error){
	var image string
	
	switch browser{
	case "chrome":
		image = "koneko-chrome:3.0"
	case "firefox":
		image = "koneko-firefox:1.0"
	}
	
	stx,scancel := context.WithTimeout(ctx, 10*time.Second)
	defer scancel()	

	config := &containerTypes.Config{
		Image:image,
	}
	hostconfig := &containerTypes.HostConfig{
		AutoRemove: true,
	}
	resp,err := cli.ContainerCreate(stx,config,hostconfig,nil,nil,"")
	if err!=nil{
		return "","","",fmt.Errorf(`Error Creating Conatainer %w`,err)
	}
	if err = cli.ContainerStart(stx,resp.ID,containerTypes.StartOptions{}); err!=nil{
		return "","","",fmt.Errorf(`Error Starting Conatainer %w`,err)

	}
	containerJSON, err := cli.ContainerInspect(stx,resp.ID)

	if containerJSON.NetworkSettings.IPAddress == "" {
        return "","", "", fmt.Errorf("container started but has no IP address")
    }

	streamURL := fmt.Sprintf(`%s:8080`,containerJSON.NetworkSettings.IPAddress)

	return resp.ID,image,streamURL ,nil

}


func RemoveContainer(ctx context.Context ,cli *client.Client ,containerID string) error {

	rtx,rcancel := context.WithTimeout(ctx,10*time.Second)
	defer rcancel()
	if err := cli.ContainerStop(rtx,containerID,containerTypes.StopOptions{}); err !=nil{
		return fmt.Errorf("Error Stopping the Container %w",err)
	}
	return nil
}


func ContainerSetup(videoTrack *webrtc.TrackLocalStaticSample, browser string, dockerCli *client.Client,) (*Container , context.CancelFunc, error){

	log.Println("Starting Container")
	containerID,image,streamURL, err := StartContainer(context.Background(),dockerCli,browser)

	if err != nil {
		log.Println("Failed to start container:", err)
		return nil,nil,err
	}
	container := &Container{
		ID : containerID,
		Address:streamURL,
		Image: image,
	}
	// Connect server to the Conatiner 
	err = container.Connect()
	if err!=nil{
		log.Println("Error connecting to the container",err)
		return nil,nil,err
	}
	log.Println("Container is Ready",container)

	// Create context (holy)
	ctx, cancel := context.WithCancel(context.Background())

	// start sending video streams to the client
	go StartStreaming(container,ctx,videoTrack)
	
	return container,cancel,nil
}
