package metadata

import (
    "time"
)

// Credentials 凭证模型（使用 GORM，表名保持为 table_credentials）
type Credentials struct {
    ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
    Type      string    `json:"type"`
    Label     string    `json:"label"`
    Username  string    `json:"username"`
    Password  string    `json:"password"`
    Host      string    `json:"host"`
    Port      int       `json:"port"`
    Database  string    `json:"database"`
    Instance  string    `json:"instance"`
    Options   string    `json:"options"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (Credentials) TableName() string { return "table_credentials" }
