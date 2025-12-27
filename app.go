package main

import (
    "dbrun/app/api"
    "dbrun/app/cache"
    "dbrun/app/sqlite"
    "context"
    "fmt"
)

// App struct
type App struct {
	ctx              context.Context
	sqliteAPI        *api.SQLiteAPI
	metadatasAPI     *api.MetadatasAPI
	appCache         *api.AppCacheApi
	systemServiceAPI *api.SystemServiceAPI
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		sqliteAPI:        api.NewSQLiteAPI(),
		metadatasAPI:     api.NewMetadatasAPI(),
		appCache:         api.NewAppCacheApi(),
		systemServiceAPI: api.NewSystemServiceAPI(),
	}
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
    a.ctx = ctx
    a.systemServiceAPI.Init(ctx)
    a.sqliteAPI.Init()

    // 取消默认目录的缓存初始化：仅在前端选择/创建项目时初始化

    // 元数据存储将按项目路径在前端创建时初始化
}

func (a *App) Shutdown(ctx context.Context) {
	fmt.Println("Shutdown started...")

	fmt.Println("Closing all metadata connections...")
	err := a.metadatasAPI.CloseAllConnections()
	if err != nil {
		fmt.Printf("Error closing metadata connections: %v\n", err)
	}
	fmt.Println("All metadata connections closed successfully")

	fmt.Println("Closing BoltDB cache...")
	if err := cache.Close(); err != nil {
		fmt.Printf("Error closing BoltDB cache: %v\n", err)
	} else {
		fmt.Println("BoltDB cache closed successfully")
	}

	fmt.Println("Closing SQLite database...")
	if sqlite.DB != nil {
		sqlite.DB.Close()
		fmt.Println("SQLite database closed successfully")
	} else {
		fmt.Println("SQLite database is nil, nothing to close")
	}

	fmt.Println("Shutdown completed")
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
