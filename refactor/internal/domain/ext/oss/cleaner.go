package oss_iface

type Cleaner interface {
	DelBatch(keys []string) error
}
