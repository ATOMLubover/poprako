package oss

// Client 定义对象存储能力的领域抽象接口
type Client interface {
	// GeneratePutPresignedURL 生成一个用于上传的预签名 URL，objectKey 是对象在 OSS 中的唯一标识
	GeneratePutPresignedURL(objectKey string) (string, error)
	// GenerateGetPresignedURL 生成一个用于下载的预签名 URL，objectKey 是对象在 OSS 中的唯一标识
	GenerateGetPresignedURL(objectKey string) (string, error)
	// Delete 删除 OSS 中的对象，objectKey 是对象在 OSS 中的唯一标识
	Delete(objectKey string) error
	// DeleteBatch 批量删除 OSS 中的对象，objectKeys 是对象在 OSS 中的唯一标识列表
	DeleteBatch(objectKeys []string) error
}
