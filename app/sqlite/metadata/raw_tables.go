package metadata

import (
	"time"
)

// RawDatabaseInfo 原始数据库信息表
type RawDatabaseInfo struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ConfigID  int64     `gorm:"not null;index" json:"config_id"`        // 关联的配置ID
	Name      string    `gorm:"not null;size:255" json:"name"`          // 数据库名称
	Comment   string    `gorm:"size:1000" json:"comment"`               // 数据库注释
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// RawSchemaInfo 原始Schema信息表
type RawSchemaInfo struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	DatabaseID int64     `gorm:"not null;index" json:"database_id"`      // 关联的数据库ID
	Name       string    `gorm:"not null;size:255" json:"name"`          // Schema名称
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// RawTableInfo 原始表信息表
type RawTableInfo struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	DatabaseID int64     `gorm:"not null;index" json:"database_id"`      // 关联的数据库ID
	SchemaID   *int64    `gorm:"index" json:"schema_id"`                 // 关联的Schema ID（可为空）
	Name       string    `gorm:"not null;size:255" json:"name"`          // 表名
	Comment    string    `gorm:"size:1000" json:"comment"`               // 表注释
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// RawFieldInfo 原始字段信息表
type RawFieldInfo struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TableID      int64     `gorm:"not null;index" json:"table_id"`        // 关联的表ID
	Name         string    `gorm:"not null;size:255" json:"name"`         // 字段名称
	Type         string    `gorm:"not null;size:100" json:"type"`         // 字段类型
	Nullable     bool      `gorm:"default:true" json:"nullable"`          // 是否可为空
	Key          string    `gorm:"size:50" json:"key"`                    // 字段索引类型
	Comment      string    `gorm:"size:1000" json:"comment"`              // 字段注释
	DefaultValue string    `gorm:"size:500" json:"default_value"`         // 默认值
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// RawViewInfo 原始视图信息表
type RawViewInfo struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	DatabaseID int64     `gorm:"not null;index" json:"database_id"`      // 关联的数据库ID
	SchemaID   *int64    `gorm:"index" json:"schema_id"`                 // 关联的Schema ID（可为空）
	Name       string    `gorm:"not null;size:255" json:"name"`          // 视图名称
	Definition string    `gorm:"type:text" json:"definition"`            // 视图定义
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 方法用于指定表名
func (RawDatabaseInfo) TableName() string {
	return "raw_database_info"
}

func (RawSchemaInfo) TableName() string {
	return "raw_schema_info"
}

func (RawTableInfo) TableName() string {
	return "raw_table_info"
}

func (RawFieldInfo) TableName() string {
	return "raw_field_info"
}

func (RawViewInfo) TableName() string {
	return "raw_view_info"
}