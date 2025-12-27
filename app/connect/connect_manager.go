package connect

import (
	"fmt"
)

// connectionPool 用于存储已创建的连接
var (
    connectionPool = make(map[int64]Connection)
)

// GetConnectionFromPool 尝试直接从连接池获取连接
func GetConnectionFromPool(id int64) (Connection, bool) {
    if conn, exists := connectionPool[id]; exists {
        fmt.Printf("[ConnectManager] hit connection pool: id=%d type=%q\n", id, conn.GetConfig().Type)
        return conn, true
    }
    return nil, false
}

// GetConnection 函数用于获取或创建数据库连接
func GetConnection(config Config) (Connection, error) {
    if config.ID > 0 {
        if conn, exists := connectionPool[config.ID]; exists {
            fmt.Printf("[ConnectManager] reuse pooled connection: id=%d type=%q\n", config.ID, conn.GetConfig().Type)
            return conn, nil
        }
    }

	var conn Connection
	var err error

    switch config.Type {
    case "mysql":
        conn, err = NewMySQLConnection(config)
    case "oracle":
        conn, err = NewOracleConnection(config)
    case "postgresql":
        conn, err = NewPostgreSQLConnection(config)
    case "sqlserver":
        conn, err = NewSQLServerConnection(config)
    case "mariadb":
        conn, err = NewMariaDBConnection(config)
    default:
        fmt.Printf("[ConnectManager] unsupported database type: id=%d type=%q\n", config.ID, config.Type)
        return nil, fmt.Errorf("unsupported database type: %s", config.Type)
    }

    if err != nil {
        fmt.Printf("[ConnectManager] create connection failed: id=%d type=%q err=%v\n", config.ID, config.Type, err)
        return nil, err
    }

    if config.ID > 0 {
        connectionPool[config.ID] = conn
        fmt.Printf("[ConnectManager] new connection created and pooled: id=%d type=%q\n", config.ID, config.Type)
    } else {
        fmt.Printf("[ConnectManager] new ephemeral connection created: id=%d type=%q\n", config.ID, config.Type)
    }
    return conn, nil
}

// CloseAllConnections 关闭所有数据库连接
func CloseAllConnections() error {
	for key, conn := range connectionPool {
		if err := conn.Close(); err != nil {
			// 即使出错也继续关闭其他连接
			return fmt.Errorf("error closing connection %d: %v", key, err)
		}
		delete(connectionPool, key)
	}
	return nil
}
