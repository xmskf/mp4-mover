package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	sourceDir     = "/xunlei"
	targetDir     = "/tmp"
	scanInterval  = 5 * time.Minute
	fileExtension = ".mp4"
)

func main() {
	log.Println("文件搬运程序已启动...")
	log.Printf("监控目录: %s", sourceDir)
	log.Printf("目标目录: %s", targetDir)
	log.Printf("扫描间隔: %v", scanInterval)

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		log.Fatalf("无法创建或访问目标目录 %s: %v", targetDir, err)
	}

	ticker := time.NewTicker(scanInterval)
	defer ticker.Stop()

	// 启动时先执行一次
	moveFiles()

	for range ticker.C {
		log.Println("--- 开始新一轮扫描 ---")
		moveFiles()
	}
}

func moveFiles() {
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		log.Printf("源目录 %s 不存在，跳过本次扫描", sourceDir)
		return
	}

	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		log.Printf("读取源目录失败: %v", err)
		return
	}

	movedCount := 0

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if strings.EqualFold(filepath.Ext(entry.Name()), fileExtension) {
			srcPath := filepath.Join(sourceDir, entry.Name())
			dstPath := filepath.Join(targetDir, entry.Name())

			// 防止 /tmp 已有同名文件导致覆盖
			if _, err := os.Stat(dstPath); err == nil {
				ext := filepath.Ext(entry.Name())
				baseName := strings.TrimSuffix(entry.Name(), ext)
				newName := fmt.Sprintf("%s_%s%s", baseName, time.Now().Format("20060102150405"), ext)
				dstPath = filepath.Join(targetDir, newName)
			}

			err := moveFile(srcPath, dstPath)
			if err != nil {
				log.Printf("移动文件失败 %s -> %s: %v", srcPath, dstPath, err)
			} else {
				log.Printf("成功移动: %s -> %s", srcPath, dstPath)
				movedCount++
			}
		}
	}

	if movedCount > 0 {
		log.Printf("本次共移动了 %d 个文件。", movedCount)
	} else {
		log.Println("本次未发现需要移动的文件。")
	}
}

// moveFile 处理跨文件系统（跨分区）的移动
func moveFile(src, dst string) error {
	// 先尝试直接重命名（同分区极快）
	err := os.Rename(src, dst)
	if err == nil {
		return nil
	}

	// 如果报错，检查是否是跨设备错误
	if linkErr, ok := err.(*os.LinkError); ok && linkErr.Err.Error() == "invalid cross-device link" {
		// 执行复制+删除
		return copyAndDelete(src, dst)
	}

	return err
}

func copyAndDelete(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	// 确保数据落盘
	out.Sync()
	out.Close()
	in.Close()

	// 删除源文件
	return os.Remove(src)
}
