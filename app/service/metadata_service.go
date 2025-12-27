package service

import (
	"dbrun/app/connect"
	"dbrun/app/models"
	meta "dbrun/app/sqlite/metadata"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"vitess.io/vitess/go/vt/sqlparser"
)

// MetadataService 将原始与VO存储整合到 service 层，替代 sqlite/metadata/manager.go
type MetadataService struct {
	db         *gorm.DB
	rawStorage *meta.RawMetadataStorage
	voStorage  *meta.VOMetadataStorage
}

// NewMetadataService 创建 service 层的元数据管理器
func NewMetadataService(db *gorm.DB) *MetadataService {
	return &MetadataService{
		db:         db,
		rawStorage: meta.NewRawMetadataStorage(db),
		voStorage:  meta.NewVOMetadataStorage(db),
	}
}

// 全局单例与访问器，兼容原有 getMgr 调用
var (
	msvc   *MetadataService
	mMutex sync.RWMutex
)

// InitMetadataStorageService 初始化元数据存储（保持 API 兼容）
func InitMetadataStorageService(cacheDir string) error {
	mMutex.Lock()
	defer mMutex.Unlock()

	if !filepath.IsAbs(cacheDir) {
		cacheDir = filepath.Join(".", cacheDir)
	}
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create metadata directory %s: %w", cacheDir, err)
	}
	dbPath := filepath.Join(cacheDir, "relation.db")

	gdb, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{Logger: logger.Default.LogMode(logger.Error)})
	if err != nil {
		return fmt.Errorf("failed to open metadata database: %w", err)
	}
	fmt.Printf("[MetadataStorage] init db at %s\n", dbPath)

	_ = gdb.Exec("PRAGMA journal_mode=WAL;")
	_ = gdb.Exec("PRAGMA foreign_keys = ON;")

	service := NewMetadataService(gdb)
	if err := service.InitTables(); err != nil {
		return fmt.Errorf("failed to initialize metadata tables: %w", err)
	}
	fmt.Printf("[MetadataStorage] tables initialized successfully at %s\n", dbPath)

	// 保险：检查关键列是否存在，必要时再迁移一次
	mig := gdb.Migrator()
	if !mig.HasColumn(&meta.VOTableInfo{}, "remark") {
		_ = gdb.AutoMigrate(&meta.VOTableInfo{})
	}
	if !mig.HasColumn(&meta.VOFieldInfo{}, "remark") {
		_ = gdb.AutoMigrate(&meta.VOFieldInfo{})
	}

	msvc = service
	return nil
}

// getMgr 保持原有命名，返回合并后的 MetadataService
func getMgr() (*MetadataService, error) {
	mMutex.RLock()
	defer mMutex.RUnlock()
	if msvc == nil {
		return nil, fmt.Errorf("metadata service not initialized")
	}
	return msvc, nil
}

// InitTables 初始化所有表结构（原始、VO、凭证）
func (m *MetadataService) InitTables() error {
	// 自动迁移原始数据表
	if err := m.db.AutoMigrate(
		&meta.RawDatabaseInfo{},
		&meta.RawSchemaInfo{},
		&meta.RawTableInfo{},
		&meta.RawFieldInfo{},
		&meta.RawViewInfo{},
	); err != nil {
		return err
	}

	// 自动迁移VO数据表
	if err := m.db.AutoMigrate(
		&meta.VODatabaseInfo{},
		&meta.VOSchemaInfo{},
		&meta.VOTableInfo{},
		&meta.VOFieldInfo{},
		&meta.VOViewInfo{},
		&meta.VODisplay{},
	); err != nil {
		return err
	}

	// 自动迁移凭证表（GORM）
	if err := m.db.AutoMigrate(&meta.Credentials{}); err != nil {
		return err
	}
	return nil
}

// UpdateRawDatabase 更新原始数据库信息
func (m *MetadataService) UpdateRawDatabase(configID int64, database connect.DatabaseInfo) error {
	_, err := m.rawStorage.SaveDatabaseInfo(configID, database)
	return err
}

// GetRawStorage 获取原始存储实例
func (m *MetadataService) GetRawStorage() *meta.RawMetadataStorage { return m.rawStorage }

