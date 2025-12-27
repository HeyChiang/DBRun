package api

import (
	sqlite "dbrun/app/sqlite"
	"dbrun/app/service"
)

// SQLiteAPI 结构体用于暴露给前端的SQLite相关方法
type SQLiteAPI struct{}

// NewSQLiteAPI 创建一个新的SQLiteAPI实例
func NewSQLiteAPI() *SQLiteAPI {
	return &SQLiteAPI{}
}

// Init 初始化SQLite数据库
func (api *SQLiteAPI) Init() error {
	return sqlite.InitSQLite()
}

// InsertCredentials 添加数据库连接凭证，由于wails无法接受指针数据，所以这里要显性的返回结构
func (api *SQLiteAPI) InsertCredentials(creds *sqlite.Credentials) (sqlite.Credentials, error) {
	return service.InsertCredentials(creds)
}

// GetAllCredentials 获取所有数据库连接凭证
func (api *SQLiteAPI) GetAllCredentials() ([]sqlite.Credentials, error) {
	return service.GetAllCredentials()
}

// UpdateCredentials 更新数据库连接凭证
func (api *SQLiteAPI) UpdateCredentials(creds *sqlite.Credentials) error {
	return service.UpdateCredentials(creds)
}

// DeleteCredentials deletes a record from the table_credentials table by ID
func (api *SQLiteAPI) DeleteCredentials(id int) error {
	return service.DeleteCredentialsByID(id)
}

// InsertProject 添加项目，并返回新项目
func (api *SQLiteAPI) InsertProject(name, path string) (*sqlite.Project, error) {
	proj := sqlite.NewProject(name, path)
	return sqlite.InsertProject(proj)
}

// GetAllProjects 获取所有项目
func (api *SQLiteAPI) GetAllProjects() ([]sqlite.Project, error) {
	return sqlite.GetAllProjects()
}

// UpdateProject 更新项目
func (api *SQLiteAPI) UpdateProject(proj *sqlite.Project) error {
	return sqlite.UpdateProject(proj)
}

// DeleteProject 删除项目
func (api *SQLiteAPI) DeleteProject(id int64) error {
	return sqlite.DeleteProjectByID(id)
}
