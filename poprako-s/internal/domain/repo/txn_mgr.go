package repo

import "context"

type TxnMgr interface {
	// RunInTxn 在一个事务上下文中执行给定的函数
	RunInTxn(func(ctx context.Context) error) error
}
