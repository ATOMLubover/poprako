package oss_iface

type Signer interface {
	GenGetURL(key string) (string, error)
	GenPutURL(key string) (string, error)
}
