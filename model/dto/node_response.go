package dto

import (
	"github.com/anhsbolic/closure-table-go/model/domain"
	"github.com/anhsbolic/closure-table-go/pkg"
	"github.com/google/uuid"
	"time"
)

type NodeCreatedResponse struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Type        string     `json:"type"`
	Description *string    `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
}

func ToNodeCreatedResponse(node domain.Node) NodeCreatedResponse {
	return NodeCreatedResponse{
		ID:          node.ID,
		Title:       node.Title,
		Type:        node.Type,
		Description: pkg.NullStringToPointer(node.Description),
		CreatedAt:   pkg.NullTimeToPointer(node.CreatedAt),
	}
}

type NodeResponse struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Type        string     `json:"type"`
	Description *string    `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

func ToNodePaginationResponse(nodes []domain.Node) []NodeResponse {
	var nodeResponses []NodeResponse

	for _, node := range nodes {
		nodeResponses = append(nodeResponses, NodeResponse{
			ID:          node.ID,
			Title:       node.Title,
			Type:        node.Type,
			Description: pkg.NullStringToPointer(node.Description),
			CreatedAt:   pkg.NullTimeToPointer(node.CreatedAt),
			UpdatedAt:   pkg.NullTimeToPointer(node.UpdatedAt),
		})
	}

	return nodeResponses
}

func ToNodeDetailResponse(node domain.Node) NodeResponse {
	return NodeResponse{
		ID:          node.ID,
		Title:       node.Title,
		Type:        node.Type,
		Description: pkg.NullStringToPointer(node.Description),
		CreatedAt:   pkg.NullTimeToPointer(node.CreatedAt),
		UpdatedAt:   pkg.NullTimeToPointer(node.UpdatedAt),
	}
}
