package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	uploadDir   = "uploads"
	maxFileSize = 2 << 20 // 2MB
)

func init() {
	// 确保上传目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		fmt.Printf("创建上传目录失败: %v\n", err)
	}
}

// UploadImage 处理图片上传
func UploadImage(c *gin.Context) {
	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取上传文件失败"})
		return
	}
	defer file.Close()

	// 检查文件大小
	if c.Request.ContentLength > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小超过限制"})
		return
	}

	// 生成唯一的文件名
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filepath := filepath.Join(uploadDir, filename)

	// 创建目标文件
	out, err := os.Create(filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文件失败"})
		return
	}
	defer out.Close()

	// 复制文件内容
	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}

	// 返回文件URL
	fileURL := fmt.Sprintf("/uploads/%s", filename)
	c.JSON(http.StatusOK, gin.H{
		"url": fileURL,
	})
}
