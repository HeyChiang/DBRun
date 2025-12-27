package sqlite

import (
	"time"
)

type Project struct {
	ID        int64     `json:"id"`         // 唯一标识，自增
	Name      string    `json:"name"`       // 项目名称
	Path      string    `json:"path"`       // 项目路径
	CreatedAt time.Time `json:"created_at"` // 创建时间
}

// NewProject 创建新的 Project 实例，CreatedAt 默认为当前时间
func NewProject(name, path string) *Project {
	return &Project{
		Name:      name,
		Path:      path,
		CreatedAt: time.Now(),
	}
}

// CreateProjectTable 创建项目表
func CreateProjectTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS table_project (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			path TEXT NOT NULL,
			created_at DATETIME NOT NULL
		)`
	_, err := DB.Exec(query)
	return err
}

// InsertProject 插入新项目，并返回新项目
func InsertProject(proj *Project) (*Project, error) {
	query := `
		INSERT INTO table_project (name, path, created_at)
			VALUES (?, ?, ?)`
	result, err := DB.Exec(query,
		proj.Name,
		proj.Path,
		proj.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	proj.ID = id
	return proj, nil
}

// GetAllProjects 查询所有项目
func GetAllProjects() ([]Project, error) {
	query := `SELECT id, name, path, created_at FROM table_project`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var proj Project
		if err := rows.Scan(&proj.ID, &proj.Name, &proj.Path, &proj.CreatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, proj)
	}
	return projects, nil
}

// UpdateProject 根据ID更新项目信息
func UpdateProject(proj *Project) error {
	query := `
		UPDATE table_project SET name=?, path=? WHERE id=?`
	_, err := DB.Exec(query, proj.Name, proj.Path, proj.ID)
	return err
}

// DeleteProjectByID 根据ID删除项目
func DeleteProjectByID(id int64) error {
	query := `DELETE FROM table_project WHERE id=?`
	_, err := DB.Exec(query, id)
	return err
}