// GetVOStorage 获取VO存储实例
func (m *MetadataService) GetVOStorage() *meta.VOMetadataStorage { return m.voStorage }

// ConvertRawToVO 将原始数据转换为VO数据（保留原始ID）
func (m *MetadataService) ConvertRawToVO(configID int64) (*models.DBInfoVO, error) {
	rawDatabases, err := m.rawStorage.GetDatabasesByConfigID(configID)
	if err != nil {
		return nil, err
	}
	var dbInfoVO models.DBInfoVO
	dbInfoVO.ConvertToVO(rawDatabases)

	display, err := m.voStorage.GetDisplayVOByConfigID(configID)
	if err != nil {
		return nil, err
	}
	dbInfoVO.Display = *display

	rawDBRows, err := m.rawStorage.GetRawDatabasesRows(configID)
	if err != nil {
		return nil, err
	}
	dbByName := make(map[string]meta.RawDatabaseInfo)
	for _, rd := range rawDBRows {
		dbByName[rd.Name] = rd
	}
	for i := range dbInfoVO.DBs {
		dbVO := &dbInfoVO.DBs[i]
		if rd, ok := dbByName[dbVO.Name]; ok {
			dbVO.ID = rd.ID
			dbVO.ConfigID = rd.ConfigID

			rawSchemas, err := m.rawStorage.GetRawSchemasRows(rd.ID)
			if err != nil {
				return nil, err
			}
			schemaByName := make(map[string]meta.RawSchemaInfo)
			for _, rs := range rawSchemas {
				schemaByName[rs.Name] = rs
			}
			for si := range dbVO.Schemas {
				sv := &dbVO.Schemas[si]
				if rs, ok := schemaByName[sv.Name]; ok {
					sv.ID = rs.ID
					sv.DatabaseID = rs.DatabaseID
				}
			}

			rawTablesNoSchema, err := m.rawStorage.GetRawTablesRows(rd.ID, nil)
			if err != nil {
				return nil, err
			}
			tableByName := make(map[string]meta.RawTableInfo)
			for _, rt := range rawTablesNoSchema {
				tableByName[rt.Name] = rt
			}
			for ti := range dbVO.Tables {
				tv := &dbVO.Tables[ti]
				if rt, ok := tableByName[tv.Name]; ok {
					tv.ID = rt.ID
					tv.DatabaseID = rt.DatabaseID
					tv.SchemaID = nil
					rawFields, err := m.rawStorage.GetRawFieldsRows(rt.ID)
					if err != nil {
						return nil, err
					}
					fByName := make(map[string]meta.RawFieldInfo)
					for _, rf := range rawFields {
						fByName[rf.Name] = rf
					}
					for fi := range tv.Fields {
						fv := &tv.Fields[fi]
						if rf, ok := fByName[fv.Name]; ok {
							fv.ID = rf.ID
							fv.TableID = rf.TableID
						}
					}
				}
			}

			rawViewsNoSchema, err := m.rawStorage.GetRawViewsRows(rd.ID, nil)
			if err != nil {
				return nil, err
			}
			viewByName := make(map[string]meta.RawViewInfo)
			for _, rv := range rawViewsNoSchema {
				viewByName[rv.Name] = rv
			}
			for vi := range dbVO.Views {
				vv := &dbVO.Views[vi]
				if rv, ok := viewByName[vv.Name]; ok {
					vv.ID = rv.ID
					vv.DatabaseID = rv.DatabaseID
					vv.SchemaID = nil
				}
			}

			for _, rs := range rawSchemas {
				rawTablesInSchema, err := m.rawStorage.GetRawTablesRows(rd.ID, &rs.ID)
				if err != nil {
					return nil, err
				}
				tblByName := make(map[string]meta.RawTableInfo)
				for _, rt := range rawTablesInSchema {
					tblByName[rt.Name] = rt
				}
				var schemaVO *models.SchemaVO
				for si := range dbVO.Schemas {
					if dbVO.Schemas[si].Name == rs.Name {
						schemaVO = &dbVO.Schemas[si]
						break
					}
				}
				if schemaVO != nil {
					for ti := range schemaVO.Tables {
						tv := &schemaVO.Tables[ti]
						if rt, ok := tblByName[tv.Name]; ok {
							tv.ID = rt.ID
							tv.DatabaseID = rt.DatabaseID
							sid := rs.ID
							tv.SchemaID = &sid
							rawFields, err := m.rawStorage.GetRawFieldsRows(rt.ID)
							if err != nil {
								return nil, err
							}
							fByName := make(map[string]meta.RawFieldInfo)
							for _, rf := range rawFields {
								fByName[rf.Name] = rf
							}
							for fi := range tv.Fields {
								fv := &tv.Fields[fi]
								if rf, ok := fByName[fv.Name]; ok {
									fv.ID = rf.ID
									fv.TableID = rf.TableID
								}
							}
						}
					}
					rawViewsInSchema, err := m.rawStorage.GetRawViewsRows(rs.DatabaseID, &rs.ID)
					if err != nil {
						return nil, err
					}
					vwByName := make(map[string]meta.RawViewInfo)
					for _, rv := range rawViewsInSchema {
						vwByName[rv.Name] = rv
					}
					for vi := range schemaVO.Views {
						vv := &schemaVO.Views[vi]
						if rv, ok := vwByName[vv.Name]; ok {
							vv.ID = rv.ID
							vv.DatabaseID = rv.DatabaseID
							sid2 := rs.ID
							vv.SchemaID = &sid2
						}
					}
				}
			}
		}
	}
	return &dbInfoVO, nil
}

