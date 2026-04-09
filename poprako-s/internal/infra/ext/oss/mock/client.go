package mock_oss

import "poprako-s/internal/domain/ext/oss"

type Client struct {
	generatePutPresignedURLFunc func(objectKey string) (string, error)
	generateGetPresignedURLFunc func(objectKey string) (string, error)
	deleteFunc                  func(objectKey string) error
	deleteBatchFunc             func(objectKeys []string) error

	putURLs        map[string]string
	getURLs        map[string]string
	deleted        []string
	deletedBatch   [][]string
	putErr         error
	getErr         error
	deleteErr      error
	deleteBatchErr error
}

func NewMockOSSClient() oss.Client {
	return &Client{}
}

var _ oss.Client = (*Client)(nil)

func (f *Client) SetGeneratePutPresignedURLFunc(fn func(objectKey string) (string, error)) {
	f.generatePutPresignedURLFunc = fn
}

func (f *Client) SetGenerateGetPresignedURLFunc(fn func(objectKey string) (string, error)) {
	f.generateGetPresignedURLFunc = fn
}

func (f *Client) SetDeleteFunc(fn func(objectKey string) error) {
	f.deleteFunc = fn
}

func (f *Client) SetDeleteBatchFunc(fn func(objectKeys []string) error) {
	f.deleteBatchFunc = fn
}

func (f *Client) SetPutURLs(urls map[string]string) {
	f.putURLs = urls
}

func (f *Client) SetGetURLs(urls map[string]string) {
	f.getURLs = urls
}

func (f *Client) SetPutErr(err error) {
	f.putErr = err
}

func (f *Client) SetGetErr(err error) {
	f.getErr = err
}

func (f *Client) SetDeleteErr(err error) {
	f.deleteErr = err
}

func (f *Client) SetDeleteBatchErr(err error) {
	f.deleteBatchErr = err
}

func (f *Client) Deleted() []string {
	return append([]string(nil), f.deleted...)
}

func (f *Client) DeletedBatch() [][]string {
	result := make([][]string, 0, len(f.deletedBatch))
	for _, batch := range f.deletedBatch {
		result = append(result, append([]string(nil), batch...))
	}
	return result
}

func (f *Client) GeneratePutPresignedURL(objectKey string) (string, error) {
	if f.generatePutPresignedURLFunc == nil {
		if f.putErr != nil {
			return "", f.putErr
		}
		if f.putURLs != nil {
			if url, ok := f.putURLs[objectKey]; ok {
				return url, nil
			}
		}
		return objectKey, nil
	}
	return f.generatePutPresignedURLFunc(objectKey)
}

func (f *Client) GenerateGetPresignedURL(objectKey string) (string, error) {
	if f.generateGetPresignedURLFunc == nil {
		if f.getErr != nil {
			return "", f.getErr
		}
		if f.getURLs != nil {
			if url, ok := f.getURLs[objectKey]; ok {
				return url, nil
			}
		}
		return objectKey, nil
	}
	return f.generateGetPresignedURLFunc(objectKey)
}

func (f *Client) Delete(objectKey string) error {
	if f.deleteFunc == nil {
		f.deleted = append(f.deleted, objectKey)
		return f.deleteErr
	}
	return f.deleteFunc(objectKey)
}

func (f *Client) DeleteBatch(objectKeys []string) error {
	if f.deleteBatchFunc == nil {
		copied := append([]string(nil), objectKeys...)
		f.deletedBatch = append(f.deletedBatch, copied)
		return f.deleteBatchErr
	}
	return f.deleteBatchFunc(objectKeys)
}
