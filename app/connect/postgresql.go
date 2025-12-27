package connect

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// PostgreSQLConnection 实现 Connection 接口
type PostgreSQLConnection struct {
	db     *sql.DB
	config Config
}

// NewPostgreSQLConnection 创建一个新的 PostgreSQL 连接
func NewPostgreSQLConnection(config Config) (Connection, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.Username, config.Password, config.Database)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open PostgreSQL connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping PostgreSQL server: %w", err)
	}

	return &PostgreSQLConnection{
		db:     db,
		config: config,
	}, nil
}

// GetDBNames 获取所有数据库信息
func (c *PostgreSQLConnection) GetDBNames() ([]DatabaseInfo, error) {
	query := `SELECT datname FROM pg_database WHERE datistemplate = false`
	rows, err := c.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get databases: %w", err)
	}
	defer rows.Close()

	var databases []DatabaseInfo
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan database name: %w", err)
		}
		databases = append(databases, DatabaseInfo{
			Name: name,
		})
	}
	return databases, nil
}

// GetConfig 获取连接配置
func (c *PostgreSQLConnection) GetConfig() Config {
	return c.config
}

// GetSchemas 获取指定数据库的所有schema
func (c *PostgreSQLConnection) GetSchemas(database string) ([]Schema, error) {
	query := `SELECT schema_name FROM information_schema.schemata 
			  WHERE schema_name NOT IN ('pg_catalog', 'information_schema')`

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
func (c *PostgreSQLConnection) GetTables(params QueryParams) ([]TableInfo, error) {
	query := `SELECT table_name, obj_description((quote_ident(table_schema) || '.' || quote_ident(table_name))::regclass, 'pg_class') as table_comment
			  FROM information_schema.tables 
			  WHERE table_schema = $1 AND table_type = 'BASE TABLE'`

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
func (c *PostgreSQLConnection) GetViews(params QueryParams) ([]ViewInfo, error) {
	query := `SELECT table_name, view_definition
			  FROM information_schema.views
			  WHERE table_schema = $1`

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
func (c *PostgreSQLConnection) GetTableFields(params QueryParams) ([]FieldInfo, error) {
	query := `SELECT column_name, data_type, 
					 col_description((quote_ident($1) || '.' || quote_ident($2))::regclass::oid, ordinal_position) as column_comment,
					 is_nullable, column_default
			  FROM information_schema.columns
			  WHERE table_schema = $1 AND table_name = $2
			  ORDER BY ordinal_position`

	rows, err := c.db.Query(query, params.Schema, params.Table)
	if err != nil {
		return nil, fmt.Errorf("failed to get fields: %w", err)
	}
	defer rows.Close()

	var fields []FieldInfo
	for rows.Next() {
		var field FieldInfo
		var comment, isNullable sql.NullString
		var defaultValue sql.NullString
		if err := rows.Scan(&field.Name, &field.Type, &comment, &isNullable, &defaultValue); err != nil {
			return nil, err
		}
		if comment.Valid {
			field.Comment = comment.String
		}
		if isNullable.Valid && isNullable.String == "YES" {
			field.Nullable = true
		}
		if defaultValue.Valid {
			field.DefaultValue = defaultValue.String
		}
		fields = append(fields, field)
	}

	return fields, nil
}

// Test 测试数据库连接是否可用
func (c *PostgreSQLConnection) Test() error {
	return c.db.Ping()
}

// Close 关闭数据库连接
func (c *PostgreSQLConnection) Close() error {
	return c.db.Close()
}
