package service

import (
    "dbrun/app/connect"
    "fmt"
)

// 更新相关操作：同步表字段、更新数据库原始信息到VO、更新表/字段备注

// SyncTableFieldsByTableID 通过表ID同步该表字段（刷新原始并同步到VO）
func SyncTableFieldsByTableID(tableID int64) error {
	manager, err := getMgr()
	if err != nil {
		return err
	}
	rs := manager.GetRawStorage()
	configID, dbName, schemaName, tableName, databaseID, schemaID, err := rs.GetTableContextByID(tableID)
	if err != nil {
		return fmt.Errorf("resolve table context failed: %w", err)
	}
	fmt.Printf("[SyncTableFieldsByTableID] tableID=%d configID=%d dbName=%s schemaName=%s tableName=%s databaseID=%d schemaID=%v\n", tableID, configID, dbName, schemaName, tableName, databaseID, schemaID)

	// 优先尝试从连接池复用现有连接（无需依赖类型字段）
	if pooled, ok := connect.GetConnectionFromPool(configID); ok {
		fmt.Printf("[SyncTableFieldsByTableID] reusing pooled connection for configID=%d\n", configID)
		rawDB, err := fetchRawDatabase(pooled, dbName)
		if err != nil {
			return err
		}
		if err := manager.UpdateRawDatabase(configID, rawDB); err != nil {
			return err
		}
		return manager.SyncRawToVO(configID)
	}

	creds, err := manager.GetCredentialsByID(configID)
	if err != nil {
		return fmt.Errorf("get credentials failed: %w", err)
	}
	fmt.Printf("[SyncTableFieldsByTableID] fetched credentials id=%d type=%q label=%q host=%q port=%d db=%q instance=%q options=%q\n", creds.ID, creds.Type, creds.Label, creds.Host, creds.Port, creds.Database, creds.Instance, creds.Options)
	if creds.Type == "" {
		return fmt.Errorf("missing database type in credentials for id %d", configID)
	}
	cfg := connect.Config{
		ID:        creds.ID,
		Type:      creds.Type,
		Label:     creds.Label,
		Username:  creds.Username,
		Password:  creds.Password,
		Host:      creds.Host,
		Port:      creds.Port,
		Database:  creds.Database,
		Instance:  creds.Instance,
		Options:   creds.Options,
		CreatedAt: creds.CreatedAt,
	}

	conn, err := connect.GetConnection(cfg)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}
	fmt.Printf("[SyncTableFieldsByTableID] created connection for configID=%d type=%q, fetching raw database=%q\n", configID, cfg.Type, dbName)
	rawDB, err := fetchRawDatabase(conn, dbName)
	if err != nil {
		return err
	}
	if err := manager.UpdateRawDatabase(configID, rawDB); err != nil {
		return err
	}
	return manager.SyncRawToVO(configID)
}

// SyncSchemaByID 通过原始SchemaID同步该Schema（刷新原始并同步到VO）
func SyncSchemaByID(schemaID int64) error {
	manager, err := getMgr()
	if err != nil {
		return err
	}
	rs := manager.GetRawStorage()
	configID, dbName, schemaName, databaseID, err := rs.GetSchemaContextByID(schemaID)
	if err != nil {
		return fmt.Errorf("resolve schema context failed: %w", err)
	}
	fmt.Printf("[SyncSchemaByID] schemaID=%d configID=%d dbName=%s schemaName=%s databaseID=%d\n", schemaID, configID, dbName, schemaName, databaseID)

	// 优先复用现有连接
	if pooled, ok := connect.GetConnectionFromPool(configID); ok {
		fmt.Printf("[SyncSchemaByID] reusing pooled connection for configID=%d\n", configID)
		rawDB, err := fetchRawDatabase(pooled, dbName)
		if err != nil {
			return err
		}
		if err := manager.UpdateRawDatabase(configID, rawDB); err != nil {
			return err
		}
		return manager.SyncRawToVO(configID)
	}

	creds, err := manager.GetCredentialsByID(configID)
	if err != nil {
		return fmt.Errorf("get credentials failed: %w", err)
	}
	fmt.Printf("[SyncSchemaByID] fetched credentials id=%d type=%q label=%q host=%q port=%d db=%q instance=%q options=%q\n", creds.ID, creds.Type, creds.Label, creds.Host, creds.Port, creds.Database, creds.Instance, creds.Options)
	if creds.Type == "" {
		return fmt.Errorf("missing database type in credentials for id %d", configID)
	}
	cfg := connect.Config{
		ID:        creds.ID,
		Type:      creds.Type,
		Label:     creds.Label,
		Username:  creds.Username,
		Password:  creds.Password,
		Host:      creds.Host,
		Port:      creds.Port,
		Database:  creds.Database,
		Instance:  creds.Instance,
		Options:   creds.Options,
		CreatedAt: creds.CreatedAt,
	}

	conn, err := connect.GetConnection(cfg)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}
	fmt.Printf("[SyncSchemaByID] created connection for configID=%d type=%q, fetching raw database=%q\n", configID, cfg.Type, dbName)
	rawDB, err := fetchRawDatabase(conn, dbName)
	if err != nil {
		return err
	}
	if err := manager.UpdateRawDatabase(configID, rawDB); err != nil {
		return err
	}
	return manager.SyncRawToVO(configID)
}