// ===== 凭证操作（迁移自 sqlite/metadata/credentials.go） =====

func (m *MetadataService) InsertCredentials(creds *meta.Credentials) error {
	return m.db.Create(creds).Error
}

func (m *MetadataService) GetAllCredentials() ([]meta.Credentials, error) {
	var list []meta.Credentials
	err := m.db.Order("id ASC").Find(&list).Error
	return list, err
}

func (m *MetadataService) UpdateCredentials(creds *meta.Credentials) error {
	return m.db.Model(&meta.Credentials{}).
		Where("id = ?", creds.ID).
		Updates(map[string]interface{}{
			"type":     creds.Type,
			"label":    creds.Label,
			"username": creds.Username,
			"password": creds.Password,
			"host":     creds.Host,
			"port":     creds.Port,
			"database": creds.Database,
			"instance": creds.Instance,
			"options":  creds.Options,
		}).Error
}

func (m *MetadataService) DeleteCredentialsByID(id int) error {
    // 先清理与该配置ID关联的所有元数据（Raw 与 VO）
    if err := m.voStorage.DeleteVOByConfigID(int64(id)); err != nil {
        return err
    }
    if err := m.rawStorage.DeleteByConfigID(int64(id)); err != nil {
        return err
    }
    // 最后删除凭证记录
    return m.db.Where("id = ?", id).Delete(&meta.Credentials{}).Error
}

