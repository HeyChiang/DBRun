package api

import (
	"dbrun/app/cache"
	"context"
)

// AppCacheApi 应用缓存API
type AppCacheApi struct {
    ctx context.Context
}

// NewAppCacheApi 创建新的应用缓存API实例
func NewAppCacheApi() *AppCacheApi {
	return &AppCacheApi{}
}

// SetContext 设置上下文
func (a *AppCacheApi) SetContext(ctx context.Context) {
    a.ctx = ctx
}

// Init 初始化应用缓存目录（由前端选择的项目目录）
func (a *AppCacheApi) Init(projectDir string) error {
    return cache.InitAppCache(projectDir)
}

// Set 设置缓存
func (a *AppCacheApi) Set(key string, value string) {
    cache.Set(key, value)
}

// Get 获取缓存值，返回值和是否存在
// 注意：为避免前端指针参数序列化问题，这里不使用指针入参
func (a *AppCacheApi) Get(key string) (string, bool) {
    var s string
    ok := cache.Get(key, &s)
    return s, ok
}
