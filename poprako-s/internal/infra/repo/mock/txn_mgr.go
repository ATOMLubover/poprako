package mock_repo

import "context"

type TxnMgr struct {
	Context context.Context
	RunErr  error
}

func NewMockTxnMgr(ctx context.Context) *TxnMgr {
	return &TxnMgr{Context: ctx}
}

func (m *TxnMgr) RunInTxn(fn func(cx context.Context) error) error {
	if m.RunErr != nil {
		return m.RunErr
	}

	ctx := m.Context
	if ctx == nil {
		ctx = context.Background()
	}

	return fn(ctx)
}
