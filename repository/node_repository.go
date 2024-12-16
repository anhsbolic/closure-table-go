package repository

import (
	"context"
	"database/sql"
	"github.com/anhsbolic/closure-table-go/model/domain"
)

type NodeRepository interface {
	Create(ctx context.Context, tx *sql.Tx, node domain.Node) (domain.Node, error)
	Update(ctx context.Context, tx *sql.Tx, id string, node domain.Node) (domain.Node, error)
	DeleteByDescendantIds(ctx context.Context, tx *sql.Tx, descendantIds []string) error
	GetRootList(ctx context.Context, db *sql.DB) ([]domain.Node, error)
	CheckByID(ctx context.Context, db *sql.DB, id string) (bool, error)
	DetailByID(ctx context.Context, db *sql.DB, id string) (domain.Node, error)
	GetDescendantList(ctx context.Context, db *sql.DB, nodeId string) ([]domain.Node, error)
}
