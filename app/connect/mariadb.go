package connect

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// MariaDBConnection 实现 Connection 接口
type MariaDBConnection struct {
	db     *sql.DB
	config Config
}

// NewMariaDBConnection 创建一个新的 MariaDB 连接
func NewMariaDBConnection(config Config) (Connection, error) {
	// 构建 DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database)

	// 打开数据库连接
	db, err := sql.Open("mysql", dsn) // MariaDB 使用 MySQL 驱动
	if err != nil {
		return nil, fmt.Errorf("failed to open MariaDB connection: %w", err)
	}

	// 测试连接并验证是否为 MariaDB
	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to get MariaDB version: %w", err)
	}

	// 验证是否为 MariaDB
	if !isMariaDB(version) {
		db.Close()
		return nil, fmt.Errorf("database is not MariaDB (version: %s)", version)
	}

	return &MariaDBConnection{db: db, config: config}, nil
}

// isMariaDB 检查版本字符串是否为 MariaDB
func isMariaDB(version string) bool {
	return len(version) >= 7 && version[0:7] == "MariaDB"
}

// GetDBNames 获取所有数据库信息
func (c *MariaDBConnection) GetDBNames() ([]DatabaseInfo, error) {
	rows, err := c.db.Query("SHOW DATABASES")
	if err != nil {
		return nil, fmt.Errorf("failed to get databases: %w", err)
	}
	defer rows.Close()

	var databases []DatabaseInfo
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, err
		}

		// 跳过系统数据库
		if dbName == "information_schema" || dbName == "mysql" || dbName == "performance_schema" || dbName == "sys" {
			continue
		}

		tables, err := c.GetTables(QueryParams{Database: dbName})
		if err != nil {
			return nil, err
		}

		views, err := c.GetViews(QueryParams{Database: dbName})
		if err != nil {
			return nil, err
		}

		databases = append(databases, DatabaseInfo{
			Name:   dbName,
			Tables: tables,
			Views:  views,
		})
	}

	return databases, nil
}

// GetTables 获取指定数据库的所有表
func (c *MariaDBConnection) GetTables(params QueryParams) ([]TableInfo, error) {
	query := `
		SELECT 
			table_name,
			table_comment
		FROM information_schema.tables 
		WHERE table_schema = ? 
		AND table_type = 'BASE TABLE'`

	rows, err := c.db.Query(query, params.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to get tables: %w", err)
	}
	defer rows.Close()

	var tables []TableInfo
	for rows.Next() {
		var table TableInfo
		if err := rows.Scan(&table.Name, &table.Comment); err != nil {
			return nil, err
		}

		// 获取表字段信息
		fields, err := c.GetTableFields(QueryParams{
			Database: params.Database,
			Table:    table.Name,
		})
		if err != nil {
			return nil, err
		}
		table.Fields = fields

		tables = append(tables, table)
	}

	return tables, nil
}

// GetViews 获取指定数据库的所有视图
func (c *MariaDBConnection) GetViews(params QueryParams) ([]ViewInfo, error) {
	query := `
		SELECT 
			table_name,
			view_definition
		FROM information_schema.views 
		WHERE table_schema = ?`

	rows, err := c.db.Query(query, params.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to get views: %w", err)
	}
	defer rows.Close()

	var views []ViewInfo
	for rows.Next() {
		var view ViewInfo
		if err := rows.Scan(&view.Name, &view.Definition); err != nil {
			return nil, err
		}
		views = append(views, view)
	}

	return views, nil
}

// GetTableFields 获取指定表的所有字段信息
func (c *MariaDBConnection) GetTableFields(params QueryParams) ([]FieldInfo, error) {
	query := `
		SELECT 
			column_name,
			column_type,
			is_nullable,
			column_key,
			column_comment,
			column_default
		FROM information_schema.columns 
		WHERE table_schema = ? 
		AND table_name = ? 
		ORDER BY ordinal_position`

	rows, err := c.db.Query(query, params.Database, params.Table)
	if err != nil {
		return nil, fmt.Errorf("failed to get fields: %w", err)
	}
	defer rows.Close()

	var fields []FieldInfo
	for rows.Next() {
		var field FieldInfo
		var isNullable string
		var defaultValue sql.NullString
		if err := rows.Scan(
			&field.Name,
			&field.Type,
			&isNullable,
			&field.Key,
			&field.Comment,
			&defaultValue,
		); err != nil {
			return nil, err
		}

		field.Nullable = (isNullable == "YES")
		if defaultValue.Valid {
			field.DefaultValue = defaultValue.String
		}

		fields = append(fields, field)
	}

	return fields, nil
}

// GetSchemas 获取指定数据库的所有schema
// MariaDB 不支持真正的 schema，返回空列表
func (c *MariaDBConnection) GetSchemas(database string) ([]Schema, error) {
	return []Schema{}, nil
}

// GetConfig 获取连接配置
func (c *MariaDBConnection) GetConfig() Config {
	return c.config
}

// Test 测试数据库连接是否可用
func (c *MariaDBConnection) Test() error {
	return c.db.Ping()
}

// Close 关闭数据库连接
func (c *MariaDBConnection) Close() error {
	return c.db.Close()
}
