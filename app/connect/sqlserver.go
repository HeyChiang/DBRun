package connect

import (
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
)

// SQLServerConnection 实现 Connection 接口
type SQLServerConnection struct {
	db     *sql.DB
	config Config
}

// NewSQLServerConnection 创建一个新的 SQL Server 连接
func NewSQLServerConnection(config Config) (Connection, error) {
	dsn := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s",
		config.Host, config.Username, config.Password, config.Port, config.Database)
	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQL Server connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping SQL Server: %w", err)
	}

	return &SQLServerConnection{
		db:     db,
		config: config,
	}, nil
}

// GetDBNames 获取所有数据库信息
func (c *SQLServerConnection) GetDBNames() ([]DatabaseInfo, error) {
	query := `SELECT name FROM sys.databases WHERE database_id > 4`
	rows, err := c.db.Query(query)
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

		schemas, err := c.GetSchemas(dbName)
		if err != nil {
			return nil, err
		}

		databases = append(databases, DatabaseInfo{
			Name:    dbName,
			Schemas: schemas,
		})
	}

	return databases, nil
}

// GetSchemas 获取指定数据库的所有schema
func (c *SQLServerConnection) GetSchemas(database string) ([]Schema, error) {
	query := fmt.Sprintf("USE [%s]; SELECT name FROM sys.schemas WHERE name NOT IN ('sys', 'guest', 'INFORMATION_SCHEMA')", database)
	rows, err := c.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get schemas: %w", err)
	}
	defer rows.Close()

	var schemas []Schema
	for rows.Next() {
		var schemaName string
		if err := rows.Scan(&schemaName); err != nil {
			return nil, err
		}
		schemas = append(schemas, Schema{Name: schemaName})
	}

	return schemas, nil
}

// GetTables 获取指定schema的所有表
func (c *SQLServerConnection) GetTables(params QueryParams) ([]TableInfo, error) {
	query := `SELECT t.name, CAST(ep.value AS NVARCHAR(MAX)) as table_comment
			  FROM sys.tables t
			  LEFT JOIN sys.extended_properties ep ON ep.major_id = t.object_id 
			  	AND ep.minor_id = 0 AND ep.name = 'MS_Description'
			  WHERE t.schema_id = SCHEMA_ID(@p1)`

	rows, err := c.db.Query(query, params.Schema)
	if err != nil {
		return nil, fmt.Errorf("failed to get tables: %w", err)
	}
	defer rows.Close()

	var tables []TableInfo
	for rows.Next() {
		var table TableInfo
		var comment sql.NullString
		if err := rows.Scan(&table.Name, &comment); err != nil {
			return nil, err
		}
		if comment.Valid {
			table.Comment = comment.String
		}
		tables = append(tables, table)
	}

	return tables, nil
}

// GetViews 获取指定schema的所有视图
func (c *SQLServerConnection) GetViews(params QueryParams) ([]ViewInfo, error) {
	query := `SELECT v.name, m.definition
			  FROM sys.views v
			  INNER JOIN sys.sql_modules m ON v.object_id = m.object_id
			  WHERE v.schema_id = SCHEMA_ID(@p1)`

	rows, err := c.db.Query(query, params.Schema)
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
func (c *SQLServerConnection) GetTableFields(params QueryParams) ([]FieldInfo, error) {
	query := `SELECT c.name, t.name as data_type,
					 CAST(ep.value AS NVARCHAR(MAX)) as column_comment,
					 c.is_nullable,
					 OBJECT_DEFINITION(c.default_object_id) as column_default
			  FROM sys.columns c
			  INNER JOIN sys.types t ON c.user_type_id = t.user_type_id
			  INNER JOIN sys.objects o ON c.object_id = o.object_id
			  LEFT JOIN sys.extended_properties ep ON ep.major_id = c.object_id 
			  	AND ep.minor_id = c.column_id AND ep.name = 'MS_Description'
			  WHERE o.schema_id = SCHEMA_ID(@p1) AND o.name = @p2
			  ORDER BY c.column_id`

	rows, err := c.db.Query(query, params.Schema, params.Table)
	if err != nil {
		return nil, fmt.Errorf("failed to get fields: %w", err)
	}
	defer rows.Close()

	var fields []FieldInfo
	for rows.Next() {
		var field FieldInfo
		var comment sql.NullString
		var isNullable bool
		var defaultValue sql.NullString
		if err := rows.Scan(&field.Name, &field.Type, &comment, &isNullable, &defaultValue); err != nil {
			return nil, err
		}
		if comment.Valid {
			field.Comment = comment.String
		}
		field.Nullable = isNullable
		if defaultValue.Valid {
			field.DefaultValue = defaultValue.String
		}
		fields = append(fields, field)
	}

	return fields, nil
}

// GetConfig 获取连接配置
func (c *SQLServerConnection) GetConfig() Config {
	return c.config
}

// Test 测试数据库连接是否可用
func (c *SQLServerConnection) Test() error {
	return c.db.Ping()
}

// Close 关闭数据库连接
func (c *SQLServerConnection) Close() error {
	return c.db.Close()
}
