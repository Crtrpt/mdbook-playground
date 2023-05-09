package internal

import (
	"context"
	"fmt"
	"net/http"
	"time"

	mdbookplayground "github.com/crtrpt/mdbook-playground"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/rs/xid"
)

func StartC1(ctx context.Context, image, code string) (res string, err error) {
	cid := xid.New().String()
	cfg := ctx.Value("cfg").(*mdbookplayground.Config)
	cli := ctx.Value("client").(*client.Client)
	v := make(map[string]struct{})

	// 挂载
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Tty: true,
		// Cmd:     strslice.StrSlice{"ls", "/"},
		Cmd:     strslice.StrSlice{"go", "run", "main.go"},
		Volumes: v,
		Image:   image},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:     mount.TypeBind,
					ReadOnly: false,
					Target:   "/var/data",
					Source:   cfg.RepoDir + "/example",
				},
			},
		}, nil, nil, cid)
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}

	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}
	time.Sleep(time.Second * 5)
	fmt.Printf("\r\nid: %s\r\n", resp.ID)
	rx, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Details:    true,
	})
	if err != nil {
		panic(err.Error())
	}
	defer rx.Close()
	target := ctx.Value("resp").(http.ResponseWriter)
	// target := ctx.Value("resp").(http.ResponseWriter)
	// _, err = stdcopy.StdCopy(target, target, rx)
	buf := make([]byte, 1024)
	for {
		len, err := rx.Read(buf[:])
		if err != nil {
			fmt.Print(err)
			fmt.Printf("err   %+v", buf)
			break
		}
		if len == 0 {
			break
		}

		fmt.Print(buf)
		if target != nil {
			target.Write(buf[0:len])
		}
	}
	if err != nil {
		return
	}
	return
}
