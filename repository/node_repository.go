package repository

import (
	"database/sql"
	"github.com/anhsbolic/closure-table-go/model/domain"
	"github.com/gofiber/fiber/v2"
)

type NodeRepository interface {
	Create(ctx *fiber.Ctx, tx *sql.Tx, node domain.Node) domain.Node
	Update(ctx *fiber.Ctx, tx *sql.Tx, id string, node domain.Node) domain.Node
	DeleteByDescendantIds(ctx *fiber.Ctx, tx *sql.Tx, descendantIds []string) error
	GetRootList(ctx *fiber.Ctx, db *sql.DB) []domain.Node
	CheckByID(ctx *fiber.Ctx, db *sql.DB, id string) bool
	DetailByID(ctx *fiber.Ctx, db *sql.DB, id string) domain.Node
	GetDescendantList(ctx *fiber.Ctx, db *sql.DB, nodeId string) []domain.Node
}
