package internal

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/rs/xid"
)

func StartC1(ctx context.Context, image, code string) (res string, err error) {
	cid := xid.New().String()
	cli := ctx.Value("client").(*client.Client)
	fmt.Printf("start image:%+v cid:%s \r\n ", image, cid)
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Cmd:   strings.Split(code, " "),
		Image: image},
		nil, nil, nil, cid)
	if err != nil {
		return
	}
	cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return
	}

	rx, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Details:    true,
	})
	if err != nil {
		panic(err.Error())

	}
	defer rx.Close()

	buf := make([]byte, 1024)
	for {
		len, err := rx.Read(buf[:])
		if err != nil {
			fmt.Print("err   ")
			break
		}
		if len == 0 {
			break
		}
		resp := ctx.Value("resp").(http.ResponseWriter)
		if resp != nil {
			resp.Write(buf[0:len])
		}
	}
	if err != nil {
		return
	}
	return
}