func (m *MetadataService) GetCredentialsByID(id int64) (*meta.Credentials, error) {
	var c meta.Credentials
	if err := m.db.Where("id = ?", id).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

// ===== 以下为合并自 metadata_operation.go 的 API 函数与辅助方法 =====

// ListDatabasesByConfig 获取数据库信息（优先使用原始存储；VO仅用于叠加业务配置）
func ListDatabasesByConfig(config connect.Config) (models.DBInfoVO, error) {
	manager, err := getMgr()
	if err != nil {
		return models.DBInfoVO{}, err
	}

	// 1) 先尝试读取原始元数据
	rawStored, err := manager.rawStorage.GetDatabasesByConfigID(int64(config.ID))
	if err == nil && len(rawStored) > 0 {
		baseVO, err := manager.ConvertRawToVO(int64(config.ID))
		if err != nil {
			return models.DBInfoVO{}, err
		}
		overlayVO, _ := manager.voStorage.GetDatabasesVOByConfigID(int64(config.ID))
		if overlayVO != nil {
			merged := mergeVOOverlay(*baseVO, models.DBInfoVO{DBs: overlayVO, Display: baseVO.Display})
			return merged, nil
		}
		return *baseVO, nil
	}

	// 2) 若原始存储没有，则从数据库拉取并仅保存原始数据，并初始化VO以提供稳定ID
	conn, err := connect.GetConnection(config)
	if err != nil {
		return models.DBInfoVO{}, fmt.Errorf("failed to get connection: %w", err)
	}
	rawFetched, err := fetchRawDatabases(conn)
	if err != nil {
		return models.DBInfoVO{}, err
	}
	rs := manager.rawStorage
	for _, db := range rawFetched {
		if _, err := rs.SaveDatabaseInfo(int64(config.ID), db); err != nil {
			return models.DBInfoVO{}, err
		}
	}
	baseVO, err := manager.ConvertRawToVO(int64(config.ID))
	if err != nil {
		return models.DBInfoVO{}, err
	}
	overlayVO, _ := manager.voStorage.GetDatabasesVOByConfigID(int64(config.ID))
	if overlayVO != nil {
		merged := mergeVOOverlay(*baseVO, models.DBInfoVO{DBs: overlayVO, Display: baseVO.Display})
		return merged, nil
	}
	return *baseVO, nil
}

// GetFieldsVOByTableID 根据表ID获取字段VO（仅VO，不触发全量加载）
func GetFieldsVOByTableID(tableID int64) ([]models.FieldInfoVO, error) {
	manager, err := getMgr()
	if err != nil {
		return nil, err
	}
	vs := manager.voStorage
	rs := manager.rawStorage

	voFields, err := vs.GetFieldsVOByTableID(tableID)
	if err != nil {
		return nil, err
	}
	rawRows, err := rs.GetRawFieldsRows(tableID)
	if err != nil {
		return nil, err
	}
	rawByID := make(map[int64]meta.RawFieldInfo, len(rawRows))
	for _, rf := range rawRows {
		rawByID[rf.ID] = rf
	}
	merged := make([]models.FieldInfoVO, 0, len(voFields))
	for _, vf := range voFields {
		mf := vf
		if rf, ok := rawByID[vf.ID]; ok {
			mf.Type = rf.Type
			mf.Nullable = rf.Nullable
			mf.Key = rf.Key
			mf.Comment = rf.Comment
			mf.DefaultValue = rf.DefaultValue
		} else {
			for _, r := range rawRows {
				if r.Name == vf.Name {
					mf.Type = r.Type
					mf.Nullable = r.Nullable
					mf.Key = r.Key
					mf.Comment = r.Comment
					mf.DefaultValue = r.DefaultValue
					break
				}
			}
		}
		merged = append(merged, mf)
	}
	return merged, nil
}

// GetTablesVOByDatabaseID 根据数据库ID与可选的SchemaID获取表VO（仅VO，不触发全量加载）
func GetTablesVOByDatabaseID(databaseID int64, schemaID *int64) ([]models.TableInfoVO, error) {
    manager, err := getMgr()
    if err != nil {
        return nil, err
    }
    vs := manager.voStorage
    return vs.GetTablesVOByDatabaseID(databaseID, schemaID)
}

// UpdateFieldsSortByTableID 根据表ID批量更新字段排序（排序值按数组顺序从1开始）
func UpdateFieldsSortByTableID(tableID int64, fieldIDs []int64) error {
    manager, err := getMgr()
    if err != nil {
        return err
    }
    vs := manager.voStorage
    return vs.UpdateFieldsSortByTableID(tableID, fieldIDs)
}

// TestConnection 测试数据库连接
func TestConnection(config connect.Config) error {
	conn, err := connect.GetConnection(config)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}
	return conn.Test()
}

// CloseAllConnections 关闭所有连接
// 移除冗余包装：直接在 API 层调用 connect.CloseAllConnections

// VO 扩展信息（备注/显示）
type TableCacheVO struct {
	Remark   string            `json:"remark"`
	Version  int64             `json:"version"`
	FieldMap map[string]string `json:"fieldMap"`
}

// GetTableVOCacheByTableID 基于表ID直接获取VO缓存（备注/字段备注）
func GetTableVOCacheByTableID(tableID int64) TableCacheVO {
	manager, err := getMgr()
	if err != nil {
		return TableCacheVO{}
	}
	vs := manager.voStorage
	remark, _ := vs.GetTableRemarkByID(tableID)
	fields, _ := vs.GetFieldsVOByTableID(tableID)
	fmap := make(map[string]string)
	for _, f := range fields {
		if f.Remark != "" {
			fmap[f.Name] = f.Remark
		}
	}
	return TableCacheVO{Remark: remark, Version: 0, FieldMap: fmap}
}

