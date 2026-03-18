package external

type OSSClient interface {
	GeneratePutPresignedURL(objectKey string) (string, error)
	GenerateGetPresignedURL(objectKey string) (string, error)
	Delete(objectKey string) error
	DeleteBatch(objectKeys []string) error
}
