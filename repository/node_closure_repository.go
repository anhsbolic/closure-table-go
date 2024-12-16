package repository

import (
	"context"
	"database/sql"
	"github.com/anhsbolic/closure-table-go/model/domain"
)

type NodeClosureRepository interface {
	Save(ctx context.Context, tx *sql.Tx, nodeClosures domain.NodeClosure) (domain.NodeClosure, error)
	DeleteByDescendantIds(ctx context.Context, tx *sql.Tx, descendantIds []string) error
	FindDescendantIdsByAncestor(ctx context.Context, tx *sql.Tx, ancestorId string) ([]string, error)
	FindByDescendant(ctx context.Context, db *sql.DB, nodeID string) ([]domain.NodeClosure, error)
	GetNewClosures(ctx context.Context, tx *sql.Tx, nodeId string, newAncestorId string) ([]domain.NodeClosure, error)
}
