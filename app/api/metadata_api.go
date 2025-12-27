package api

import (
    "dbrun/app/connect"
    "dbrun/app/models"
    "dbrun/app/service"
)

// API 结构体，用于暴露 前端使用的方法
type MetadatasAPI struct{}

// NewAPI 创建一个新的API实例
func NewMetadatasAPI() *MetadatasAPI {
	return &MetadatasAPI{}
}

// ListDatabasesByConfig 方法获取数据库表的元数据
func (a *MetadatasAPI) ListDatabasesByConfig(config connect.Config) (models.DBInfoVO, error) {
	return service.ListDatabasesByConfig(config)
}

// Test 方法测试连接是否可用
func (a *MetadatasAPI) TestConnection(config connect.Config) error {
	return service.TestConnection(config)
}

// InitCache 初始化缓存（按项目目录）
func (a *MetadatasAPI) InitCache(projectDir string) error {
	return service.InitMetadataStorageService(projectDir)
}

func (a *MetadatasAPI) CloseAllConnections() error {
    return connect.CloseAllConnections()
}

// SyncTableFieldsByTableID 通过表ID同步该表字段
func (a *MetadatasAPI) SyncTableFieldsByTableID(tableID int64) error {
    return service.SyncTableFieldsByTableID(tableID)
}

// SyncSchemaByID 通过schemaID同步整个schema
func (a *MetadatasAPI) SyncSchemaByID(schemaID int64) error {
    return service.SyncSchemaByID(schemaID)
}

func (a *MetadatasAPI) SyncDatabaseByID(databaseID int64) error {
    return service.SyncDatabaseByID(databaseID)
}

// GetFieldsVOByTableID 根据表ID获取字段VO
func (a *MetadatasAPI) GetFieldsVOByTableID(tableID int64) ([]models.FieldInfoVO, error) {
    return service.GetFieldsVOByTableID(tableID)
}

// GetTablesVOByDatabaseID 根据数据库ID与可选的SchemaID获取表VO
func (a *MetadatasAPI) GetTablesVOByDatabaseID(databaseID int64, schemaID *int64) ([]models.TableInfoVO, error) {
    return service.GetTablesVOByDatabaseID(databaseID, schemaID)
}

// GetTableVOCacheByTableID 通过表ID获取TableVO缓存
func (a *MetadatasAPI) GetTableVOCacheByTableID(tableID int64) (service.TableCacheVO, bool) {
    tableVO := service.GetTableVOCacheByTableID(tableID)
    return tableVO, tableVO.Remark != "" || len(tableVO.FieldMap) > 0
}

// SetTableVOCacheByTableID 缓存TableVO信息（ID关联）
func (a *MetadatasAPI) SetTableVOCacheByTableID(tableID int64, tableVO service.TableCacheVO) error {
    service.SetTableVOCacheByTableID(tableID, tableVO)
    return nil
}

// ClearTableVOCacheByTableID 清除特定表的缓存（ID关联）
func (a *MetadatasAPI) ClearTableVOCacheByTableID(tableID int64) error {
    service.ClearTableVOCacheByTableID(tableID)
    return nil
}

// ParseViewSQL 解析视图SQL语句，返回字段信息
func (a *MetadatasAPI) ParseViewSQL(sql string) ([]models.FieldInfoVO, error) {
    return service.ParseViewSQL(sql)
}

// UpdateFieldsSortByTableID 根据表ID批量更新字段排序
func (a *MetadatasAPI) UpdateFieldsSortByTableID(tableID int64, fieldIDs []int64) error {
    return service.UpdateFieldsSortByTableID(tableID, fieldIDs)
}
