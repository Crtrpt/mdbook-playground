package internal

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/robfig/cron/v3"
	"github.com/rs/xid"
)

var Cli *client.Client
var ImageList map[string]*types.ImageSummary

var Cfg *Config

func InitConfig() {
	Cfg = &Config{}
	_, err := toml.DecodeFile("./app.toml", Cfg)
	if err != nil {
		panic(err)
	}
}

// 初始化docker客户端
func InitDockerClient() {
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

// 自动删除容器
func AutoRemoveCloseContainer() (res any, err error) {
	c := cron.New(cron.WithSeconds())
	fmt.Printf("启动自动删除 %s\r\n", Cfg.AutoRemoveCron)
	_, err = c.AddFunc(Cfg.AutoRemoveCron, func() {
		fmt.Printf("执行自动删除容器 %s\r\n", Cfg.AutoRemoveBefore)
		ctx := context.Background()
		list, err := Cli.ContainerList(ctx, types.ContainerListOptions{
			Before: Cfg.AutoRemoveBefore,
			All:    true,
		})
		if err != nil {
			panic(err)
		}
		for _, c := range list {
			err = Cli.ContainerRemove(context.Background(), c.ID, types.ContainerRemoveOptions{})
			if err != nil {
				fmt.Printf("err:%+v", err)
			}
		}
	})
	if err != nil {
		panic(err)
	}
	c.Start()
	return
}

// 启动一个容器
func StartContainer(ctx context.Context, image, code string) (res string, err error) {
	cid := xid.New().String()
	cfg := ctx.Value("cfg").(*Config)
	cli := ctx.Value("client").(*client.Client)
	req := ctx.Value("req").(ReqForm)
	v := make(map[string]struct{})
	fmt.Printf("image:%v source:%v \r\n", image, cfg.RepoDir+"/example/"+req.Project)
	// 挂载
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Cmd:     strslice.StrSlice{"make"},
		Volumes: v,
		Image:   image},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:     mount.TypeBind,
					ReadOnly: false,
					Target:   "/var/data",
					Source:   cfg.RepoDir + "/example/" + req.Project,
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
	time.Sleep(time.Second * 2)
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
	// _, err = stdcopy.StdCopy(target, target, rx)
	buf := make([]byte, 1024)
	len := 0
	for {
		len, err = rx.Read(buf[:])
		if err != nil {

			break
		}
		if len == 0 {
			break
		}

		if target != nil {
			target.Write(buf[0:len])
		}
	}
	if err != nil {
		return
	}
	return
}
