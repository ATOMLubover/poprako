package external

type OSSClient interface {
	GeneratePutPresignedURL(objectKey string, contentType string) (string, error)
	GenerateGetPresignedURL(objectKey string) (string, error)
	Delete(objectKey string) error
	DeleteBatch(objectKeys []string) error
}
