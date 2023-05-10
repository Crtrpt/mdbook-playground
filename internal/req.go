package internal

type ReqForm struct {
	Code    string `json:"code"` //代码
	Repo    string `json:"repo"` //仓库
	Path    string `json:"file"` //路径
	Project string `json:"project"`
	Image   string `json:"image"` //使用的镜像
}
