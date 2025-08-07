package controllers

import (
	"gin-fataMorgana/services"
	"gin-fataMorgana/utils"

	"github.com/gin-gonic/gin"
)

// MessageController 消息控制器
type MessageController struct {
	messageService *services.MessageService
}

// NewMessageController 创建消息控制器实例
func NewMessageController() *MessageController {
	return &MessageController{
		messageService: services.NewMessageService(),
	}
}

// GetUserMessage 获取用户消息推送
// @Summary 获取用户消息推送
// @Description 从Redis消息队列中获取一条待推送的消息，并更新消息状态为已读
// @Tags 消息推送
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=models.UserMessageResponse} "成功返回消息"
// @Success 200 {object} utils.Response{data=null} "没有待推送的消息"
// @Failure 401 {object} utils.Response "认证失败"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /api/v2/message/pop [post]
func (mc *MessageController) GetUserMessage(c *gin.Context) {
	// 从上下文获取用户UID
	userUID, exists := c.Get("uid")
	if !exists {
		utils.Unauthorized(c)
		return
	}

	uid, ok := userUID.(string)
	if !ok {
		utils.InvalidParamsWithMessage(c, "用户UID格式错误")
		return
	}

	// 获取用户消息
	message, err := mc.messageService.GetUserMessage(c.Request.Context(), uid)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeServer, "获取消息失败")
		return
	}

	// 如果没有消息，返回空
	if message == nil {
		utils.Success(c, nil)
		return
	}

	// 返回消息
	utils.Success(c, message)
}