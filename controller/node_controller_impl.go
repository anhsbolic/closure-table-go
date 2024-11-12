package controller

import (
	"closure-table-go/model/dto"
	"closure-table-go/service"
	"github.com/gofiber/fiber/v2"
)

type NodeControllerImpl struct {
	NodeService service.NodeService
}

func NewNodeController(categoryService service.NodeService) NodeController {
	return &NodeControllerImpl{
		NodeService: categoryService,
	}
}

func (controller *NodeControllerImpl) Create(ctx *fiber.Ctx) error {
	request := new(dto.NodeCreateRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		return err
	}

	result, err := controller.NodeService.Create(ctx, *request)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(dto.ApiResponseSuccess{
		Success: true,
		Message: "Node has been created",
		Data:    result,
	})
}

func (controller *NodeControllerImpl) RootList(ctx *fiber.Ctx) error {
	result, err := controller.NodeService.RootList(ctx)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ApiResponseSuccess{
		Success: true,
		Message: "List of root nodes",
		Data:    result,
	})
}

func (controller *NodeControllerImpl) DetailNode(ctx *fiber.Ctx) error {
	nodeId := ctx.Params("nodeId")
	result, err := controller.NodeService.DetailNode(ctx, nodeId)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ApiResponseSuccess{
		Success: true,
		Message: "Detail of node",
		Data:    result,
	})
}

func (controller *NodeControllerImpl) UpdateNode(ctx *fiber.Ctx) error {
	nodeId := ctx.Params("nodeId")
	request := new(dto.NodeUpdateRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		return err
	}

	result, err := controller.NodeService.UpdateNode(ctx, nodeId, *request)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ApiResponseSuccess{
		Success: true,
		Message: "Node detail has been updated",
		Data:    result,
	})
}

func (controller *NodeControllerImpl) DeleteNode(ctx *fiber.Ctx) error {
	nodeId := ctx.Params("nodeId")
	err := controller.NodeService.DeleteNode(ctx, nodeId)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ApiResponseSuccess{
		Success: true,
		Message: "Node with all descendants has been deleted",
	})
}

func (controller *NodeControllerImpl) DescendantList(ctx *fiber.Ctx) error {
	nodeId := ctx.Params("nodeId")
	result, err := controller.NodeService.DescendantList(ctx, nodeId)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ApiResponseSuccess{
		Success: true,
		Message: "List of descendant nodes",
		Data:    result,
	})
}

func (controller *NodeControllerImpl) MoveNode(ctx *fiber.Ctx) error {
	nodeId := ctx.Params("nodeId")
	request := new(dto.NodeMoveRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		return err
	}

	err = controller.NodeService.MoveNode(ctx, nodeId, *request)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ApiResponseSuccess{
		Success: true,
		Message: "Node has been moved",
	})
}