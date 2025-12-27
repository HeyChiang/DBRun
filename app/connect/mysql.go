package connect

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLConnection 实现 Connection 接口
type MySQLConnection struct {
	db     *sql.DB
	config Config
}

// NewMySQLConnection 创建一个新的 MySQL 连接
func NewMySQLConnection(config Config) (Connection, error) {
	// 构建 DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database)

	// 打开数据库连接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open MySQL connection: %w", err)
	}

	// 测试连接
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping MySQL server: %w", err)
	}

	return &MySQLConnection{
		db:     db,
		config: config,
	}, nil
}

// 实现 Connection 接口的方法

func (c *MySQLConnection) GetDBNames() ([]DatabaseInfo, error) {
	rows, err := c.db.Query("SHOW DATABASES")
	if err != nil {
		return nil, fmt.Errorf("failed to get databases: %w", err)
	}
	defer rows.Close()

	var databases []DatabaseInfo
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, fmt.Errorf("failed to scan database name: %w", err)
		}
		databases = append(databases, DatabaseInfo{Name: dbName})
	}

	return databases, nil
}

func (c *MySQLConnection) GetTables(params QueryParams) ([]TableInfo, error) {
	query := `
        SELECT TABLE_NAME, TABLE_COMMENT
        FROM INFORMATION_SCHEMA.TABLES
        WHERE TABLE_SCHEMA = ?
        AND TABLE_TYPE = 'BASE TABLE'
    `
	rows, err := c.db.Query(query, params.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to get tables: %w", err)
	}
	defer rows.Close()

	var tables []TableInfo
	for rows.Next() {
		var table TableInfo
		if err := rows.Scan(&table.Name, &table.Comment); err != nil {
			return nil, fmt.Errorf("failed to scan table info: %w", err)
		}
		tables = append(tables, table)
	}

	return tables, nil
}

func (c *MySQLConnection) GetViews(params QueryParams) ([]ViewInfo, error) {
	query := `
        SELECT TABLE_NAME, VIEW_DEFINITION
        FROM INFORMATION_SCHEMA.VIEWS
        WHERE TABLE_SCHEMA = ?
    `
	rows, err := c.db.Query(query, params.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to get views: %w", err)
	}
	defer rows.Close()

	var views []ViewInfo
	for rows.Next() {
		var view ViewInfo
		if err := rows.Scan(&view.Name, &view.Definition); err != nil {
			return nil, fmt.Errorf("failed to scan view info: %w", err)
		}
		views = append(views, view)
	}

	return views, nil
}

func (c *MySQLConnection) GetTableFields(params QueryParams) ([]FieldInfo, error) {
	query := `
        SELECT COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, COLUMN_KEY, COLUMN_COMMENT
        FROM INFORMATION_SCHEMA.COLUMNS
        WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
    `
	rows, err := c.db.Query(query, params.Database, params.Table)
	if err != nil {
		return nil, fmt.Errorf("failed to get table fields: %w", err)
	}
	defer rows.Close()

	var fields []FieldInfo
	for rows.Next() {
		var field FieldInfo
		var isNullable string
		if err := rows.Scan(&field.Name, &field.Type, &isNullable, &field.Key, &field.Comment); err != nil {
			return nil, fmt.Errorf("failed to scan field info: %w", err)
		}
		field.Nullable = (isNullable == "YES")
		fields = append(fields, field)
	}

	return fields, nil
}

// GetSchemas 获取指定数据库的所有schema
// MySQL doesn't have true schema support like Oracle, so we return an empty schema list
func (c *MySQLConnection) GetSchemas(database string) ([]Schema, error) {
	// MySQL doesn't use schemas in the same way as Oracle
	// For MySQL, database is equivalent to schema
	return []Schema{}, nil
}

// GetConfig 获取 MySQL 连接配置
func (c *MySQLConnection) GetConfig() Config {
	return c.config
}

// Test 测试数据库连接是否可用
func (c *MySQLConnection) Test() error {
	return c.db.Ping()
}

func (c *MySQLConnection) Close() error {
	return c.db.Close()
}
