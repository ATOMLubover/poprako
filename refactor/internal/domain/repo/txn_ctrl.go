package repo_iface

import "context"

type TxnFn func(context.Context) error

type TxnCtrl interface {
	// `RunWithTxn` executes the provided function within a transaction context.
	// If the function returns an error, the transaction will be rolled back.
	// If the function returns nil, the transaction will be committed.
	RunWithTxn(fn TxnFn) error
}
