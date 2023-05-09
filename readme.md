mdbook playground 提供对各种语言的支持 后端使用docker

参考 https://cs.opensource.google/go/x/playground/+/14ebe15b:sandbox/sandbox.go

请求参数
```json
{
    "image":"go",
    "path":"",
    "code":""
}
```

文件替换和运行

删除镜像
```
docker rm $(docker ps -aq)
```