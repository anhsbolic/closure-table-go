package repository

import (
	"context"
	"database/sql"
	"github.com/anhsbolic/closure-table-go/model/domain"
	"github.com/anhsbolic/closure-table-go/pkg"
	"github.com/lib/pq"
)

type NodeRepositoryImpl struct {
}

func NewNodeRepository() NodeRepository {
	return &NodeRepositoryImpl{}
}

func (repository *NodeRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, node domain.Node) (domain.Node, error) {
	// Save Root Node
	query := `INSERT INTO nodes (id, title, type, description, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := tx.QueryRowContext(ctx, query,
		node.ID,
		node.Title,
		node.Type,
		node.Description,
		node.CreatedAt,
	).Scan(&node.ID)

	if err != nil {
		return domain.Node{}, err
	}

	return node, nil
}

func (repository *NodeRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, id string, node domain.Node) (domain.Node, error) {
	query := `UPDATE nodes SET title = $1, type = $2, description = $3, updated_at = $4 WHERE id = $5`
	_, err := tx.ExecContext(ctx, query,
		node.Title,
		node.Type,
		node.Description,
		node.UpdatedAt,
		id,
	)
	if err != nil {
		return domain.Node{}, err
	}

	return node, nil
}

func (repository *NodeRepositoryImpl) DeleteByDescendantIds(ctx context.Context, tx *sql.Tx, descendantIds []string) error {
	query := `DELETE FROM nodes WHERE id = ANY($1)`
	_, err := tx.ExecContext(ctx, query, pq.Array(descendantIds))
	if err != nil {
		return err
	}

	return nil
}

func (repository *NodeRepositoryImpl) GetRootList(ctx context.Context, db *sql.DB) ([]domain.Node, error) {
	// Get Root List
	query := `SELECT n.id, n.title, n.type, n.description, n.created_at, n.updated_at
			FROM nodes n
			    JOIN node_closure nc ON n.id = nc.descendant
			WHERE nc.ancestor = nc.descendant
			  AND nc.depth = 0
			  AND NOT EXISTS (SELECT 1
			                  FROM node_closure nc2
			                  WHERE nc2.descendant = nc.descendant
			                    AND nc2.ancestor != nc.descendant)
			ORDER BY n.created_at DESC`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer pkg.CloseRows(rows)

	var nodes []domain.Node
	for rows.Next() {
		node := domain.Node{}
		err := rows.Scan(
			&node.ID,
			&node.Title,
			&node.Type,
			&node.Description,
			&node.CreatedAt,
			&node.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func (repository *NodeRepositoryImpl) CheckByID(ctx context.Context, db *sql.DB, id string) (bool, error) {
	query := `SELECT id FROM nodes WHERE id = $1`
	rows, err := db.QueryContext(ctx, query, id)
	if err != nil {
		return false, err
	}

	return rows.Next(), nil
}

func (repository *NodeRepositoryImpl) DetailByID(ctx context.Context, db *sql.DB, id string) (domain.Node, error) {
	query := `SELECT id, title, type, description, created_at, updated_at FROM nodes WHERE id = $1`
	row := db.QueryRowContext(ctx, query, id)

	node := domain.Node{}
	err := row.Scan(
		&node.ID,
		&node.Title,
		&node.Type,
		&node.Description,
		&node.CreatedAt,
		&node.UpdatedAt,
	)
	if err != nil {
		return domain.Node{}, err
	}

	return node, nil
}

func (repository *NodeRepositoryImpl) GetDescendantList(ctx context.Context, db *sql.DB, nodeId string) ([]domain.Node, error) {
	// Get Descendant List
	query := `SELECT n.id, n.title, n.type, n.description, n.created_at, n.updated_at
			FROM nodes n
			    JOIN node_closure nc ON n.id = nc.descendant
			WHERE nc.ancestor = $1
			  AND nc.depth > 0
			ORDER BY n.created_at DESC`
	rows, err := db.QueryContext(ctx, query, nodeId)
	if err != nil {
		return nil, err
	}
	defer pkg.CloseRows(rows)

	var nodes []domain.Node
	for rows.Next() {
		node := domain.Node{}
		err := rows.Scan(
			&node.ID,
			&node.Title,
			&node.Type,
			&node.Description,
			&node.CreatedAt,
			&node.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}
