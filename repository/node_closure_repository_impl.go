package repository

import (
	"database/sql"
	"github.com/anhsbolic/closure-table-go/model/domain"
	"github.com/anhsbolic/closure-table-go/pkg"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

type NodeClosureRepositoryImpl struct {
}

func NewNodeClosureRepository() NodeClosureRepository {
	return &NodeClosureRepositoryImpl{}
}

func (repository *NodeClosureRepositoryImpl) Save(ctx *fiber.Ctx, tx *sql.Tx, nodeClosure domain.NodeClosure) domain.NodeClosure {
	SQL := `INSERT INTO node_closure (ancestor, descendant, depth) VALUES ($1, $2, $3)`
	_, err := tx.ExecContext(ctx.Context(), SQL,
		nodeClosure.Ancestor,
		nodeClosure.Descendant,
		nodeClosure.Depth)

	// Panic if error
	pkg.PanicIfError(err)

	return nodeClosure
}

func (repository *NodeClosureRepositoryImpl) DeleteByDescendantIds(ctx *fiber.Ctx, tx *sql.Tx, descendantIds []string) error {
	// Delete Node Closure By Descendant Ids
	SQL := `DELETE FROM node_closure WHERE descendant = ANY($1)`
	_, err := tx.ExecContext(ctx.Context(), SQL, pq.Array(descendantIds))

	// Panic if error
	pkg.PanicIfError(err)

	return nil
}

func (repository *NodeClosureRepositoryImpl) FindDescendantIdsByAncestor(ctx *fiber.Ctx, tx *sql.Tx, ancestorId string) []string {
	SQL := `SELECT descendant FROM node_closure WHERE ancestor = $1`
	rows, err := tx.QueryContext(ctx.Context(), SQL, ancestorId)

	// Panic if error
	pkg.PanicIfError(err)

	// Close Rows
	defer pkg.CloseRows(rows)

	// Create descendantIds
	var descendantIds []string
	for rows.Next() {
		var descendantID string
		err := rows.Scan(&descendantID)
		pkg.PanicIfError(err)

		descendantIds = append(descendantIds, descendantID)
	}

	return descendantIds
}

func (repository *NodeClosureRepositoryImpl) FindByDescendant(ctx *fiber.Ctx, db *sql.DB, nodeID string) []domain.NodeClosure {
	SQL := `SELECT ancestor, descendant, depth FROM node_closure WHERE descendant = $1 ORDER BY depth`
	rows, err := db.QueryContext(ctx.Context(), SQL, nodeID)

	// Panic if error
	pkg.PanicIfError(err)

	// Close Rows
	defer pkg.CloseRows(rows)

	// Create nodeClosures
	var nodeClosures []domain.NodeClosure
	for rows.Next() {
		nodeClosure := domain.NodeClosure{}
		err := rows.Scan(&nodeClosure.Ancestor, &nodeClosure.Descendant, &nodeClosure.Depth)
		pkg.PanicIfError(err)

		nodeClosures = append(nodeClosures, nodeClosure)
	}

	return nodeClosures
}

func (repository *NodeClosureRepositoryImpl) GetNewClosures(ctx *fiber.Ctx, tx *sql.Tx, nodeId string, newAncestorId string) []domain.NodeClosure {
	SQL := `SELECT
				super_tree.ancestor,
				sub_tree.descendant,
				super_tree.depth + sub_tree.depth + 1 as depth
			FROM
				node_closure AS super_tree
			JOIN
				node_closure AS sub_tree ON sub_tree.ancestor = $1
			WHERE
				super_tree.descendant = $2`
	rows, err := tx.QueryContext(ctx.Context(), SQL, nodeId, newAncestorId)

	// Panic if error
	pkg.PanicIfError(err)

	// Close Rows
	defer pkg.CloseRows(rows)

	// Create nodeClosures
	var nodeClosures []domain.NodeClosure
	for rows.Next() {
		nodeClosure := domain.NodeClosure{}
		err := rows.Scan(&nodeClosure.Ancestor, &nodeClosure.Descendant, &nodeClosure.Depth)
		pkg.PanicIfError(err)

		nodeClosures = append(nodeClosures, nodeClosure)
	}

	return nodeClosures
}
