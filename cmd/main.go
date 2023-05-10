package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/crtrpt/mdbook-playground/internal"
)

type DockerKey string

type ReqForm struct {
	Code  string `json:"code"`  //代码
	Repo  string `json:"repo"`  //仓库
	Path  string `json:"path"`  //路径
	Image string `json:"image"` //使用的镜像
}

// 启动容器
func StartC(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // 设置允许访问所有域
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
	w.Header().Set("Access-Control-Max-Age", "172800")
	w.Header().Set("Access-Control-Allow-Credentials", "false")

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
	// fmt.Printf("//req: %+v", req.Code)
	// fmt.Printf("//req: %+v", req.Image)

	if internal.ImageList[req.Image] == nil {

		fmt.Printf("//%+v 出错了", req.Image)

		w.WriteHeader(500)
		return
	}

	ctx := context.WithValue(context.Background(), "client", internal.Cli)
	ctx = context.WithValue(ctx, "resp", w)
	ctx = context.WithValue(ctx, "cfg", internal.Cfg)
	_, err = internal.StartContainer(ctx, req.Image, req.Code)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	// // //获取输出日志

}

// 主程序
func main() {
	internal.InitConfig()
	internal.InitDockerClient()
	internal.AutoRemoveCloseContainer()
	http.HandleFunc("/", StartC)
	fmt.Printf("listen :9080 \r\n")
	log.Fatal(http.ListenAndServe(":9080", nil))
}
