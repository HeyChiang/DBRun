package metadata

import (
	"dbrun/app/models"
	"gorm.io/gorm"
)

// VOMetadataStorage VO元数据存储管理器
type VOMetadataStorage struct {
	db *gorm.DB
}

// NewVOMetadataStorage 创建VO元数据存储管理器
func NewVOMetadataStorage(db *gorm.DB) *VOMetadataStorage {
	return &VOMetadataStorage{db: db}
}

// SaveDatabaseInfoVO 保存数据库信息VO（VO仅记录存在关系，不复制原始字段）
func (v *VOMetadataStorage) SaveDatabaseInfoVO(configID int64, dbInfo models.DatabaseInfoVO) (*VODatabaseInfo, error) {
    // 严格通过原始数据按名称解析原始ID
    var rawDB RawDatabaseInfo
    if err := v.db.Where("config_id = ? AND name = ?", configID, dbInfo.Name).First(&rawDB).Error; err != nil {
        return nil, err
    }

    voDB := &VODatabaseInfo{
        ID:       rawDB.ID,
        ConfigID: configID,
    }

    // 仅按ID查找并更新/创建
    var existing VODatabaseInfo
    if err := v.db.Where("id = ?", voDB.ID).First(&existing).Error; err == nil {
        existing.ConfigID = configID
        if err := v.db.Save(&existing).Error; err != nil {
            return nil, err
        }
        voDB = &existing
    } else if err == gorm.ErrRecordNotFound {
        if err := v.db.Create(voDB).Error; err != nil {
            return nil, err
        }
    } else {
        return nil, err
    }

    // 保存Schema信息
    for _, schema := range dbInfo.Schemas {
        if _, err := v.SaveSchemaInfoVO(voDB.ID, schema); err != nil {
            return nil, err
        }
    }

    // 保存表信息
    for _, table := range dbInfo.Tables {
        if _, err := v.SaveTableInfoVO(voDB.ID, nil, table); err != nil {
            return nil, err
        }
    }

    // 保存视图信息
    for _, view := range dbInfo.Views {
        if _, err := v.SaveViewInfoVO(voDB.ID, nil, view); err != nil {
            return nil, err
        }
    }

    return voDB, nil
}

