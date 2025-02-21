package oss

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	defaultExpiryTime = time.Second * 24 * 60 * 60 // 1 day
)

var MinioClientGlobal *MinioClient

type MinioClient struct {
	Client *minio.Client
}

// NewMinioClient 初始化minio
func NewMinioClient(endpoint, accessKey, secretKey string, secure bool) (*MinioClient, error) {
	// 初始化 Minio 客户端

	// 跳过证书验证, 如果证书正常, 删除该代码
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 跳过证书验证
		},
	}
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:     credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure:    secure,
		Transport: transport,
	})
	if err != nil {
		log.Println("new minio client fail: ", err)
		return nil, err
	}
	client := &MinioClient{
		Client: minioClient,
	}
	MinioClientGlobal = client
	return client, nil
}

func (m *MinioClient) PostPresignedUrl(ctx context.Context, bucketName, objectName string) (string, map[string]string, error) {
	expiry := defaultExpiryTime

	policy := minio.NewPostPolicy()
	_ = policy.SetBucket(bucketName)
	_ = policy.SetKey(objectName)
	_ = policy.SetExpires(time.Now().UTC().Add(expiry))

	presignedURL, formData, err := m.Client.PresignedPostPolicy(ctx, policy)
	if err != nil {
		log.Fatalln(err)
		return "", map[string]string{}, err
	}

	return presignedURL.String(), formData, nil
}

func (m *MinioClient) PutPresignedUrl(ctx context.Context, bucketName, objectName string) (string, error) {
	expiry := defaultExpiryTime

	presignedURL, err := m.Client.PresignedPutObject(ctx, bucketName, objectName, expiry)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}

	return presignedURL.String(), nil
}

// UploadFile 上传文件
func (m *MinioClient) UploadFile(ctx context.Context, bucketName, objectName, filePath string, options minio.PutObjectOptions) error {
	// 打开本地文件
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("open filePath: %s fail: %s", filePath, err)
		return err
	}
	defer file.Close()

	// 上传文件到存储桶
	_, err = m.Client.FPutObject(ctx, bucketName, objectName, filePath, options)
	if err != nil {
		log.Println("putObject fail: ", err)
		return err
	}

	fmt.Println("Successfully uploaded", objectName)

	return nil
}

// DownloadFile 下载文件
func (m *MinioClient) DownloadFile(ctx context.Context, bucketName, objectName, filePath string, options minio.GetObjectOptions) error {
	// 创建本地文件
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 下载存储桶中的文件到本地
	err = m.Client.FGetObject(ctx, bucketName, objectName, filePath, options)
	if err != nil {
		return err
	}

	fmt.Println("Successfully downloaded", objectName)
	return nil
}

// DeleteFile 删除文件
func (m *MinioClient) DeleteFile(ctx context.Context, bucketName, objectName string, options minio.RemoveObjectOptions) (bool, error) {
	// 删除存储桶中的文件
	err := m.Client.RemoveObject(ctx, bucketName, objectName, options)
	if err != nil {
		log.Println("remove object fail: ", err)
		return false, err
	}

	fmt.Println("Successfully deleted", objectName)
	return true, err
}

// ListObjects 列出文件
func (m *MinioClient) ListObjects(ctx context.Context, bucketName string, options minio.ListObjectsOptions) ([]string, error) {
	var objectNames []string

	for object := range m.Client.ListObjects(ctx, bucketName, options) {
		if object.Err != nil {
			return nil, object.Err
		}

		objectNames = append(objectNames, object.Key)
	}

	return objectNames, nil
}

// GetPresignedGetObject 返回对象的url地址，有效期时间为expires
func (m *MinioClient) GetPresignedGetObject(ctx context.Context, bucketName string, objectName string, expiry time.Duration, reqParams url.Values) (string, error) {
	object, err := m.Client.PresignedGetObject(ctx, bucketName, objectName, expiry, reqParams)
	if err != nil {
		log.Println("get object fail: ", err)
		return "", err
	}

	return object.String(), nil
}