// SyncDatabaseByID 通过原始数据库ID同步该数据库（刷新原始并同步到VO）
func SyncDatabaseByID(databaseID int64) error {
    manager, err := getMgr()
    if err != nil {
        return err
    }
    rs := manager.GetRawStorage()
    configID, dbName, err := rs.GetDatabaseContextByID(databaseID)
    if err != nil {
        return fmt.Errorf("resolve database context failed: %w", err)
    }
    fmt.Printf("[SyncDatabaseByID] databaseID=%d configID=%d dbName=%s\n", databaseID, configID, dbName)

    // 优先复用现有连接
    if pooled, ok := connect.GetConnectionFromPool(configID); ok {
        fmt.Printf("[SyncDatabaseByID] reusing pooled connection for configID=%d\n", configID)
        rawDB, err := fetchRawDatabase(pooled, dbName)
        if err != nil {
            return err
        }
        if err := manager.UpdateRawDatabase(configID, rawDB); err != nil {
            return err
        }
        return manager.SyncRawToVO(configID)
    }

    creds, err := manager.GetCredentialsByID(configID)
    if err != nil {
        return fmt.Errorf("get credentials failed: %w", err)
    }
    if creds.Type == "" {
        return fmt.Errorf("missing database type in credentials for id %d", configID)
    }
    cfg := connect.Config{
        ID:        creds.ID,
        Type:      creds.Type,
        Label:     creds.Label,
        Username:  creds.Username,
        Password:  creds.Password,
        Host:      creds.Host,
        Port:      creds.Port,
        Database:  creds.Database,
        Instance:  creds.Instance,
        Options:   creds.Options,
        CreatedAt: creds.CreatedAt,
    }

    conn, err := connect.GetConnection(cfg)
    if err != nil {
        return fmt.Errorf("failed to get connection: %w", err)
    }
    fmt.Printf("[SyncDatabaseByID] created connection for configID=%d type=%q, fetching raw database=%q\n", configID, cfg.Type, dbName)
    rawDB, err := fetchRawDatabase(conn, dbName)
    if err != nil {
        return err
    }
    if err := manager.UpdateRawDatabase(configID, rawDB); err != nil {
        return err
    }
    return manager.SyncRawToVO(configID)
}

// 已移除名称版同步与全量更新方法，统一采用 ID 驱动的同步接口

// （已移除）SetTableVOCache/ClearTableVOCache：前端未调用的名称版方法

// SetTableVOCacheByTableID 写入VO备注（通过表ID）
func SetTableVOCacheByTableID(tableID int64, tableVO TableCacheVO) {
	manager, err := getMgr()
	if err != nil {
		return
	}
	vs := manager.GetVOStorage()
	_ = vs.UpdateTableRemarkByID(tableID, tableVO.Remark)
	for fieldName, remark := range tableVO.FieldMap {
		_ = vs.UpdateFieldRemarkByTableIDName(tableID, fieldName, remark)
	}
}

// ClearTableVOCacheByTableID 清除表备注（通过表ID）
func ClearTableVOCacheByTableID(tableID int64) {
    manager, err := getMgr()
    if err != nil {
        return
    }
    vs := manager.GetVOStorage()
    _ = vs.UpdateTableRemarkByID(tableID, "")
}

// （已清理）移除未使用的扩展信息更新方法，改由缓存接口统一覆盖

// SyncRawToVO 将原始数据同步到VO存储（合并自 metadata_service.go）
func (m *MetadataService) SyncRawToVO(configID int64) error {
    dbInfoVO, err := m.ConvertRawToVO(configID)
    if err != nil {
        return err
    }
    for _, db := range dbInfoVO.DBs {
        if _, err := m.voStorage.SaveDatabaseInfoVO(configID, db); err != nil {
            return err
        }
    }
    return nil
}
