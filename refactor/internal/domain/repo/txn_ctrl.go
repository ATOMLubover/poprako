package repo_iface

type TxnFn[T any] func(Prov) (T, error)

type RawTxnFn func(Prov) error

type TxnCtrl interface {
	// `RunWithTxn` executes the provided function within a transaction context.
	// If the function returns an error, the transaction will be rolled back.
	// If the function returns nil, the transaction will be committed.
	RunWithTxn(fn RawTxnFn) error
}

func RunWithTxn[T any](ctrl TxnCtrl, fn TxnFn[T]) (T, error) {
	var re T

	err := ctrl.RunWithTxn(func(prov Prov) error {
		var err error

		re, err = fn(prov)
		if err != nil {
			// Make the transaction roll back by returning the error.
			return err
		}

		return nil
	})

	return re, err
}
