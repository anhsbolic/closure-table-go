package service

import (
	"context"
	"database/sql"
	"github.com/anhsbolic/closure-table-go/model/domain"
	"github.com/anhsbolic/closure-table-go/model/dto"
	"github.com/anhsbolic/closure-table-go/pkg"
	"github.com/anhsbolic/closure-table-go/repository"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
)

type NodeServiceImpl struct {
	NodeRepository        repository.NodeRepository
	NodeClosureRepository repository.NodeClosureRepository
	DB                    *sql.DB
	Validate              *validator.Validate
}

func NewNodeService(
	nodeRepository repository.NodeRepository,
	nodeClosureRepository repository.NodeClosureRepository,
	db *sql.DB,
	validate *validator.Validate,
) NodeService {
	return &NodeServiceImpl{
		NodeRepository:        nodeRepository,
		NodeClosureRepository: nodeClosureRepository,
		DB:                    db,
		Validate:              validate,
	}
}

func (service *NodeServiceImpl) Create(ctx context.Context, request dto.NodeCreateRequest) (dto.NodeCreatedResponse, error) {
	// Validate request
	err := service.Validate.Struct(request)
	if err != nil {
		return dto.NodeCreatedResponse{}, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Check Ancestor Node
	if request.AncestorID != nil {
		isAncestorNodeExist, err := service.NodeRepository.CheckByID(ctx, service.DB, *request.AncestorID)
		if err != nil {
			return dto.NodeCreatedResponse{}, err
		}
		if !isAncestorNodeExist {
			return dto.NodeCreatedResponse{}, fiber.NewError(
				fiber.StatusUnprocessableEntity,
				"Ancestor node is not found",
			)
		}
	}

	// Start transaction
	tx, err := service.DB.Begin()
	if err != nil {
		return dto.NodeCreatedResponse{}, err
	}

	// Defer commit or rollback
	defer pkg.CommitOrRollback(tx)

	// Save node
	description := sql.NullString{Valid: false}
	if request.Description != nil {
		description = sql.NullString{String: *request.Description, Valid: true}
	}
	node := domain.Node{
		ID:          uuid.New(),
		Title:       request.Title,
		Type:        request.Type,
		Description: description,
		CreatedAt:   sql.NullTime{Time: time.Now(), Valid: true},
	}
	createdNode, err := service.NodeRepository.Create(ctx, tx, node)
	if err != nil {
		return dto.NodeCreatedResponse{}, err
	}

	// Save NodeClosure : Self Reference
	closure := domain.NodeClosure{
		Ancestor:   createdNode.ID,
		Descendant: createdNode.ID,
		Depth:      0,
	}
	_, err = service.NodeClosureRepository.Save(ctx, tx, closure)
	if err != nil {
		return dto.NodeCreatedResponse{}, err
	}

	// When Node Have Ancestor
	if request.AncestorID != nil {
		// Get Ancestor Closures
		ancestorClosures, err := service.NodeClosureRepository.FindByDescendant(ctx, service.DB, *request.AncestorID)
		if err != nil {
			return dto.NodeCreatedResponse{}, err
		}

		// Save NodeClosure : Ancestor Reference
		depth := 1
		for _, ancestorClosure := range ancestorClosures {
			closure := domain.NodeClosure{
				Ancestor:   ancestorClosure.Ancestor,
				Descendant: createdNode.ID,
				Depth:      depth,
			}
			_, err := service.NodeClosureRepository.Save(ctx, tx, closure)
			if err != nil {
				return dto.NodeCreatedResponse{}, err
			}
			depth++
		}
	}

	// return response
	return dto.ToNodeCreatedResponse(createdNode), nil
}

func (service *NodeServiceImpl) RootList(ctx context.Context) ([]dto.NodeResponse, error) {
	// Get Root Nodes
	rootNodes, err := service.NodeRepository.GetRootList(ctx, service.DB)
	if err != nil {
		return []dto.NodeResponse{}, err
	}

	// return response
	return dto.ToNodePaginationResponse(rootNodes), nil
}

func (service *NodeServiceImpl) DetailNode(ctx context.Context, nodeId string) (dto.NodeResponse, error) {
	// Get Node By ID
	node, err := service.NodeRepository.DetailByID(ctx, service.DB, nodeId)
	if err != nil {
		return dto.NodeResponse{}, err
	}
	if node.ID == uuid.Nil {
		return dto.NodeResponse{}, fiber.ErrNotFound
	}

	// return response
	return dto.ToNodeDetailResponse(node), nil
}

func (service *NodeServiceImpl) UpdateNode(ctx context.Context, nodeId string, request dto.NodeUpdateRequest) (dto.NodeResponse, error) {
	// Get Detail Node By ID
	node, err := service.NodeRepository.DetailByID(ctx, service.DB, nodeId)
	if err != nil {
		return dto.NodeResponse{}, err
	}
	if node.ID == uuid.Nil {
		return dto.NodeResponse{}, fiber.ErrNotFound
	}

	// Validate request
	err = service.Validate.Struct(request)
	if err != nil {
		return dto.NodeResponse{}, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Start transaction
	tx, err := service.DB.Begin()
	if err != nil {
		return dto.NodeResponse{}, err
	}

	// Defer commit or rollback
	defer pkg.CommitOrRollback(tx)

	// Update Node
	node.Title = request.Title
	node.Type = request.Type
	if request.Description != nil {
		node.Description = sql.NullString{String: *request.Description, Valid: true}
	}
	node.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	updatedNode, err := service.NodeRepository.Update(ctx, tx, nodeId, node)
	if err != nil {
		return dto.NodeResponse{}, err
	}

	// return response
	return dto.ToNodeDetailResponse(updatedNode), nil
}

func (service *NodeServiceImpl) DeleteNode(ctx context.Context, nodeId string) error {
	// Check Node By ID
	isNodeExist, err := service.NodeRepository.CheckByID(ctx, service.DB, nodeId)
	if err != nil {
		return err
	}
	if !isNodeExist {
		return fiber.ErrNotFound
	}

	// Start transaction
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}

	// Defer commit or rollback
	defer pkg.CommitOrRollback(tx)

	// Get Descendant IDs
	descendantIds, err := service.NodeClosureRepository.FindDescendantIdsByAncestor(ctx, tx, nodeId)
	if err != nil {
		return err
	}

	// Delete Node Closure : Self with All Descendants
	err = service.NodeClosureRepository.DeleteByDescendantIds(ctx, tx, descendantIds)
	if err != nil {
		return err
	}

	// Delete Node with All Descendants
	err = service.NodeRepository.DeleteByDescendantIds(ctx, tx, descendantIds)
	if err != nil {
		return err
	}

	// return response
	return nil
}

func (service *NodeServiceImpl) DescendantList(ctx context.Context, nodeId string) ([]dto.NodeResponse, error) {
	// Check Node By ID
	isNodeExist, err := service.NodeRepository.CheckByID(ctx, service.DB, nodeId)
	if err != nil {
		return []dto.NodeResponse{}, err
	}
	if !isNodeExist {
		return []dto.NodeResponse{}, fiber.ErrNotFound
	}

	// Get Descendant Nodes
	descendantNodes, err := service.NodeRepository.GetDescendantList(ctx, service.DB, nodeId)
	if err != nil {
		return []dto.NodeResponse{}, err
	}

	// return response
	return dto.ToNodePaginationResponse(descendantNodes), nil
}

func (service *NodeServiceImpl) MoveNode(ctx context.Context, nodeId string, request dto.NodeMoveRequest) error {
	// Check Node By ID
	isNodeExist, err := service.NodeRepository.CheckByID(ctx, service.DB, nodeId)
	if err != nil {
		return err
	}
	if !isNodeExist {
		return fiber.ErrNotFound
	}

	// Validate request
	err = service.Validate.Struct(request)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Check Ancestor Node
	isAncestorNodeExist, err := service.NodeRepository.CheckByID(ctx, service.DB, request.ToAncestorID)
	if err != nil {
		return err
	}
	if !isAncestorNodeExist {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Ancestor node is not found")
	}

	// Start transaction
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}

	// Defer commit or rollback
	defer pkg.CommitOrRollback(tx)

	// Get New Path For Node
	newClosures, err := service.NodeClosureRepository.GetNewClosures(ctx, tx, nodeId, request.ToAncestorID)
	if err != nil {
		return err
	}

	// Get Descendant IDs
	descendantIds, err := service.NodeClosureRepository.FindDescendantIdsByAncestor(ctx, tx, nodeId)
	if err != nil {
		return err
	}

	// Delete Node Closure : Self with All Descendants
	err = service.NodeClosureRepository.DeleteByDescendantIds(ctx, tx, descendantIds)
	if err != nil {
		return err
	}

	// Save New Node Closure For Self and All Descendants Under New Ancestor
	for _, closure := range newClosures {
		_, err := service.NodeClosureRepository.Save(ctx, tx, closure)
		if err != nil {
			return err
		}
	}

	// return success
	return nil
}
