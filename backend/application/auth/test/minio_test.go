package test

import (
	"backend/pkg/oss"
	"context"
	"github.com/minio/minio-go/v7"
	"log"
	"net/url"
	"testing"
	"time"
)

const (
	filePath         = "example_video.mp4"
	objectName       = "example_video.mp4"
	downloadFilePath = "download_example_video.mp4"
)

type MinioKeys struct {
	MinioEndPoint   string
	MinioAccessKey  string
	MinioSecretKey  string
	MinioSecure     bool
	MinioBucketName string
	MinioLocation   string
}

// MinioInitTest 初始化minio client
func MinioInitTest(c MinioKeys) *oss.MinioClient {
	// 加载minio的endpoint等

	// log.Println(c.MinioEndPoint, c.MinioAccessKey, c.MinioSecretKey)

	client, err := oss.NewMinioClient(c.MinioEndPoint, c.MinioAccessKey, c.MinioSecretKey, c.MinioSecure)
	if err != nil {
		log.Println("init minio fail")
		return nil
	}
	return client
}

var minioKeys = MinioKeys{
	"192.168.3.125:9000",
	// "pupcFveqqthin0Q3eGsw",
	// "vJV6ryU3Qzkla6PeZta7MrfzpLxfduItjYgyzCns",
	"LH1I4a83WdghqEmTnZmw",
	"3xRjjBTJ3ZAcOIvgcDAmkOZQvY8NYAIefNSFTpkK",
	true,
	"tiktok",
	"cn-east-1",
}

var ctx = context.Background()

// 测试上传文件
func TestMinioUpload(t *testing.T) {
	client := MinioInitTest(minioKeys)

	putObjectOptions := minio.PutObjectOptions{}
	err := client.UploadFile(ctx, minioKeys.MinioBucketName, objectName, filePath, putObjectOptions)
	if err != nil {
		t.Error("upload file: ", err)
	}
}

// 测试下载文件
func TestMinioDownloadFile(t *testing.T) {
	client := MinioInitTest(minioKeys)
	getObjectOptions := minio.GetObjectOptions{}
	err := client.DownloadFile(ctx, minioKeys.MinioBucketName, objectName, downloadFilePath, getObjectOptions)
	if err != nil {
		t.Error("download file: ", err)
	}
}

// 测试列出bucket下所有的对象
func TestListObjects(t *testing.T) {
	client := MinioInitTest(minioKeys)
	listObjectsOptions := minio.ListObjectsOptions{}
	objects, err := client.ListObjects(ctx, minioKeys.MinioBucketName, listObjectsOptions)
	if err != nil {
		t.Error("list object: ", err)
	}
	t.Logf("objects: %+v", objects)
}

// 删除对象
func TestDeleteFile(t *testing.T) {
	client := MinioInitTest(minioKeys)
	removeObjectOptions := minio.RemoveObjectOptions{}
	ret, err := client.DeleteFile(ctx, minioKeys.MinioBucketName, objectName, removeObjectOptions)
	if err != nil {
		t.Error("delete object: ", ret, err)
	}
	log.Println("delete object: ", ret)
}

// 获取对象，返回url,
func TestGetPresignedGetObject(t *testing.T) {
	client := MinioInitTest(minioKeys)
	reqParams := url.Values{}
	// reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", objectName)) // 生成直接下载的文件
	object, err := client.GetPresignedGetObject(ctx, minioKeys.MinioBucketName, objectName, 24*time.Hour, reqParams)
	if err != nil {
		t.Error("GetPresignedGetObject: ", err)
	}
	log.Println("GetPresignedGetObject: ", object)
}
