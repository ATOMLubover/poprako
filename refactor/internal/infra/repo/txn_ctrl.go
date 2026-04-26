package repo_infra

import (
	"context"

	repo_iface "poprako-s/internal/domain/repo"

	"gorm.io/gorm"
)

const TXN_GDB_KEY = "txn_gdb"

type txnCtrlImpl struct {
	gdb *gorm.DB
}

func NewTxnCtrl(gdb *gorm.DB) repo_iface.TxnCtrl {
	return &txnCtrlImpl{gdb: gdb}
}

func (t *txnCtrlImpl) RunWithTxn(fn repo_iface.TxnFn) error {
	return t.gdb.Transaction(func(tx *gorm.DB) error {
		cx := context.WithValue(context.Background(), TXN_GDB_KEY, tx)
		return fn(cx)
	})
}

func takeTxnGdb(cx context.Context) *gorm.DB {
	if gdb, ok := cx.Value(TXN_GDB_KEY).(*gorm.DB); ok {
		return gdb
	}

	return nil
}
