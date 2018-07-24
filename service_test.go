package kbimgsvr

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"./api"
	"./client"
	"./server"
)

func TestImageService(t *testing.T) {
	server, err := server.NewServer("127.0.0.1:9999", "./", 1024*500)
	if err != nil {
		panic(err)
	}
	defer server.Stop()
	c, err := client.NewClient("127.0.0.1:9999")
	if err != nil {
		panic(err)
	}
	defer c.Close()
	f, err := os.Open("dalian_home7.jpeg")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	imgData := bytes.NewBuffer(make([]byte, 0, 102400))
	encoder := base64.NewEncoder(base64.StdEncoding, imgData)
	_, err = encoder.Write(b)
	if err != nil {
		panic(err)
	}
	encoder.Close()
	putResp, err := c.PutImage(context.Background(), &api.PutRequest{Data: imgData.Bytes(), Width: 1920, Height: 1080, Class: "nb"})
	if err != nil {
		panic(err)
	}
	fmt.Println(putResp.Name)
	getResp, err := c.GetImage(context.Background(), &api.GetRequest{Name: putResp.Name})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(getResp.Data))
}
