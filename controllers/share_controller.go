package controllers

import (
	"context"
	"gin-fataMorgana/services"
	"gin-fataMorgana/utils"

	"github.com/gin-gonic/gin"
)

// ShareController 分享相关控制器
type ShareController struct {
	shareService *services.ShareService
}

func NewShareController() *ShareController {
	return &ShareController{
		shareService: services.NewShareService(),
	}
}

// GetShareLink 获取分享链接
func (c *ShareController) GetShareLink(ctx *gin.Context) {
	link, err := c.shareService.GetShareLink(context.Background())
	if err != nil {
		utils.ErrorWithMessage(ctx, utils.CodeDatabaseError, "获取分享链接失败")
		return
	}
	utils.Success(ctx, gin.H{"share_link": link})
}
