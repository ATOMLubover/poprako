package repo_infra

import (
	repo_iface "poprako-s/internal/domain/repo"

	"gorm.io/gorm"
)

type txnCtrlImpl struct {
	gdb *gorm.DB
}

func NewTxnCtrl(gdb *gorm.DB) repo_iface.TxnCtrl {
	return &txnCtrlImpl{gdb: gdb}
}

func (t *txnCtrlImpl) RunWithTxn(fn repo_iface.RawTxnFn) error {
	return t.gdb.Transaction(func(tx *gorm.DB) error {
		// Create a provider with the transaction DB so transaction functions can
		// get repositories safely.
		prov := newProv(tx)

		return fn(prov)
	})
}
