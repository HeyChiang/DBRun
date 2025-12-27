package service

import (
    meta "dbrun/app/sqlite/metadata"
)

// InsertCredentials 添加数据库连接凭证，返回插入后的结构
func InsertCredentials(creds *meta.Credentials) (meta.Credentials, error) {
    manager, err := getMgr()
    if err != nil {
        return meta.Credentials{}, err
    }
    err = manager.InsertCredentials(creds)
    return *creds, err
}

// GetAllCredentials 获取所有数据库连接凭证
func GetAllCredentials() ([]meta.Credentials, error) {
    manager, err := getMgr()
    if err != nil {
        return nil, err
    }
    return manager.GetAllCredentials()
}

// UpdateCredentials 更新数据库连接凭证
func UpdateCredentials(creds *meta.Credentials) error {
    manager, err := getMgr()
    if err != nil {
        return err
    }
    return manager.UpdateCredentials(creds)
}

// DeleteCredentialsByID 根据ID删除数据库连接凭证
func DeleteCredentialsByID(id int) error {
    manager, err := getMgr()
    if err != nil {
        return err
    }
    return manager.DeleteCredentialsByID(id)
}