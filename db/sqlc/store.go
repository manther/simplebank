package db

import (
	"context"
	"database/sql"
	"fmt"
)

type SQLStore struct {
	*Queries
	db *sql.DB
}

type Store interface {
	TransferTX(ctx context.Context, arg TransferTxParms) (TransferTxResult, error)
	Querier
}

func NewStore(db *sql.DB) Store {
	store := SQLStore{
		db:      db,
		Queries: New(db),
	}
	return &store
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if errRb := tx.Rollback(); errRb != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, errRb)
		}
	}
	return tx.Commit()
}

// TransferTxParms contains the input parameters of the transfer transaction
type TransferTxParms struct {
	FromAccountID int64 `json:"from_account"`
	ToAccountID   int64 `json:"to_account"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult holds the results of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *SQLStore) TransferTX(ctx context.Context, arg TransferTxParms) (TransferTxResult, error) {
	var transfResult TransferTxResult

	fn := func(q *Queries) error {
		var err error
		// fmt.Println("CreateTransfer:", txVal)
		transfResult.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccount: arg.FromAccountID,
			ToAccount:   arg.ToAccountID,
			Amount:      arg.Amount,
		})
		if err != nil {
			return err
		}

		// fmt.Println("CreateEntryFrom:", txVal)
		transfResult.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// fmt.Println("CreateEntryTo:", txVal)
		transfResult.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		var fnFrm = func() error {
			transfResult.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				Accountid: arg.FromAccountID,
				Amount:    -arg.Amount,
			})
			if err != nil {
				return err
			}
			return nil
		}

		var fnTo = func() error {
			transfResult.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				Accountid: arg.ToAccountID,
				Amount:    arg.Amount,
			})
			if err != nil {
				return err
			}
			return nil
		}

		return addAccountBalanceOrdered(arg.FromAccountID, arg.ToAccountID, fnFrm, fnTo)
	}

	err := store.execTx(ctx, fn)
	return transfResult, err
}

func addAccountBalanceOrdered(id1 int64, id2 int64, funcFrom func() error, funcTo func() error) error {
	if id1 < id2 {
		err := funcFrom()
		if err != nil {
			return err
		}

		err = funcTo()
		if err != nil {
			return err
		}
	} else {
		err := funcTo()
		if err != nil {
			return err
		}

		err = funcFrom()
		if err != nil {
			return err
		}
	}

	return nil
}