// SaveSchemaInfoVO 保存Schema信息VO
func (v *VOMetadataStorage) SaveSchemaInfoVO(databaseID int64, schema models.SchemaVO) (*VOSchemaInfo, error) {
	voSchema := &VOSchemaInfo{
		DatabaseID: databaseID,
	}

	// 通过原始数据定位原始Schema的ID
	var rawSchema RawSchemaInfo
	if err := v.db.Where("database_id = ? AND name = ?", databaseID, schema.Name).First(&rawSchema).Error; err == nil {
		voSchema.ID = rawSchema.ID
	}

	// 先查找是否已存在（优先按ID）
	var existing VOSchemaInfo
	var err error
	if voSchema.ID != 0 {
		err = v.db.Where("id = ?", voSchema.ID).First(&existing).Error
	} else {
		// 不再按名称匹配VO，若无法解析ID，则视为不存在
		err = gorm.ErrRecordNotFound
	}
	if err == nil {
		// 更新扩展字段（目前仅保持 DatabaseID）
		existing.DatabaseID = databaseID
		if err = v.db.Save(&existing).Error; err != nil {
			return nil, err
		}
		voSchema = &existing
	} else if err == gorm.ErrRecordNotFound {
		if err = v.db.Create(voSchema).Error; err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	// 保存Schema下的表信息
	for _, table := range schema.Tables {
		if _, err := v.SaveTableInfoVO(databaseID, &voSchema.ID, table); err != nil {
			return nil, err
		}
	}

	// 保存Schema下的视图信息
	for _, view := range schema.Views {
		if _, err := v.SaveViewInfoVO(databaseID, &voSchema.ID, view); err != nil {
			return nil, err
		}
	}

	return voSchema, nil
}

// SaveTableInfoVO 保存表信息VO（仅扩展字段）
func (v *VOMetadataStorage) SaveTableInfoVO(databaseID int64, schemaID *int64, table models.TableInfoVO) (*VOTableInfo, error) {
	voTable := &VOTableInfo{
		DatabaseID: databaseID,
		SchemaID:   schemaID,
		Color:      table.Color,
		Remark:     table.Remark,
	}

	// 通过原始数据定位对应的表ID
	var rawTable RawTableInfo
	queryRaw := v.db.Where("database_id = ? AND name = ?", databaseID, table.Name)
	if schemaID != nil {
		queryRaw = queryRaw.Where("schema_id = ?", *schemaID)
	} else {
		queryRaw = queryRaw.Where("schema_id IS NULL")
	}
	if err := queryRaw.First(&rawTable).Error; err == nil {
		voTable.ID = rawTable.ID
	}

	// 先查找是否已存在（优先按ID）
	var existing VOTableInfo
	var err error
	if voTable.ID != 0 {
		err = v.db.Where("id = ?", voTable.ID).First(&existing).Error
	} else {
		// 不再按名称匹配VO，若无法解析ID，则视为不存在
		err = gorm.ErrRecordNotFound
	}
	if err == nil {
		// 更新扩展字段
		existing.Color = table.Color
		existing.Remark = table.Remark
		existing.DatabaseID = databaseID
		existing.SchemaID = schemaID
		if err = v.db.Save(&existing).Error; err != nil {
			return nil, err
		}
		voTable = &existing
	} else if err == gorm.ErrRecordNotFound {
		if err = v.db.Create(voTable).Error; err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	// 保存字段信息
	for _, field := range table.Fields {
		if _, err := v.SaveFieldInfoVO(voTable.ID, field); err != nil {
			return nil, err
		}
	}

	return voTable, nil
}

// SaveFieldInfoVO 保存字段信息VO（仅扩展字段）
func (v *VOMetadataStorage) SaveFieldInfoVO(tableID int64, field models.FieldInfoVO) (*VOFieldInfo, error) {
	voField := &VOFieldInfo{
		TableID: tableID,
		Display: field.Display,
		Remark:  field.Remark,
	}

	// 通过原始数据定位字段ID
	var rawField RawFieldInfo
	if err := v.db.Where("table_id = ? AND name = ?", tableID, field.Name).First(&rawField).Error; err == nil {
		voField.ID = rawField.ID
	}

	// 先查找是否已存在（优先按ID）
	var existing VOFieldInfo
	var err error
	if voField.ID != 0 {
		err = v.db.Where("id = ?", voField.ID).First(&existing).Error
	} else {
		// 不再按名称匹配VO，若无法解析ID，则视为不存在
		err = gorm.ErrRecordNotFound
	}
	if err == nil {
		// 更新扩展字段
		existing.Display = field.Display
		existing.Remark = field.Remark
		existing.TableID = tableID
		if err = v.db.Save(&existing).Error; err != nil {
			return nil, err
		}
		voField = &existing
	} else if err == gorm.ErrRecordNotFound {
		if err = v.db.Create(voField).Error; err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return voField, nil
}

// SaveViewInfoVO 保存视图信息VO（仅扩展字段）
func (v *VOMetadataStorage) SaveViewInfoVO(databaseID int64, schemaID *int64, view models.ViewInfoVO) (*VOViewInfo, error) {
	voView := &VOViewInfo{
		DatabaseID: databaseID,
		SchemaID:   schemaID,
		Color:      view.Color,
		Remark:     view.Remark,
	}

	// 通过原始数据定位视图ID
	var rawView RawViewInfo
	queryRaw := v.db.Where("database_id = ? AND name = ?", databaseID, view.Name)
	if schemaID != nil {
		queryRaw = queryRaw.Where("schema_id = ?", *schemaID)
	} else {
		queryRaw = queryRaw.Where("schema_id IS NULL")
	}
	if err := queryRaw.First(&rawView).Error; err == nil {
		voView.ID = rawView.ID
	}

	// 先查找是否已存在（优先按ID）
	var existing VOViewInfo
	var err error
	if voView.ID != 0 {
		err = v.db.Where("id = ?", voView.ID).First(&existing).Error
	} else {
		// 不再按名称匹配VO，若无法解析ID，则视为不存在
		err = gorm.ErrRecordNotFound
	}
	if err == nil {
		// 更新扩展字段
		existing.Color = view.Color
		existing.Remark = view.Remark
		existing.DatabaseID = databaseID
		existing.SchemaID = schemaID
		if err = v.db.Save(&existing).Error; err != nil {
			return nil, err
		}
		voView = &existing
	} else if err == gorm.ErrRecordNotFound {
		if err = v.db.Create(voView).Error; err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return voView, nil
}

// SaveDisplayVO 保存显示配置VO
func (v *VOMetadataStorage) SaveDisplayVO(configID int64, display models.Display) (*VODisplay, error) {
	voDisplay := &VODisplay{
		ConfigID: configID,
		DBCnt:    display.DBCnt,
	}

	// 先查找是否已存在
	var existing VODisplay
	err := v.db.Where("config_id = ?", configID).First(&existing).Error
	if err == nil {
		// 更新现有记录
		existing.DBCnt = display.DBCnt
		err = v.db.Save(&existing).Error
		if err != nil {
			return nil, err
		}
		voDisplay = &existing
	} else if err == gorm.ErrRecordNotFound {
		// 创建新记录
		err = v.db.Create(voDisplay).Error
		if err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	// 已移除样式VO（VOStyle），此处不再保存 display.Style
	return voDisplay, nil
}

// SaveStyleVO 已移除：样式配置不再持久化（保留注释占位，避免误用）

// GetDatabasesVOByConfigID 根据配置ID获取数据库信息VO
func (v *VOMetadataStorage) GetDatabasesVOByConfigID(configID int64) ([]models.DatabaseInfoVO, error) {
    type dbRow struct {
        ID      int64
        ConfigID int64
        Alias   string
        Name    string
        Comment string
    }
    var rows []dbRow
    err := v.db.Table("vo_database_info").
        Select("vo_database_info.id as id, vo_database_info.config_id as config_id, vo_database_info.alias as alias, raw_database_info.name as name, raw_database_info.comment as comment").
        Joins("JOIN raw_database_info ON raw_database_info.id = vo_database_info.id").
        Where("vo_database_info.config_id = ?", configID).
        Scan(&rows).Error
    if err != nil {
        return nil, err
    }

    var databases []models.DatabaseInfoVO
    for _, r := range rows {
        db := models.DatabaseInfoVO{
            ID:      r.ID,
            ConfigID: r.ConfigID,
            Alias:   r.Alias,
            Name:    r.Name,
            Comment: r.Comment,
        }

        schemas, err := v.GetSchemasVOByDatabaseID(r.ID)
        if err != nil { return nil, err }
        db.Schemas = schemas

        tables, err := v.GetTablesVOByDatabaseID(r.ID, nil)
        if err != nil { return nil, err }
        db.Tables = tables

        views, err := v.GetViewsVOByDatabaseID(r.ID, nil)
        if err != nil { return nil, err }
        db.Views = views

        databases = append(databases, db)
    }

    return databases, nil
}

// GetSchemasVOByDatabaseID 根据数据库ID获取Schema信息VO
func (v *VOMetadataStorage) GetSchemasVOByDatabaseID(databaseID int64) ([]models.SchemaVO, error) {
    type schemaRow struct {
        ID   int64
        Alias string
        Name string
    }
    var rows []schemaRow
    err := v.db.Table("vo_schema_info").
        Select("vo_schema_info.id as id, vo_schema_info.alias as alias, raw_schema_info.name as name").
        Joins("JOIN raw_schema_info ON raw_schema_info.id = vo_schema_info.id").
        Where("vo_schema_info.database_id = ?", databaseID).
        Scan(&rows).Error
    if err != nil {
        return nil, err
    }

    var schemas []models.SchemaVO
    for _, r := range rows {
        schema := models.SchemaVO{
            ID:         r.ID,
            DatabaseID: databaseID,
            Alias:      r.Alias,
            Name: r.Name,
        }

        tables, err := v.GetTablesVOByDatabaseID(databaseID, &r.ID)
        if err != nil { return nil, err }
        schema.Tables = tables

        views, err := v.GetViewsVOByDatabaseID(databaseID, &r.ID)
        if err != nil { return nil, err }
        schema.Views = views

        schemas = append(schemas, schema)
    }

    return schemas, nil
}

// GetTablesVOByDatabaseID 根据数据库ID获取表信息VO（仅扩展字段）
func (v *VOMetadataStorage) GetTablesVOByDatabaseID(databaseID int64, schemaID *int64) ([]models.TableInfoVO, error) {
    type tableRow struct {
        ID      int64
        Name    string
        Comment string
        Color   string
        Remark  string
        DatabaseID int64
        SchemaID   *int64
        Alias      string
    }
    query := v.db.Table("vo_table_info").
        Select("vo_table_info.id as id, raw_table_info.name as name, raw_table_info.comment as comment, vo_table_info.color as color, vo_table_info.remark as remark, vo_table_info.database_id as database_id, vo_table_info.schema_id as schema_id, vo_table_info.alias as alias").
        Joins("JOIN raw_table_info ON raw_table_info.id = vo_table_info.id").
        Where("vo_table_info.database_id = ?", databaseID)
    if schemaID != nil {
        query = query.Where("vo_table_info.schema_id = ?", *schemaID)
    } else {
        query = query.Where("vo_table_info.schema_id IS NULL")
    }

    var rows []tableRow
    if err := query.Scan(&rows).Error; err != nil {
        return nil, err
    }

    var tables []models.TableInfoVO
    for _, r := range rows {
        table := models.TableInfoVO{
            ID:      r.ID,
            DatabaseID: r.DatabaseID,
            SchemaID:   r.SchemaID,
            Alias:      r.Alias,
            Name:    r.Name,
            Color:   r.Color,
            Remark:  r.Remark,
            Comment: r.Comment,
        }
        fields, err := v.GetFieldsVOByTableID(r.ID)
        if err != nil { return nil, err }
        table.Fields = fields
        tables = append(tables, table)
    }

    return tables, nil
}

// GetFieldsVOByTableID 根据表ID获取字段信息VO
func (v *VOMetadataStorage) GetFieldsVOByTableID(tableID int64) ([]models.FieldInfoVO, error) {
    type fieldRow struct {
        ID        int64
        Name      string
        Display   bool
        Remark    string
        TableID   int64
        Alias     string
        FontColor string
        BgColor   string
        Sort      int
    }
    var rows []fieldRow
    err := v.db.Table("vo_field_info").
        Select("vo_field_info.id as id, raw_field_info.name as name, vo_field_info.display as display, vo_field_info.remark as remark, vo_field_info.table_id as table_id, vo_field_info.alias as alias, vo_field_info.font_color as font_color, vo_field_info.bg_color as bg_color, vo_field_info.sort as sort").
        Joins("JOIN raw_field_info ON raw_field_info.id = vo_field_info.id").
        Where("vo_field_info.table_id = ?", tableID).
        Order("vo_field_info.sort ASC, vo_field_info.id ASC").
        Scan(&rows).Error
    if err != nil {
        return nil, err
    }
    var fields []models.FieldInfoVO
    for _, r := range rows {
        fields = append(fields, models.FieldInfoVO{
            ID:           r.ID,
            TableID:      r.TableID,
            Alias:        r.Alias,
            Name:         r.Name,
            Display:      r.Display,
            Remark:       r.Remark,
            FontColor:    r.FontColor,
            BgColor:      r.BgColor,
            Sort:         r.Sort,
        })
    }
    return fields, nil
}

// GetViewsVOByDatabaseID 根据数据库ID获取视图信息VO（仅扩展字段）
func (v *VOMetadataStorage) GetViewsVOByDatabaseID(databaseID int64, schemaID *int64) ([]models.ViewInfoVO, error) {
    type viewRow struct {
        ID    int64
        Name  string
        Color string
        Remark string
        DatabaseID int64
        SchemaID   *int64
        Alias      string
        Sort       int
    }
    query := v.db.Table("vo_view_info").
        Select("vo_view_info.id as id, raw_view_info.name as name, vo_view_info.color as color, vo_view_info.remark as remark, vo_view_info.database_id as database_id, vo_view_info.schema_id as schema_id, vo_view_info.alias as alias, vo_view_info.sort as sort").
        Joins("JOIN raw_view_info ON raw_view_info.id = vo_view_info.id").
        Where("vo_view_info.database_id = ?", databaseID)
    if schemaID != nil {
        query = query.Where("vo_view_info.schema_id = ?", *schemaID)
    } else {
        query = query.Where("vo_view_info.schema_id IS NULL")
    }

    var rows []viewRow
    if err := query.Order("vo_view_info.sort ASC, vo_view_info.id ASC").Scan(&rows).Error; err != nil {
        return nil, err
    }

    var views []models.ViewInfoVO
    for _, r := range rows {
        views = append(views, models.ViewInfoVO{
            ID:         r.ID,
            DatabaseID: r.DatabaseID,
            SchemaID:   r.SchemaID,
            Alias:      r.Alias,
            Name:       r.Name,
            Color:      r.Color,
            Remark:     r.Remark,
        })
    }
    return views, nil
}

// DeleteVOByConfigID 根据配置ID删除所有VO相关数据
func (v *VOMetadataStorage) DeleteVOByConfigID(configID int64) error {
	// 先查找所有数据库ID
	var voDatabases []VODatabaseInfo
	if err := v.db.Where("config_id = ?", configID).Find(&voDatabases).Error; err != nil {
		return err
	}

	// 按数据库级联删除相关数据
	for _, db := range voDatabases {
		if err := v.DeleteVOByDatabaseID(db.ID); err != nil {
			return err
		}
	}

	// 删除数据库记录
	return v.db.Where("config_id = ?", configID).Delete(&VODatabaseInfo{}).Error
}

// DeleteVOByDatabaseID 根据数据库ID删除所有VO相关数据
func (v *VOMetadataStorage) DeleteVOByDatabaseID(databaseID int64) error {
	// 删除字段信息
	if err := v.db.Where("table_id IN (SELECT id FROM vo_table_info WHERE database_id = ?)", databaseID).Delete(&VOFieldInfo{}).Error; err != nil {
		return err
	}

	// 删除表信息
	if err := v.db.Where("database_id = ?", databaseID).Delete(&VOTableInfo{}).Error; err != nil {
		return err
	}

	// 删除视图信息
	if err := v.db.Where("database_id = ?", databaseID).Delete(&VOViewInfo{}).Error; err != nil {
		return err
	}

	// 删除Schema信息
	return v.db.Where("database_id = ?", databaseID).Delete(&VOSchemaInfo{}).Error
}



// UpdateTableRemarkByID 根据表ID更新备注
func (v *VOMetadataStorage) UpdateTableRemarkByID(tableID int64, remark string) error {
	res := v.db.Model(&VOTableInfo{}).Where("id = ?", tableID).Update("remark", remark)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		// 若 VO 记录不存在，按 Raw 信息懒创建后再写入
		var rt RawTableInfo
		if err := v.db.Where("id = ?", tableID).First(&rt).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil
			}
			return err
		}
		vo := VOTableInfo{
			ID:         tableID,
			DatabaseID: rt.DatabaseID,
			SchemaID:   rt.SchemaID,
			Remark:     remark,
		}
		if err := v.db.Create(&vo).Error; err != nil {
			return err
		}
	}
	return nil
}

// UpdateFieldRemarkByID 根据字段ID更新备注
func (v *VOMetadataStorage) UpdateFieldRemarkByID(fieldID int64, remark string) error {
    res := v.db.Model(&VOFieldInfo{}).Where("id = ?", fieldID).Update("remark", remark)
    if res.Error != nil {
        return res.Error
    }
    if res.RowsAffected == 0 {
        var rf RawFieldInfo
        if err := v.db.Where("id = ?", fieldID).First(&rf).Error; err != nil {
            if err == gorm.ErrRecordNotFound {
                return nil
            }
            return err
        }
        vo := VOFieldInfo{
            ID:      fieldID,
            TableID: rf.TableID,
            Remark:  remark,
        }
        if err := v.db.Create(&vo).Error; err != nil {
            return err
        }
    }
    return nil
}

// UpdateFieldSortByID 根据字段ID更新排序（如VO不存在则懒创建）
func (v *VOMetadataStorage) UpdateFieldSortByID(fieldID int64, sort int) error {
    res := v.db.Model(&VOFieldInfo{}).Where("id = ?", fieldID).Update("sort", sort)
    if res.Error != nil {
        return res.Error
    }
    if res.RowsAffected == 0 {
        var rf RawFieldInfo
        if err := v.db.Where("id = ?", fieldID).First(&rf).Error; err != nil {
            if err == gorm.ErrRecordNotFound {
                return nil
            }
            return err
        }
        vo := VOFieldInfo{
            ID:      fieldID,
            TableID: rf.TableID,
            Sort:    sort,
        }
        if err := v.db.Create(&vo).Error; err != nil {
            return err
        }
    }
    return nil
}

// UpdateFieldsSortByTableID 批量按表ID更新字段排序（排序值按传入顺序从1开始）
func (v *VOMetadataStorage) UpdateFieldsSortByTableID(tableID int64, fieldIDs []int64) error {
    for idx, fid := range fieldIDs {
        sort := idx + 1
        // 仅针对该表的字段更新排序，避免误更新其他表同ID（理论上不会有）
        // 这里复用按ID更新，懒创建时会设置正确的 TableID
        if err := v.UpdateFieldSortByID(fid, sort); err != nil {
            return err
        }
    }
    return nil
}

// GetTableRemarkByID 根据表ID获取备注
func (v *VOMetadataStorage) GetTableRemarkByID(tableID int64) (string, error) {
    var t VOTableInfo
    if err := v.db.Select("remark").Where("id = ?", tableID).First(&t).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return "", nil
        }
        return "", err
    }
    return t.Remark, nil
}

