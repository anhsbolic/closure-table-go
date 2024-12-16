package repository

import (
	"context"
	"database/sql"
	"github.com/anhsbolic/closure-table-go/model/domain"
	"github.com/anhsbolic/closure-table-go/pkg"
	"github.com/lib/pq"
)

type NodeClosureRepositoryImpl struct {
}

func NewNodeClosureRepository() NodeClosureRepository {
	return &NodeClosureRepositoryImpl{}
}

func (repository *NodeClosureRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, nodeClosure domain.NodeClosure) (domain.NodeClosure, error) {
	query := `INSERT INTO node_closure (ancestor, descendant, depth) VALUES ($1, $2, $3)`
	_, err := tx.ExecContext(ctx, query,
		nodeClosure.Ancestor,
		nodeClosure.Descendant,
		nodeClosure.Depth)

	if err != nil {
		return domain.NodeClosure{}, err
	}
	return nodeClosure, nil
}

func (repository *NodeClosureRepositoryImpl) DeleteByDescendantIds(ctx context.Context, tx *sql.Tx, descendantIds []string) error {
	query := `DELETE FROM node_closure WHERE descendant = ANY($1)`
	_, err := tx.ExecContext(ctx, query, pq.Array(descendantIds))

	if err != nil {
		return err
	}
	return nil
}

func (repository *NodeClosureRepositoryImpl) FindDescendantIdsByAncestor(ctx context.Context, tx *sql.Tx, ancestorId string) ([]string, error) {
	query := `SELECT descendant FROM node_closure WHERE ancestor = $1`
	rows, err := tx.QueryContext(ctx, query, ancestorId)
	if err != nil {
		return nil, err
	}
	defer pkg.CloseRows(rows)

	var descendantIds []string
	for rows.Next() {
		var descendantID string
		err := rows.Scan(&descendantID)
		if err != nil {
			return nil, err
		}
		descendantIds = append(descendantIds, descendantID)
	}

	return descendantIds, nil
}

func (repository *NodeClosureRepositoryImpl) FindByDescendant(ctx context.Context, db *sql.DB, nodeID string) ([]domain.NodeClosure, error) {
	query := `SELECT ancestor, descendant, depth FROM node_closure WHERE descendant = $1 ORDER BY depth`
	rows, err := db.QueryContext(ctx, query, nodeID)
	if err != nil {
		return nil, err
	}
	defer pkg.CloseRows(rows)

	var nodeClosures []domain.NodeClosure
	for rows.Next() {
		nodeClosure := domain.NodeClosure{}
		err := rows.Scan(&nodeClosure.Ancestor, &nodeClosure.Descendant, &nodeClosure.Depth)
		if err != nil {
			return nil, err
		}
		nodeClosures = append(nodeClosures, nodeClosure)
	}

	return nodeClosures, nil
}

func (repository *NodeClosureRepositoryImpl) GetNewClosures(ctx context.Context, tx *sql.Tx, nodeId string, newAncestorId string) ([]domain.NodeClosure, error) {
	query := `SELECT
				super_tree.ancestor,
				sub_tree.descendant,
				super_tree.depth + sub_tree.depth + 1 as depth
			FROM
				node_closure AS super_tree
			JOIN
				node_closure AS sub_tree ON sub_tree.ancestor = $1
			WHERE
				super_tree.descendant = $2`
	rows, err := tx.QueryContext(ctx, query, nodeId, newAncestorId)
	if err != nil {
		return nil, err
	}
	defer pkg.CloseRows(rows)

	var nodeClosures []domain.NodeClosure
	for rows.Next() {
		nodeClosure := domain.NodeClosure{}
		err := rows.Scan(&nodeClosure.Ancestor, &nodeClosure.Descendant, &nodeClosure.Depth)
		if err != nil {
			return nil, err
		}
		nodeClosures = append(nodeClosures, nodeClosure)
	}

	return nodeClosures, nil
}
