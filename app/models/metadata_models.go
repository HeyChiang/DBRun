package models

import "dbrun/app/connect"

// DatabaseExtendedInfo 数据库扩展信息
type DatabaseExtendedInfo struct {
	Remark string `json:"remark"`
	Color  string `json:"color"`
}

// TableExtendedInfo 表扩展信息
type TableExtendedInfo struct {
	Remark string `json:"remark"`
	Color  string `json:"color"`
}

// FieldExtendedInfo 字段扩展信息
type FieldExtendedInfo struct {
    Remark  string `json:"remark"`
    Display bool   `json:"display"`
}

// ViewExtendedInfo 视图扩展信息
type ViewExtendedInfo struct {
	Remark string `json:"remark"`
	Color  string `json:"color"`
}

// FieldInfoVO 字段信息VO
type FieldInfoVO struct {
    ID           int64  `json:"id"`
    TableID      int64  `json:"tableId"`
    Alias        string `json:"alias"`
    Name         string `json:"name"`
    Type         string `json:"type"`
    Display      bool   `json:"display"`
    Nullable     bool   `json:"nullable,omitempty"`
    Key          string `json:"key"`
    Comment      string `json:"comment,omitempty"`
    DefaultValue string `json:"default_value,omitempty"`
    Remark       string `json:"remark"`
    FontColor    string `json:"fontColor"`
    BgColor      string `json:"bgColor"`
    Sort         int    `json:"sort"`
}

// TableInfoVO 表信息VO
type TableInfoVO struct {
    ID      int64        `json:"id"`
    DatabaseID int64     `json:"databaseId"`
    SchemaID   *int64    `json:"schemaId,omitempty"`
    Alias   string        `json:"alias"`
    Name    string        `json:"name"`
    Color   string        `json:"color"`
    Comment string        `json:"comment"`
    Remark  string        `json:"remark"`
    Fields  []FieldInfoVO `json:"fields"`
}

// ViewInfoVO 视图信息VO
type ViewInfoVO struct {
    ID         int64  `json:"id"`
    DatabaseID int64  `json:"databaseId"`
    SchemaID   *int64 `json:"schemaId,omitempty"`
    Alias      string `json:"alias"`
    Name       string `json:"name"`
    Color      string `json:"color"`
    Definition string `json:"definition"`
    Remark     string `json:"remark"`
}

// SchemaVO Schema信息VO
type SchemaVO struct {
    ID        int64        `json:"id"`
    DatabaseID int64       `json:"databaseId"`
    Alias     string       `json:"alias"`
    Name   string        `json:"name"`
    Tables []TableInfoVO `json:"tables"`
    Views  []ViewInfoVO  `json:"views"`
}

// DatabaseInfoVO 数据库信息VO
type DatabaseInfoVO struct {
    ID      int64         `json:"id"`
    ConfigID int64        `json:"configId"`
    Alias   string        `json:"alias"`
    Name    string        `json:"name"`
    Comment string        `json:"comment"`
    Schemas []SchemaVO    `json:"schemas"`
    Tables  []TableInfoVO `json:"tables"`
    Views   []ViewInfoVO  `json:"views"`
}

// Style 样式配置
type Style struct {
	Color  string `json:"color"`
	IsShow bool   `json:"isShow"`
}

// Display 显示配置
type Display struct {
	DBCnt int               `json:"dbCnt"`
	Style map[string]Style  `json:"style"`
}

// DBInfoVO 数据库信息VO的完整结构
type DBInfoVO struct {
	DBs     []DatabaseInfoVO `json:"dbs"`
	Display Display          `json:"display"`
}

// ConvertToVO 将connect.DatabaseInfo转换为DatabaseInfoVO
func (vo *DBInfoVO) ConvertToVO(dbInfos []connect.DatabaseInfo) {
	vo.DBs = make([]DatabaseInfoVO, len(dbInfos))
	for i, db := range dbInfos {
		vo.DBs[i] = convertDatabaseInfoToVO(db)
	}
}

// 将connect.DatabaseInfo转换为DatabaseInfoVO
func convertDatabaseInfoToVO(db connect.DatabaseInfo) DatabaseInfoVO {
	dbVO := DatabaseInfoVO{
		Name:    db.Name,
		Comment: db.Comment,
		Tables:  make([]TableInfoVO, len(db.Tables)),
		Views:   make([]ViewInfoVO, len(db.Views)),
		Schemas: make([]SchemaVO, len(db.Schemas)),
	}

	// 转换表信息
	for i, table := range db.Tables {
		dbVO.Tables[i] = convertTableInfoToVO(table)
	}

	// 转换视图信息
	for i, view := range db.Views {
		dbVO.Views[i] = convertViewInfoToVO(view)
	}

	// 转换Schema信息
	for i, schema := range db.Schemas {
		dbVO.Schemas[i] = convertSchemaToVO(schema)
	}

	return dbVO
}

// 将connect.Schema转换为SchemaVO
func convertSchemaToVO(schema connect.Schema) SchemaVO {
	schemaVO := SchemaVO{
		Name:   schema.Name,
		Tables: make([]TableInfoVO, len(schema.Tables)),
		Views:  make([]ViewInfoVO, len(schema.Views)),
	}

	// 转换表信息
	for i, table := range schema.Tables {
		schemaVO.Tables[i] = convertTableInfoToVO(table)
	}

	// 转换视图信息
	for i, view := range schema.Views {
		schemaVO.Views[i] = convertViewInfoToVO(view)
	}

	return schemaVO
}

// 将connect.TableInfo转换为TableInfoVO
func convertTableInfoToVO(table connect.TableInfo) TableInfoVO {
	tableVO := TableInfoVO{
		Name:    table.Name,
		Comment: table.Comment,
		Remark:  "", // 默认为空，可以根据需要从其他地方获取
		Fields:  make([]FieldInfoVO, len(table.Fields)),
	}

	// 转换字段信息
	for i, field := range table.Fields {
		tableVO.Fields[i] = convertFieldInfoToVO(field)
	}

	return tableVO
}

// 将connect.FieldInfo转换为FieldInfoVO
func convertFieldInfoToVO(field connect.FieldInfo) FieldInfoVO {
	return FieldInfoVO{
		Name:         field.Name,
		Type:         field.Type,
		Nullable:     field.Nullable,
		Key:          field.Key,
		Display:      true,
		Comment:      field.Comment,
		DefaultValue: field.DefaultValue,
		Remark:       "", // 默认为空，可以根据需要从其他地方获取
	}
}

// 将connect.ViewInfo转换为ViewInfoVO
func convertViewInfoToVO(view connect.ViewInfo) ViewInfoVO {
	return ViewInfoVO{
		Name:       view.Name,
		Definition: view.Definition,
		Remark:     "", // 默认为空，可以根据需要从其他地方获取
	}
}

// TableCacheVO 表缓存VO
type TableCacheVO struct {
	Tables []TableInfoVO `json:"tables"`
	Views  []ViewInfoVO  `json:"views"`
}