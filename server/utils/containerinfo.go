package utils

import (
    "fmt"
	"context"
	containerTypes "github.com/docker/docker/api/types/container"
	"time"
	"github.com/docker/docker/client"
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