// ParseViewSQL 解析视图SQL，获取字段列表
func ParseViewSQL(sql string) ([]models.FieldInfoVO, error) {
	parser := sqlparser.NewTestParser()
	stmt, err := parser.Parse(sql)
	if err != nil {
		return nil, fmt.Errorf("解析SQL失败: %w", err)
	}
	selectStmt, ok := stmt.(*sqlparser.Select)
	if !ok {
		return nil, fmt.Errorf("无法解析非SELECT语句，当前语句类型: %T", stmt)
	}
	var fields []models.FieldInfoVO
	for _, se := range selectStmt.SelectExprs {
		switch expr := se.(type) {
		case *sqlparser.AliasedExpr:
			var fieldName string
			if expr.As.IsEmpty() {
				if cn, ok := expr.Expr.(*sqlparser.ColName); ok {
					fieldName = cn.Name.String()
				} else {
					fieldName = sqlparser.String(expr.Expr)
				}
			} else {
				fieldName = expr.As.String()
			}
			fields = append(fields, models.FieldInfoVO{Name: fieldName, Type: "VARCHAR", Display: true, Nullable: true, Key: ""})
		case *sqlparser.StarExpr:
			fields = append(fields, models.FieldInfoVO{Name: "*", Type: "UNKNOWN", Display: true, Nullable: true, Key: "", Comment: "星号表达式需要访问数据库才能获取具体字段"})
		default:
			fields = append(fields, models.FieldInfoVO{Name: sqlparser.String(expr), Type: "UNKNOWN", Display: true, Nullable: true, Key: ""})
		}
	}
	return fields, nil
}

// mergeVOOverlay 将VO中的业务配置（备注/显示等）叠加到原始转换后的VO
func mergeVOOverlay(base models.DBInfoVO, overlay models.DBInfoVO) models.DBInfoVO {
	oDB := make(map[string]models.DatabaseInfoVO)
	for _, db := range overlay.DBs {
		oDB[db.Name] = db
	}
	for i := range base.DBs {
		b := &base.DBs[i]
		odb, ok := oDB[b.Name]
		if !ok {
			continue
		}
		b.ID = odb.ID
		b.ConfigID = odb.ConfigID
		if b.Alias == "" {
			b.Alias = odb.Alias
		}
		mergeTablesOverlay(&b.Tables, odb.Tables)
		mergeViewsOverlay(&b.Views, odb.Views)
		oSchemas := make(map[string]models.SchemaVO)
		for _, s := range odb.Schemas {
			oSchemas[s.Name] = s
		}
		for si := range b.Schemas {
			bs := &b.Schemas[si]
			if os, ok := oSchemas[bs.Name]; ok {
				bs.ID = os.ID
				bs.DatabaseID = os.DatabaseID
				if bs.Alias == "" {
					bs.Alias = os.Alias
				}
				mergeTablesOverlay(&bs.Tables, os.Tables)
				mergeViewsOverlay(&bs.Views, os.Views)
			}
		}
	}
	if overlay.Display.Style != nil && len(overlay.Display.Style) > 0 {
		base.Display = overlay.Display
	}
	return base
}

func mergeTablesOverlay(base *[]models.TableInfoVO, over []models.TableInfoVO) {
    om := make(map[string]models.TableInfoVO)
    for _, t := range over {
        om[t.Name] = t
    }
    for i := range *base {
        bt := &(*base)[i]
        if ot, ok := om[bt.Name]; ok {
            bt.ID = ot.ID
            bt.DatabaseID = ot.DatabaseID
            bt.SchemaID = ot.SchemaID
            if bt.Alias == "" {
                bt.Alias = ot.Alias
            }
            if ot.Color != "" {
                bt.Color = ot.Color
            }
            if ot.Remark != "" {
                bt.Remark = ot.Remark
            }
            // 构建叠加字段映射（按名称）
            fm := make(map[string]models.FieldInfoVO)
            for _, f := range ot.Fields {
                fm[f.Name] = f
            }
            // 基表字段的快速索引，便于重排与属性同步
            bm := make(map[string]models.FieldInfoVO)
            for _, bf := range bt.Fields {
                bm[bf.Name] = bf
            }
            // 按 VO 字段顺序重建字段列表，并同步 Display/Remark/Sort
            newFields := make([]models.FieldInfoVO, 0, len(bt.Fields))
            for _, of := range ot.Fields {
                if bf, ok := bm[of.Name]; ok {
                    // 同步属性
                    bf.Display = of.Display
                    if of.Remark != "" {
                        bf.Remark = of.Remark
                    }
                    bf.Sort = of.Sort
                    newFields = append(newFields, bf)
                    delete(bm, of.Name)
                } else {
                    // VO 中存在但原始字段中不存在的条目，直接附加（保留VO属性）
                    newFields = append(newFields, of)
                }
            }
            // 追加未出现在 VO 中的剩余字段，保持其原始相对顺序
            for _, bf := range bt.Fields {
                if _, ok := fm[bf.Name]; !ok {
                    newFields = append(newFields, bf)
                }
            }
            bt.Fields = newFields
        }
    }
}

