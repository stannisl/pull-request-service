package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/stannisl/pull-request-service/pkg/db"
)

type BaseRepository struct {
	db        *sqlx.DB
	txManager db.TransactionManager
}

func NewBaseRepository(db *sqlx.DB, txManager db.TransactionManager) *BaseRepository {
	return &BaseRepository{
		db:        db,
		txManager: txManager,
	}
}

func (r *BaseRepository) GetExecutor(ctx context.Context) sqlx.ExtContext {
	if tx, ok := ctx.Value(db.TxKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return r.db
}

func (r *BaseRepository) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.txManager.WithTransaction(ctx, fn)
}
