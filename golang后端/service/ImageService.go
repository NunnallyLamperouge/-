package service

import (
	"bytes"
	"encoding/base64"
	"errors"
	"example/config"
	"example/entity"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"gorm.io/gorm"
)

var uploadConfig = &struct {
	maxFileSize  int64
	allowedTypes []string
	imageQuality int
	uploadDir    string
}{
	maxFileSize:  8, // 8MB
	allowedTypes: []string{".jpg", ".jpeg", ".png", ".gif"},
	imageQuality: 75,
	uploadDir:    "uploads",
}

type ImageUpload struct {
	gorm.Model
	ID     int    `gorm:"primaryKey;autoIncrement;bigint" json:"id"`
	UserId int    `gorm:"bigint" json:"user_id"`
	Path   string `gorm:"varchar(100)" json:"path"`
	Size   int    `gorm:"bigint" json:"size"`
}

func validateFileType(filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowedType := range uploadConfig.allowedTypes {
		if ext == allowedType {
			return nil
		}
	}
	return fmt.Errorf("不支持的文件类型: %s，仅支持: %s", ext, strings.Join(uploadConfig.allowedTypes, ", "))
}

func validateFileSize(file *multipart.FileHeader) error {
	maxSize := uploadConfig.maxFileSize << 20
	if file.Size > maxSize {
		return fmt.Errorf("文件大小超过限制 %dMB", uploadConfig.maxFileSize)
	}
	return nil
}

func compressImage(file multipart.File) ([]byte, error) {
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("解码图片失败: %v", err)
	}

	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	if width > 1920 || height > 1080 {
		ratio := float64(width) / float64(height)
		if width > height {
			width = 1920
			height = int(float64(width) / ratio)
		} else {
			height = 1080
			width = int(float64(height) * ratio)
		}
		img = resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	}

	buf := new(bytes.Buffer)
	switch format {
	case "jpeg":
		err = jpeg.Encode(buf, img, &jpeg.Options{Quality: uploadConfig.imageQuality})
	case "png":
		err = png.Encode(buf, img)
	default:
		return nil, errors.New("不支持的图片格式")
	}

	if err != nil {
		return nil, fmt.Errorf("压缩图片失败: %v", err)
	}

	return buf.Bytes(), nil
}

func UploadImage(c *gin.Context) {
	userIdUnFmt, existed := c.Get("userId")
	if !existed {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse[string](401, "校验用户码时发生错误，请重新登陆后尝试。"))
		return
	}
	userId := userIdUnFmt.(int)

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, entity.ErrorResponse[string](400, "未上传文件"))
		return
	}

	if err := validateFileType(file.Filename); err != nil {
		c.JSON(400, entity.ErrorResponse[string](400, err.Error()))

		return
	}

	if err := validateFileSize(file); err != nil {
		c.JSON(400, entity.ErrorResponse[string](400, err.Error()))
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse[string](500, "打开上传图片失败"))

		return
	}
	defer src.Close()

	compressedData, err := compressImage(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse[string](500, "压缩图片失败"))
		return
	}

	ext := filepath.Ext(file.Filename)
	newFilename := fmt.Sprintf("%s_%s%s",
		time.Now().Format("20060102150405"),
		uuid.New().String()[:8],
		ext,
	)

	if err := ensureDir(uploadConfig.uploadDir); err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse[string](500, "创建上传图片文件夹时发生错误"))
		return
	}

	dst := filepath.Join(uploadConfig.uploadDir, newFilename)
	if err := os.WriteFile(dst, compressedData, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse[string](500, "保存文件时发生错误"))
		return
	}

	imageUpload := ImageUpload{
		UserId: userId,
		Path:   dst,
		Size:   int(len(compressedData)),
	}

	tx := config.MysqlDataBase.Begin()
	if err := tx.Create(&imageUpload).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse[string](500, "存储数据库信息时发生错误，请稍后重试"))
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, entity.SuccessResponseWithMessage[ImageUpload]("文件上传成功！", imageUpload))
}
func ensureDir(dirName string) error {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		return os.MkdirAll(dirName, 0755)
	}
	return nil
}

// GetImageBase64 接收文件名，读取文件并转换为 Base64 编码
func GetImageBase64(c *gin.Context) {
	fileName := c.Query("filename")

	// 检查文件是否存在
	_, err := os.Stat("./uploads/" + fileName)
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "图片资源不存在"))
		return
	}

	// 读取文件内容
	fileContent, err := ioutil.ReadFile("./uploads/" + fileName)
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "无法读取文件内容"))
		return
	}

	// 将文件内容转换为 Base64 编码
	base64Content := base64.StdEncoding.EncodeToString(fileContent)

	// 构造 data URI 格式
	dataURI := "data:image/webp;base64," + base64Content

	c.JSON(http.StatusOK, entity.SuccessResponse(dataURI))
}
