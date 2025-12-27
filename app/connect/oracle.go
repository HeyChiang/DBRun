package connect

import (
	"database/sql"
	"fmt"

	go_ora "github.com/sijms/go-ora/v2"
)

// OracleConnection 实现 Connection 接口
type OracleConnection struct {
	db     *sql.DB
	config Config
}

// NewOracleConnection 创建一个新的 Oracle 连接
func NewOracleConnection(config Config) (Connection, error) {
	// 使用go-ora的BuildUrl方法构建连接字符串
	connStr := go_ora.BuildUrl(config.Host, config.Port, config.Database, config.Username, config.Password, nil)

	// 打开数据库连接
	db, err := sql.Open("oracle", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open Oracle connection: %w", err)
	}

	// 测试连接
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping Oracle server: %w", err)
	}

	return &OracleConnection{
		db:     db,
		config: config,
	}, nil
}

// GetDBNames 获取所有数据库信息
func (c *OracleConnection) GetDBNames() ([]DatabaseInfo, error) {
    return []DatabaseInfo{{Name: c.config.Database}}, nil
}

// GetConfig 返回连接配置信息
func (c *OracleConnection) GetConfig() Config {
	return c.config
}

// GetSchemas 获取指定数据库的所有schema
func (c *OracleConnection) GetSchemas(database string) ([]Schema, error) {
	rows, err := c.db.Query(`
		SELECT USERNAME 
		FROM ALL_USERS 
		WHERE ORACLE_MAINTAINED = 'N' 
		ORDER BY USERNAME
	`)
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
		schemas = append(schemas, Schema{
			Name: schemaName,
		})
	}
	return schemas, nil
}

// GetTables 获取指定schema的所有表
func (c *OracleConnection) GetTables(params QueryParams) ([]TableInfo, error) {
	query := `
		SELECT t.TABLE_NAME, c.COMMENTS
		FROM ALL_TABLES t
		LEFT JOIN ALL_TAB_COMMENTS c ON t.TABLE_NAME = c.TABLE_NAME AND t.OWNER = c.OWNER
		WHERE t.OWNER = :1 
		ORDER BY t.TABLE_NAME
	`

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
func (c *OracleConnection) GetViews(params QueryParams) ([]ViewInfo, error) {
	query := `
		SELECT VIEW_NAME, TEXT 
		FROM ALL_VIEWS 
		WHERE OWNER = :1 
		ORDER BY VIEW_NAME
	`
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
func (c *OracleConnection) GetTableFields(params QueryParams) ([]FieldInfo, error) {
	// First, get all columns
	columnsQuery := `
		SELECT c.COLUMN_NAME, c.DATA_TYPE, c.NULLABLE
		FROM ALL_TAB_COLUMNS c
		WHERE c.OWNER = :1 AND c.TABLE_NAME = :2
	`
    rows, err := c.db.Query(columnsQuery, params.Schema, params.Table)
	if err != nil {
		return nil, fmt.Errorf("failed to get table fields: %w", err)
	}
	defer rows.Close()

	// Create a map to store field information
	fieldMap := make(map[string]FieldInfo)

	for rows.Next() {
		var field FieldInfo
		var nullable string
		if err := rows.Scan(&field.Name, &field.Type, &nullable); err != nil {
			return nil, err
		}
		field.Nullable = nullable == "Y"
		fieldMap[field.Name] = field
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Query for constraints (PRI, FOR, UNI)
	constraintsQuery := `
		SELECT cols.COLUMN_NAME, cons.CONSTRAINT_TYPE
		FROM ALL_CONSTRAINTS cons
		JOIN ALL_CONS_COLUMNS cols ON cons.OWNER = cols.OWNER
			AND cons.CONSTRAINT_NAME = cols.CONSTRAINT_NAME
			AND cons.TABLE_NAME = cols.TABLE_NAME
		WHERE cons.OWNER = :1
			AND cons.TABLE_NAME = :2
			AND cons.CONSTRAINT_TYPE IN ('P', 'R', 'U')
	`
    rows, err = c.db.Query(constraintsQuery, params.Schema, params.Table)
	if err != nil {
		return nil, fmt.Errorf("failed to get constraints: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var columnName string
		var constraintType string
		if err := rows.Scan(&columnName, &constraintType); err != nil {
			return nil, err
		}

		if field, ok := fieldMap[columnName]; ok {
			switch constraintType {
			case "P":
				field.Key = "PRI"
			case "R":
				if field.Key != "PRI" {
					field.Key = "FOR"
				}
			case "U":
				if field.Key != "PRI" && field.Key != "FOR" {
					field.Key = "UNI"
				}
			}
			fieldMap[columnName] = field // Update the map
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Query for indexes
	indexQuery := `
		SELECT DISTINCT ind.column_name
		FROM ALL_IND_COLUMNS ind
		WHERE ind.table_owner = :1
			AND ind.table_name = :2
	`
    rows, err = c.db.Query(indexQuery, params.Schema, params.Table)
	if err != nil {
		return nil, fmt.Errorf("failed to get indexes: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var columnName string
		if err := rows.Scan(&columnName); err != nil {
			return nil, err
		}

		if field, ok := fieldMap[columnName]; ok {
			if field.Key == "" {
				field.Key = "IDX"
			}
			fieldMap[columnName] = field // Update the map
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Build the fields slice from the map
	var fields []FieldInfo
	for _, field := range fieldMap {
		fields = append(fields, field)
	}

	return fields, nil
}

// Test 测试数据库连接是否可用
func (c *OracleConnection) Test() error {
	return c.db.Ping()
}

// Close 关闭数据库连接
func (c *OracleConnection) Close() error {
	return c.db.Close()
}
