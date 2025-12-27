package service

import (
    "context"
    "fmt"
    "github.com/wailsapp/wails/v2/pkg/runtime"
    "os"
    "path/filepath"
)

// 打开系统目录选择对话框
func OpenDirectory(ctx context.Context) (string, error) {
	if ctx == nil {
		return "", fmt.Errorf("context is not set")
	}
	result, err := runtime.OpenDirectoryDialog(ctx, runtime.OpenDialogOptions{
		Title: "选择项目文件夹",
	})
	if err != nil {
		return "", err
	}
	return result, nil
}

// CreateDirectory 创建目录，如果目录不存在则创建
func CreateDirectory(dirPath string) error {
    fmt.Printf("CreateDirectory: %s\n", dirPath)

    // 确保路径是绝对路径
    absPath, err := filepath.Abs(dirPath)
    if err != nil {
        return fmt.Errorf("failed to get absolute path for %s: %w", dirPath, err)
    }

    // 检查目录是否已存在
    if _, err := os.Stat(absPath); os.IsNotExist(err) {
        // 创建目录，包括所有必要的父目录
        if err := os.MkdirAll(absPath, 0755); err != nil {
            return fmt.Errorf("failed to create directory %s: %w", absPath, err)
        }
        fmt.Printf("Directory created successfully: %s\n", absPath)
    } else if err != nil {
        return fmt.Errorf("failed to check directory %s: %w", absPath, err)
    } else {
        // 目录已存在则报错，避免重复创建
        return fmt.Errorf("当前项目目录已存在，不能重复创建")
    }
    return nil
}

// PathExists 检查给定路径是否存在且为目录
func PathExists(dirPath string) (bool, error) {
    absPath, err := filepath.Abs(dirPath)
    if err != nil {
        return false, fmt.Errorf("failed to get absolute path for %s: %w", dirPath, err)
    }
    info, err := os.Stat(absPath)
    if os.IsNotExist(err) {
        return false, nil
    }
    if err != nil {
        return false, fmt.Errorf("failed to check path %s: %w", absPath, err)
    }
    return info.IsDir(), nil
}
