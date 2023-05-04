package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var Cli *client.Client
var ContainerList map[string]*types.Container

func main() {
	ContainerList = make(map[string]*types.Container)
	Cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	containers, err := Cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		ContainerList[container.Names[0]] = &container
		fmt.Printf("name:%s id:%s\r\n", container.Names[0], container.ID)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		buf, _ := io.ReadAll(r.Body)
		req, _ := url.ParseQuery(string(buf))

		if ContainerList[req.Get("image")] == nil {
			w.WriteHeader(500)
			return
		}
		fmt.Printf("start image %+v \r\n", req.Get("image"))
		err := Cli.ContainerStart(context.Background(), ContainerList[req.Get("image")].ID, types.ContainerStartOptions{})
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		w.WriteHeader(200)
		w.Write(buf)
	})
	log.Fatal(http.ListenAndServe(":9080", nil))
}
