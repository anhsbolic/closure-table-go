package service

import (
	"context"
	"github.com/anhsbolic/closure-table-go/model/dto"
)

type NodeService interface {
	Create(ctx context.Context, request dto.NodeCreateRequest) (dto.NodeCreatedResponse, error)
	RootList(ctx context.Context) ([]dto.NodeResponse, error)
	DetailNode(ctx context.Context, nodeId string) (dto.NodeResponse, error)
	UpdateNode(ctx context.Context, nodeId string, request dto.NodeUpdateRequest) (dto.NodeResponse, error)
	DeleteNode(ctx context.Context, nodeId string) error
	DescendantList(ctx context.Context, nodeId string) ([]dto.NodeResponse, error)
	MoveNode(ctx context.Context, nodeId string, request dto.NodeMoveRequest) error
}
