package metadata

import (
	"time"
)

// VODatabaseInfo VO数据库信息表
type VODatabaseInfo struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ConfigID  int64     `gorm:"not null;index" json:"config_id"`        // 关联的配置ID
	Alias     string    `gorm:"size:255" json:"alias"`                  // 数据库别名（扩展字段）
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// VOSchemaInfo VO Schema信息表（VO仅做扩展，不复制原始字段）
type VOSchemaInfo struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	DatabaseID int64     `gorm:"not null;index" json:"database_id"`      // 关联的数据库ID
	Alias      string    `gorm:"size:255" json:"alias"`                  // Schema别名（扩展字段）
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// VOTableInfo VO表信息表（仅扩展字段，如颜色与备注）
type VOTableInfo struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	DatabaseID int64     `gorm:"not null;index" json:"database_id"`      // 关联的数据库ID
	SchemaID   *int64    `gorm:"index" json:"schema_id"`                 // 关联的Schema ID（可为空）
	Alias      string    `gorm:"size:255" json:"alias"`                  // 表别名（扩展字段）
	Color      string    `gorm:"size:50" json:"color"`                   // 表颜色
	Remark     string    `gorm:"size:2000" json:"remark"`                // 表备注
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// VOFieldInfo VO字段信息表（仅扩展字段，如显示与备注）
type VOFieldInfo struct {
    ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
    TableID    int64     `gorm:"not null;index" json:"table_id"`        // 关联的表ID
    Alias      string    `gorm:"size:255" json:"alias"`                 // 字段别名（扩展字段）
    Display    bool      `gorm:"default:true" json:"display"`           // 是否显示
    Remark     string    `gorm:"size:2000" json:"remark"`               // 字段备注
    FontColor  string    `gorm:"size:50" json:"font_color"`             // 字体颜色
    BgColor    string    `gorm:"size:50" json:"bg_color"`               // 背景颜色
    Sort       int       `gorm:"default:0;index" json:"sort"`           // 排序（升序）
    CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// VOViewInfo VO视图信息表（仅扩展字段，如颜色与备注）
type VOViewInfo struct {
    ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
    DatabaseID int64     `gorm:"not null;index" json:"database_id"`      // 关联的数据库ID
    SchemaID   *int64    `gorm:"index" json:"schema_id"`                 // 关联的Schema ID（可为空）
    Alias      string    `gorm:"size:255" json:"alias"`                  // 视图别名（扩展字段）
    Color      string    `gorm:"size:50" json:"color"`                   // 视图颜色
    Remark     string    `gorm:"size:2000" json:"remark"`               // 视图备注
    Sort       int       `gorm:"default:0;index" json:"sort"`           // 排序（升序）
    CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// VODisplay 显示配置表
type VODisplay struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ConfigID  int64     `gorm:"not null;index" json:"config_id"`        // 关联的配置ID
	DBCnt     int       `gorm:"default:0" json:"db_cnt"`                // 数据库数量
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 方法用于指定表名
func (VODatabaseInfo) TableName() string {
	return "vo_database_info"
}

func (VOSchemaInfo) TableName() string {
	return "vo_schema_info"
}

func (VOTableInfo) TableName() string {
	return "vo_table_info"
}

func (VOFieldInfo) TableName() string {
	return "vo_field_info"
}

func (VOViewInfo) TableName() string {
	return "vo_view_info"
}

func (VODisplay) TableName() string {
	return "vo_display"
}