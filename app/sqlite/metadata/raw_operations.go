package metadata

import (
	"dbrun/app/connect"
	"gorm.io/gorm"
)

// RawMetadataStorage 原始元数据存储管理器
type RawMetadataStorage struct {
    db *gorm.DB
}

// NewRawMetadataStorage 创建原始元数据存储管理器
func NewRawMetadataStorage(db *gorm.DB) *RawMetadataStorage {
	return &RawMetadataStorage{db: db}
}

// SaveDatabaseInfo 保存数据库信息
func (r *RawMetadataStorage) SaveDatabaseInfo(configID int64, dbInfo connect.DatabaseInfo) (*RawDatabaseInfo, error) {
	rawDB := &RawDatabaseInfo{
		ConfigID: configID,
		Name:     dbInfo.Name,
		Comment:  dbInfo.Comment,
	}

	// 先查找是否已存在
	var existing RawDatabaseInfo
	err := r.db.Where("config_id = ? AND name = ?", configID, dbInfo.Name).First(&existing).Error
	if err == nil {
		// 更新现有记录
		existing.Comment = dbInfo.Comment
		err = r.db.Save(&existing).Error
		if err != nil {
			return nil, err
		}
		rawDB = &existing
	} else if err == gorm.ErrRecordNotFound {
		// 创建新记录
		err = r.db.Create(rawDB).Error
		if err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	// 保存Schema信息
	for _, schema := range dbInfo.Schemas {
		_, err := r.SaveSchemaInfo(rawDB.ID, schema)
		if err != nil {
			return nil, err
		}
	}

	// 保存表信息
	for _, table := range dbInfo.Tables {
		_, err := r.SaveTableInfo(rawDB.ID, nil, table)
		if err != nil {
			return nil, err
		}
	}

	// 保存视图信息
	for _, view := range dbInfo.Views {
		_, err := r.SaveViewInfo(rawDB.ID, nil, view)
		if err != nil {
			return nil, err
		}
	}

	return rawDB, nil
}

// SaveSchemaInfo 保存Schema信息
func (r *RawMetadataStorage) SaveSchemaInfo(databaseID int64, schema connect.Schema) (*RawSchemaInfo, error) {
	rawSchema := &RawSchemaInfo{
		DatabaseID: databaseID,
		Name:       schema.Name,
	}

	// 先查找是否已存在
	var existing RawSchemaInfo
	err := r.db.Where("database_id = ? AND name = ?", databaseID, schema.Name).First(&existing).Error
	if err == nil {
		rawSchema = &existing
	} else if err == gorm.ErrRecordNotFound {
		// 创建新记录
		err = r.db.Create(rawSchema).Error
		if err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	// 保存Schema下的表信息
	for _, table := range schema.Tables {
		_, err := r.SaveTableInfo(databaseID, &rawSchema.ID, table)
		if err != nil {
			return nil, err
		}
	}

	// 保存Schema下的视图信息
	for _, view := range schema.Views {
		_, err := r.SaveViewInfo(databaseID, &rawSchema.ID, view)
		if err != nil {
			return nil, err
		}
	}

	return rawSchema, nil
}

// SaveTableInfo 保存表信息
func (r *RawMetadataStorage) SaveTableInfo(databaseID int64, schemaID *int64, table connect.TableInfo) (*RawTableInfo, error) {
	rawTable := &RawTableInfo{
		DatabaseID: databaseID,
		SchemaID:   schemaID,
		Name:       table.Name,
		Comment:    table.Comment,
	}

	// 先查找是否已存在
	var existing RawTableInfo
	query := r.db.Where("database_id = ? AND name = ?", databaseID, table.Name)
	if schemaID != nil {
		query = query.Where("schema_id = ?", *schemaID)
	} else {
		query = query.Where("schema_id IS NULL")
	}
	
	err := query.First(&existing).Error
	if err == nil {
		// 更新现有记录
		existing.Comment = table.Comment
		err = r.db.Save(&existing).Error
		if err != nil {
			return nil, err
		}
		rawTable = &existing
	} else if err == gorm.ErrRecordNotFound {
		// 创建新记录
		err = r.db.Create(rawTable).Error
		if err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	// 先删除该表的所有现有字段记录
	if err := r.db.Where("table_id = ?", rawTable.ID).Delete(&RawFieldInfo{}).Error; err != nil {
		return nil, err
	}

	// 然后重新插入所有字段
	for _, field := range table.Fields {
		rawField := &RawFieldInfo{
			TableID:      rawTable.ID,
			Name:         field.Name,
			Type:         field.Type,
			Nullable:     field.Nullable,
			Key:          field.Key,
			DefaultValue: field.DefaultValue,
			Comment:      field.Comment,
		}
		if err := r.db.Create(rawField).Error; err != nil {
			return nil, err
		}
	}

	return rawTable, nil
}

// SaveFieldInfo 保存字段信息
func (r *RawMetadataStorage) SaveFieldInfo(tableID int64, field connect.FieldInfo) (*RawFieldInfo, error) {
	rawField := &RawFieldInfo{
		TableID:      tableID,
		Name:         field.Name,
		Type:         field.Type,
		Nullable:     field.Nullable,
		Key:          field.Key,
		Comment:      field.Comment,
		DefaultValue: field.DefaultValue,
	}

	// 先查找是否已存在
	var existing RawFieldInfo
	err := r.db.Where("table_id = ? AND name = ?", tableID, field.Name).First(&existing).Error
	if err == nil {
		// 更新现有记录
		existing.Type = field.Type
		existing.Nullable = field.Nullable
		existing.Key = field.Key
		existing.Comment = field.Comment
		existing.DefaultValue = field.DefaultValue
		err = r.db.Save(&existing).Error
		if err != nil {
			return nil, err
		}
		rawField = &existing
	} else if err == gorm.ErrRecordNotFound {
		// 创建新记录
		err = r.db.Create(rawField).Error
		if err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return rawField, nil
}

// SaveViewInfo 保存视图信息
func (r *RawMetadataStorage) SaveViewInfo(databaseID int64, schemaID *int64, view connect.ViewInfo) (*RawViewInfo, error) {
	rawView := &RawViewInfo{
		DatabaseID: databaseID,
		SchemaID:   schemaID,
		Name:       view.Name,
		Definition: view.Definition,
	}

	// 先查找是否已存在
	var existing RawViewInfo
	query := r.db.Where("database_id = ? AND name = ?", databaseID, view.Name)
	if schemaID != nil {
		query = query.Where("schema_id = ?", *schemaID)
	} else {
		query = query.Where("schema_id IS NULL")
	}
	
	err := query.First(&existing).Error
	if err == nil {
		// 更新现有记录
		existing.Definition = view.Definition
		err = r.db.Save(&existing).Error
		if err != nil {
			return nil, err
		}
		rawView = &existing
	} else if err == gorm.ErrRecordNotFound {
		// 创建新记录
		err = r.db.Create(rawView).Error
		if err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return rawView, nil
}

// GetDatabasesByConfigID 根据配置ID获取数据库信息
func (r *RawMetadataStorage) GetDatabasesByConfigID(configID int64) ([]connect.DatabaseInfo, error) {
	var rawDatabases []RawDatabaseInfo
	err := r.db.Where("config_id = ?", configID).Find(&rawDatabases).Error
	if err != nil {
		return nil, err
	}

	var databases []connect.DatabaseInfo
	for _, rawDB := range rawDatabases {
		db := connect.DatabaseInfo{
			Name:    rawDB.Name,
			Comment: rawDB.Comment,
		}

		// 获取Schema信息
		schemas, err := r.GetSchemasByDatabaseID(rawDB.ID)
		if err != nil {
			return nil, err
		}
		db.Schemas = schemas

		// 获取表信息（不属于任何Schema的）
		tables, err := r.GetTablesByDatabaseID(rawDB.ID, nil)
		if err != nil {
			return nil, err
		}
		db.Tables = tables

		// 获取视图信息（不属于任何Schema的）
		views, err := r.GetViewsByDatabaseID(rawDB.ID, nil)
		if err != nil {
			return nil, err
		}
		db.Views = views

		databases = append(databases, db)
	}

	return databases, nil
}

// GetSchemasByDatabaseID 根据数据库ID获取Schema信息
func (r *RawMetadataStorage) GetSchemasByDatabaseID(databaseID int64) ([]connect.Schema, error) {
	var rawSchemas []RawSchemaInfo
	err := r.db.Where("database_id = ?", databaseID).Find(&rawSchemas).Error
	if err != nil {
		return nil, err
	}

	var schemas []connect.Schema
	for _, rawSchema := range rawSchemas {
		schema := connect.Schema{
			Name: rawSchema.Name,
		}

		// 获取Schema下的表信息
		tables, err := r.GetTablesByDatabaseID(databaseID, &rawSchema.ID)
		if err != nil {
			return nil, err
		}
		schema.Tables = tables

		// 获取Schema下的视图信息
		views, err := r.GetViewsByDatabaseID(databaseID, &rawSchema.ID)
		if err != nil {
			return nil, err
		}
		schema.Views = views

		schemas = append(schemas, schema)
	}

	return schemas, nil
}

// GetTablesByDatabaseID 根据数据库ID和Schema ID获取表信息
func (r *RawMetadataStorage) GetTablesByDatabaseID(databaseID int64, schemaID *int64) ([]connect.TableInfo, error) {
	var rawTables []RawTableInfo
	query := r.db.Where("database_id = ?", databaseID)
	if schemaID != nil {
		query = query.Where("schema_id = ?", *schemaID)
	} else {
		query = query.Where("schema_id IS NULL")
	}
	
	err := query.Find(&rawTables).Error
	if err != nil {
		return nil, err
	}

	var tables []connect.TableInfo
	for _, rawTable := range rawTables {
		table := connect.TableInfo{
			Name:    rawTable.Name,
			Comment: rawTable.Comment,
		}

		// 获取字段信息
		fields, err := r.GetFieldsByTableID(rawTable.ID)
		if err != nil {
			return nil, err
		}
		table.Fields = fields

		tables = append(tables, table)
	}

	return tables, nil
}

// GetFieldsByTableID 根据表ID获取字段信息
func (r *RawMetadataStorage) GetFieldsByTableID(tableID int64) ([]connect.FieldInfo, error) {
	var rawFields []RawFieldInfo
	err := r.db.Where("table_id = ?", tableID).Find(&rawFields).Error
	if err != nil {
		return nil, err
	}

	var fields []connect.FieldInfo
	for _, rawField := range rawFields {
		field := connect.FieldInfo{
			Name:         rawField.Name,
			Type:         rawField.Type,
			Nullable:     rawField.Nullable,
			Key:          rawField.Key,
			Comment:      rawField.Comment,
			DefaultValue: rawField.DefaultValue,
		}
		fields = append(fields, field)
	}

	return fields, nil
}

// GetViewsByDatabaseID 根据数据库ID和Schema ID获取视图信息
func (r *RawMetadataStorage) GetViewsByDatabaseID(databaseID int64, schemaID *int64) ([]connect.ViewInfo, error) {
	var rawViews []RawViewInfo
	query := r.db.Where("database_id = ?", databaseID)
	if schemaID != nil {
		query = query.Where("schema_id = ?", *schemaID)
	} else {
		query = query.Where("schema_id IS NULL")
	}
	
	err := query.Find(&rawViews).Error
	if err != nil {
		return nil, err
	}

	var views []connect.ViewInfo
	for _, rawView := range rawViews {
		view := connect.ViewInfo{
			Name:       rawView.Name,
			Definition: rawView.Definition,
		}
		views = append(views, view)
	}

	return views, nil
}

// DeleteByConfigID 根据配置ID删除所有相关数据
func (r *RawMetadataStorage) DeleteByConfigID(configID int64) error {
	// 获取所有数据库
	var databases []RawDatabaseInfo
	err := r.db.Where("config_id = ?", configID).Find(&databases).Error
	if err != nil {
		return err
	}

	// 删除每个数据库的相关数据
	for _, db := range databases {
		err = r.DeleteByDatabaseID(db.ID)
		if err != nil {
			return err
		}
	}

	// 删除数据库记录
	return r.db.Where("config_id = ?", configID).Delete(&RawDatabaseInfo{}).Error
}

// DeleteByDatabaseID 根据数据库ID删除所有相关数据
func (r *RawMetadataStorage) DeleteByDatabaseID(databaseID int64) error {
	// 删除字段信息
	err := r.db.Where("table_id IN (SELECT id FROM raw_table_info WHERE database_id = ?)", databaseID).Delete(&RawFieldInfo{}).Error
	if err != nil {
		return err
	}

	// 删除表信息
	err = r.db.Where("database_id = ?", databaseID).Delete(&RawTableInfo{}).Error
	if err != nil {
		return err
	}

	// 删除视图信息
	err = r.db.Where("database_id = ?", databaseID).Delete(&RawViewInfo{}).Error
	if err != nil {
		return err
	}

	// 删除Schema信息
    return r.db.Where("database_id = ?", databaseID).Delete(&RawSchemaInfo{}).Error
}

// GetTableContextByID 通过原始表ID获取其上下文信息（configID、数据库/Schema/表名称及相关ID）
func (r *RawMetadataStorage) GetTableContextByID(tableID int64) (configID int64, dbName string, schemaName string, tableName string, databaseID int64, schemaID *int64, err error) {
    var rt RawTableInfo
    if err = r.db.Where("id = ?", tableID).First(&rt).Error; err != nil {
        return 0, "", "", "", 0, nil, err
    }
    var rd RawDatabaseInfo
    if err = r.db.Where("id = ?", rt.DatabaseID).First(&rd).Error; err != nil {
        return 0, "", "", "", 0, nil, err
    }
    var rs RawSchemaInfo
    if rt.SchemaID != nil {
        if err = r.db.Where("id = ?", *rt.SchemaID).First(&rs).Error; err != nil {
            return 0, "", "", "", 0, nil, err
        }
        schemaName = rs.Name
    }
    return rd.ConfigID, rd.Name, schemaName, rt.Name, rt.DatabaseID, rt.SchemaID, nil
}

// GetSchemaContextByID 通过原始SchemaID获取其上下文信息（configID、数据库/Schema名称及相关ID）
func (r *RawMetadataStorage) GetSchemaContextByID(schemaID int64) (configID int64, dbName string, schemaName string, databaseID int64, err error) {
    var rs RawSchemaInfo
    if err = r.db.Where("id = ?", schemaID).First(&rs).Error; err != nil {
        return 0, "", "", 0, err
    }
    var rd RawDatabaseInfo
    if err = r.db.Where("id = ?", rs.DatabaseID).First(&rd).Error; err != nil {
        return 0, "", "", 0, err
    }
    return rd.ConfigID, rd.Name, rs.Name, rs.DatabaseID, nil
}

// GetDatabaseContextByID 通过原始数据库ID获取其上下文信息（configID、数据库名称）
func (r *RawMetadataStorage) GetDatabaseContextByID(databaseID int64) (configID int64, dbName string, err error) {
    var rd RawDatabaseInfo
    if err := r.db.Where("id = ?", databaseID).First(&rd).Error; err != nil {
        return 0, "", err
    }
    return rd.ConfigID, rd.Name, nil
}

// ===== 便捷方法：直接返回原始行（含ID），供VO转换阶段使用 =====

// GetRawDatabasesRows 根据配置ID返回原始数据库行
func (r *RawMetadataStorage) GetRawDatabasesRows(configID int64) ([]RawDatabaseInfo, error) {
    var rows []RawDatabaseInfo
    if err := r.db.Where("config_id = ?", configID).Find(&rows).Error; err != nil {
        return nil, err
    }
    return rows, nil
}

// GetRawSchemasRows 根据数据库ID返回原始Schema行
func (r *RawMetadataStorage) GetRawSchemasRows(databaseID int64) ([]RawSchemaInfo, error) {
    var rows []RawSchemaInfo
    if err := r.db.Where("database_id = ?", databaseID).Find(&rows).Error; err != nil {
        return nil, err
    }
    return rows, nil
}

// GetRawTablesRows 根据数据库ID和可选SchemaID返回原始表行
func (r *RawMetadataStorage) GetRawTablesRows(databaseID int64, schemaID *int64) ([]RawTableInfo, error) {
    var rows []RawTableInfo
    q := r.db.Where("database_id = ?", databaseID)
    if schemaID != nil {
        q = q.Where("schema_id = ?", *schemaID)
    } else {
        q = q.Where("schema_id IS NULL")
    }
    if err := q.Find(&rows).Error; err != nil {
        return nil, err
    }
    return rows, nil
}

// GetRawViewsRows 根据数据库ID和可选SchemaID返回原始视图行
func (r *RawMetadataStorage) GetRawViewsRows(databaseID int64, schemaID *int64) ([]RawViewInfo, error) {
    var rows []RawViewInfo
    q := r.db.Where("database_id = ?", databaseID)
    if schemaID != nil {
        q = q.Where("schema_id = ?", *schemaID)
    } else {
        q = q.Where("schema_id IS NULL")
    }
    if err := q.Find(&rows).Error; err != nil {
        return nil, err
    }
    return rows, nil
}

// GetRawFieldsRows 根据表ID返回原始字段行
func (r *RawMetadataStorage) GetRawFieldsRows(tableID int64) ([]RawFieldInfo, error) {
    var rows []RawFieldInfo
    if err := r.db.Where("table_id = ?", tableID).Find(&rows).Error; err != nil {
        return nil, err
    }
    return rows, nil
}