// （已移除）通过名称组合解析表ID的方法，前端与API均改为ID直连

// UpdateFieldRemarkByTableIDName 根据表ID与字段名更新字段备注
func (v *VOMetadataStorage) UpdateFieldRemarkByTableIDName(tableID int64, fieldName string, remark string) error {
    // 先在 Raw 中定位字段ID
    var rf RawFieldInfo
    if err := v.db.Where("table_id = ? AND name = ?", tableID, fieldName).First(&rf).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil
        }
        return err
    }

    // 用字段ID更新 VO 备注（不做 JOIN，避免 SQLite UPDATE 引用不存在的列）
    res := v.db.Model(&VOFieldInfo{}).Where("id = ?", rf.ID).Update("remark", remark)
    if res.Error != nil {
        return res.Error
    }

    // 若 VO 记录不存在则按 Raw 信息懒创建
    if res.RowsAffected == 0 {
        vo := VOFieldInfo{
            ID:      rf.ID,
            TableID: tableID,
            Remark:  remark,
        }
        if err := v.db.Create(&vo).Error; err != nil {
            return err
        }
    }
    return nil
}

// GetDisplayVOByConfigID 根据配置ID获取显示配置VO
func (v *VOMetadataStorage) GetDisplayVOByConfigID(configID int64) (*models.Display, error) {
    var voDisplay VODisplay
    err := v.db.Where("config_id = ?", configID).First(&voDisplay).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            // 返回默认显示配置
            return &models.Display{
                DBCnt: 0,
                Style: make(map[string]models.Style),
            }, nil
        } 
        return nil, err
    }

    display := &models.Display{
        DBCnt: int(voDisplay.DBCnt),
        Style: make(map[string]models.Style),
    }

    // 已移除 VOStyle，样式返回空映射
    return display, nil
}