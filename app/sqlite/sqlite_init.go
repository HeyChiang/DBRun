package sqlite

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB
var currentDBPath string

func InitSQLite() error {
	return InitSQLiteWithPath("./data")
}

func InitSQLiteWithPath(dbDir string) error {
	// 确保数据目录存在
	log.Println("Creating data directory:", dbDir)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Println("InitSQLite: Error creating data directory:", err)
		return err
	}

	dbPath := filepath.Join(dbDir, "app.db")
	
	// 如果数据库路径没有变化，不需要重新初始化
	if currentDBPath == dbPath && DB != nil {
		log.Println("Database path unchanged, skipping reinitialization:", dbPath)
		return nil
	}
	
	// 如果已有数据库连接，先关闭
	if DB != nil {
		log.Println("Closing existing database connection")
		DB.Close()
	}
	
	log.Println("Opening database at:", dbPath)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("InitSQLite: Error opening database:", err)
		return err
	}

	// 测试连接
	log.Println("Pinging database...")
	if err = db.Ping(); err != nil {
		log.Println("InitSQLite: Error pinging database:", err)
		return err
	}

	DB = db
	currentDBPath = dbPath

	// 创建项目表
	log.Println("Creating project table...")
	if err := CreateProjectTable(); err != nil {
		log.Println("InitSQLite: Error creating project table:", err)
		return err
	}

	log.Println("SQLite initialized successfully at:", dbPath)
	return nil
}
