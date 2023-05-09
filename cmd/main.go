package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/crtrpt/mdbook-playground/internal"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var Cli *client.Client
var ImageList map[string]*types.ImageSummary

type DockerKey string

type ReqForm struct {
	Code  string `json:"code"`  //代码
	Repo  string `json:"repo"`  //仓库
	Path  string `json:"path"`  //路径
	Image string `json:"image"` //使用的镜像
}

// 启动容器
func StartC(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*") // 设置允许访问所有域
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
		w.Header().Set("Access-Control-Max-Age", "172800")
		w.Header().Set("Access-Control-Allow-Credentials", "false")
		w.Header().Set("content-type", "application/json") //// 设置返回格式是json
		return
	}

	buf, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	if len(buf) == 0 {
		w.Write([]byte("empty body"))
		return
	}
	fmt.Printf("%+v", string(buf))
	req := ReqForm{}
	json.Unmarshal(buf, &req)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Printf("//req:%+v", req.Code)
	fmt.Printf("//req:%+v", req.Image)

	if ImageList[req.Image] == nil {

		fmt.Printf("//%+v 出错了", req.Image)

		w.WriteHeader(500)
		return
	}

	ctx := context.WithValue(context.Background(), "client", Cli)
	ctx = context.WithValue(ctx, "resp", w)
	_, err = internal.StartC1(ctx, req.Image, req.Code)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	// //获取输出日志

	w.WriteHeader(200)
	// w.Write(buf)
}
func InitDocker() {
	ImageList = make(map[string]*types.ImageSummary)
	Cli, _ = client.NewClientWithOpts(client.FromEnv)

	images, err := Cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	for _, image := range images {
		ImageList[image.RepoTags[0]] = &image
		fmt.Printf("name:%+v \r\n", image.RepoTags[0])
	}
}

// 主程序
func main() {
	InitDocker()
	http.HandleFunc("/", StartC)
	fmt.Printf("listen :9080 \r\n")
	log.Fatal(http.ListenAndServe(":9080", nil))
}
