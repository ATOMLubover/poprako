package oss

// URLSigner 定义 OSS 预签名 URL 生成能力（用于上传/下载）
// 业务 app、读路径均只依赖此接口，不接触删除能力
type URLSigner interface {
	// GeneratePutPresignedURL 生成一个用于上传的预签名 URL，objectKey 是对象在 OSS 中的唯一标识
	GeneratePutPresignedURL(objectKey string) (string, error)
	// GenerateGetPresignedURL 生成一个用于下载的预签名 URL，objectKey 是对象在 OSS 中的唯一标识
	GenerateGetPresignedURL(objectKey string) (string, error)
}

// Deleter 定义 OSS 对象删除能力
// 仅供 OSSManagerSvc 的后台 worker 使用，业务 app 不得直接依赖
// 契约：删除不存在的对象必须视为成功；批量删除中不存在的对象不应导致整体失败
type Deleter interface {
	// Delete 删除 OSS 中的对象，objectKey 是对象在 OSS 中的唯一标识
	Delete(objectKey string) error
	// DeleteBatch 批量删除 OSS 中的对象，objectKeys 是对象在 OSS 中的唯一标识列表
	DeleteBatch(objectKeys []string) error
}

// Client 同时具备 URL 生成与对象删除能力（用于构造阶段 DI，不对外暴露）
// 新代码应优先依赖 URLSigner 或 Deleter，而非直接依赖 Client
type Client interface {
	URLSigner
	Deleter
}
