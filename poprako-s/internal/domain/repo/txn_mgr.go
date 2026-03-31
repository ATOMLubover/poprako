package repo

import "context"

// TxnMgr 负责控制事务的范围和执行，提供一个在事务上下文中执行函数的接口
type TxnMgr interface {
	// RunInTxn 在一个事务上下文中执行给定的函数
	RunInTxn(func(ctx context.Context) error) error
}