func mergeViewsOverlay(base *[]models.ViewInfoVO, over []models.ViewInfoVO) {
	om := make(map[string]models.ViewInfoVO)
	for _, v := range over {
		om[v.Name] = v
	}
	for i := range *base {
		bv := &(*base)[i]
		if ov, ok := om[bv.Name]; ok {
			bv.ID = ov.ID
			bv.DatabaseID = ov.DatabaseID
			bv.SchemaID = ov.SchemaID
			if bv.Alias == "" {
				bv.Alias = ov.Alias
			}
			if ov.Color != "" {
				bv.Color = ov.Color
			}
			if ov.Remark != "" {
				bv.Remark = ov.Remark
			}
		}
	}
}

// 原始数据拉取
func fetchRawDatabases(conn connect.Connection) ([]connect.DatabaseInfo, error) {
	dbNames, err := conn.GetDBNames()
	if err != nil {
		return nil, fmt.Errorf("failed to get database names: %w", err)
	}
	var raws []connect.DatabaseInfo
	for _, db := range dbNames {
		r, err := fetchRawDatabase(conn, db.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to get raw info for database %s: %w", db.Name, err)
		}
		raws = append(raws, r)
	}
	return raws, nil
}

func fetchRawDatabase(conn connect.Connection, dbName string) (connect.DatabaseInfo, error) {
	info := connect.DatabaseInfo{Name: dbName}
	dbType := conn.GetConfig().Type
    if dbType == "postgresql" || dbType == "sqlserver" || dbType == "oracle" {
		schemas, err := conn.GetSchemas(dbName)
		if err != nil {
			return info, fmt.Errorf("failed to get schemas for database %s: %w", dbName, err)
		}
		for _, s := range schemas {
			sc := connect.Schema{Name: s.Name}
			tables, err := conn.GetTables(connect.QueryParams{Database: dbName, Schema: s.Name})
			if err != nil {
				return info, fmt.Errorf("failed to get tables for schema %s in %s: %w", s.Name, dbName, err)
			}
			for i := range tables {
				fields, err := conn.GetTableFields(connect.QueryParams{Database: dbName, Schema: s.Name, Table: tables[i].Name})
				if err != nil {
					return info, fmt.Errorf("failed to get fields for table %s: %w", tables[i].Name, err)
				}
				tables[i].Fields = fields
			}
			views, err := conn.GetViews(connect.QueryParams{Database: dbName, Schema: s.Name})
			if err != nil {
				return info, fmt.Errorf("failed to get views for schema %s in %s: %w", s.Name, dbName, err)
			}
			sc.Tables = tables
			sc.Views = views
			info.Schemas = append(info.Schemas, sc)
		}
	} else {
		tables, err := conn.GetTables(connect.QueryParams{Database: dbName})
		if err != nil {
			return info, fmt.Errorf("failed to get tables for database %s: %w", dbName, err)
		}
		for i := range tables {
			fields, err := conn.GetTableFields(connect.QueryParams{Database: dbName, Table: tables[i].Name})
			if err != nil {
				return info, fmt.Errorf("failed to get fields for table %s: %w", tables[i].Name, err)
			}
			tables[i].Fields = fields
		}
		views, err := conn.GetViews(connect.QueryParams{Database: dbName})
		if err != nil {
			return info, fmt.Errorf("failed to get views for database %s: %w", dbName, err)
		}
		info.Tables = tables
		info.Views = views
	}
	return info, nil
}
