package connect

import "time"

// Connection 接口定义了所有数据库连接应该实现的方法
type Connection interface {
	GetDBNames() ([]DatabaseInfo, error)
	GetTables(params QueryParams) ([]TableInfo, error)
	GetViews(params QueryParams) ([]ViewInfo, error)
	GetTableFields(params QueryParams) ([]FieldInfo, error)
	GetSchemas(database string) ([]Schema, error)
	GetConfig() Config
	Test() error
	Close() error
}

// QueryParams 统一的查询参数
type QueryParams struct {
	Database string // 数据库名
	Schema   string // schema名称，用于oracle等数据库
	Table    string // 表名
}

// DatabaseInfo 存储数据库的信息
// Updated JSON tags to use lowercase names
type DatabaseInfo struct {
	Name    string      `json:"name"`
	Comment string      `json:"comment"`
	Schemas []Schema    `json:"schemas"`
	Tables  []TableInfo `json:"tables"`
	Views   []ViewInfo  `json:"views"`
}

type Schema struct {
	Name   string      `json:"name"`
	Tables []TableInfo `json:"tables"`
	Views  []ViewInfo  `json:"views"`
}

// ViewInfo 存储视图的信息
// Updated JSON tags to use lowercase names
type ViewInfo struct {
	Name       string `json:"name"`
	Definition string `json:"definition"`
}

// TableInfo 存储表的信息
// Updated JSON tags to use lowercase names
type TableInfo struct {
	Name    string      `json:"name"`
	Comment string      `json:"comment"`
	Fields  []FieldInfo `json:"fields"`
}

// FieldInfo 存储字段的信息
// Updated JSON tags to use lowercase names
type FieldInfo struct {
	// 字段名称
	Name     string `json:"name"`
	// 字段类型
	Type     string `json:"type"`
	Nullable bool   `json:"nullable,omitempty"`
	// 字段索引
	Key          string `json:"key"`
	Comment      string `json:"comment,omitempty"`
	DefaultValue string `json:"default_value,omitempty"`
}

// Config 结构体用于存储数据库连接信息
// 不再依赖 sqlite 包，避免导入循环
type Config struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Label     string `json:"label"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Database  string `json:"database"`
	Instance  string `json:"instance"`
	Options   string `json:"options"`
	CreatedAt time.Time `json:"created_at"`
}
