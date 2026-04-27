package oss_iface

type Signer interface {
	GenGetUrl(key string) (string, error)
	GenPutUrl(key string) (string, error)
}
