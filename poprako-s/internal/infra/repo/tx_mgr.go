package repo

import (
	"context"

	iface "poprako-s/internal/domain/repo"

	"gorm.io/gorm"
)

type txnMgrImpl struct {
	gdb *gorm.DB
}

func NewTxnMgr(gdb *gorm.DB) iface.TxnMgr {
	return &txnMgrImpl{
		gdb: gdb,
	}
}

const txnKey = "txn"

func (m *txnMgrImpl) RunInTxn(fn func(ctx context.Context) error) error {
	return m.gdb.Transaction(func(tx *gorm.DB) error {
		ctx := context.WithValue(context.Background(), txnKey, tx)
		return fn(ctx)
	})
}
