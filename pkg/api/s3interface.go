package api

import "github.com/aws/aws-sdk-go/service/s3"

type S3API interface {
	// 创建一个新的存储桶
	CreateBucket(input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error)

	// 列出存储桶中的对象
	ListObjects(input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error)

	// 上传一个对象到存储桶中
	PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)

	// 下载存储桶中的一个对象
	GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error)

	// 删除存储桶中的一个对象
	DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)

	// 删除存储桶中的多个对象
	DeleteObjects(input *s3.DeleteObjectsInput) (*s3.DeleteObjectsOutput, error)

	// 复制存储桶中的一个对象
	CopyObject(input *s3.CopyObjectInput) (*s3.CopyObjectOutput, error)

	// 获取存储桶中的对象元数据
	HeadObject(input *s3.HeadObjectInput) (*s3.HeadObjectOutput, error)

	// 列出存储桶中的所有存储桶版本
	ListObjectVersions(input *s3.ListObjectVersionsInput) (*s3.ListObjectVersionsOutput, error)

	// 上传一个分块对象到存储桶中
	CreateMultipartUpload(input *s3.CreateMultipartUploadInput) (*s3.CreateMultipartUploadOutput, error)

	// 上传一个分块对象的一部分到存储桶中
	UploadPart(input *s3.UploadPartInput) (*s3.UploadPartOutput, error)

	// 完成分块对象的上传
	CompleteMultipartUpload(input *s3.CompleteMultipartUploadInput) (*s3.CompleteMultipartUploadOutput, error)

	// 取消分块对象的上传
	AbortMultipartUpload(input *s3.AbortMultipartUploadInput) (*s3.AbortMultipartUploadOutput, error)

	// 列出存储桶中的所有已上传但未完成的分块对象
	ListParts(input *s3.ListPartsInput) (*s3.ListPartsOutput, error)
}
