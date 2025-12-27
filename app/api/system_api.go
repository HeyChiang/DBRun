package api

import (
	"dbrun/app/service"
	"context"
)

// SystemServiceAPI 封装对外的系统服务 API
// 可以在这里做参数校验、统一错误处理等
// 也方便前端直接调用

type SystemServiceAPI struct {
	ctx context.Context
}

func NewSystemServiceAPI() *SystemServiceAPI {
	return &SystemServiceAPI{}
}

func (api *SystemServiceAPI) Init(ctx context.Context) {
	api.ctx = ctx
}

// 打开系统目录选择对话框
func (api *SystemServiceAPI) OpenDirectory() (string, error) {
	return service.OpenDirectory(api.ctx)
}

// 创建目录
func (api *SystemServiceAPI) CreateDirectory(dirPath string) error {
	return service.CreateDirectory(dirPath)
}

// 检查目录是否存在
func (api *SystemServiceAPI) PathExists(dirPath string) (bool, error) {
	return service.PathExists(dirPath)
}
