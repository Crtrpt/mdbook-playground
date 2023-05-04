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
var ImageList map[string]*types.ImageSummary

func main() {
	ImageList = make(map[string]*types.ImageSummary)
	Cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	images, err := Cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	for _, image := range images {
		ImageList[image.RepoTags[0]] = &image
		fmt.Printf("name:%v", image)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		buf, _ := io.ReadAll(r.Body)
		req, _ := url.ParseQuery(string(buf))

		if ImageList[req.Get("image")] == nil {
			w.WriteHeader(500)
			return
		}
		fmt.Printf("start image %+v \r\n", req.Get("image"))
		err := Cli.ContainerStart(context.Background(), ImageList[req.Get("image")].ID, types.ContainerStartOptions{})
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		w.WriteHeader(200)
		w.Write(buf)
	})
	log.Fatal(http.ListenAndServe(":9080", nil))
}
