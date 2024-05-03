package api

import (
	"go.uber.org/zap"
	"mybbs-backend/vo"
	"strconv"

	"github.com/gin-gonic/gin"

	"mybbs-backend/common"
	"mybbs-backend/config"
	"mybbs-backend/dto"
	"mybbs-backend/logic"
)

func GetPostDetail(ctx *gin.Context) {
	pidStr := ctx.Param("post_id")
	if pidStr == "" {
		zap.L().Error("post_id is empty")
		common.FailByMsg(ctx, "帖子ID为空")
		return
	}
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("failed to parse post_id", zap.Error(err))
		common.FailByMsg(ctx, "解析帖子ID失败")
		return
	}
	detail, err := logic.ReadPost(pid)
	if err != nil {
		zap.L().Error("failed to get post detail", zap.Error(err))
		common.FailByMsg(ctx, "获取帖子详情失败")
		return
	}

	user := logic.GetUserById(detail.UserID)
	common.Success(ctx,
		vo.PostDetail{Username: user.Username,
			Post: detail})
}

func GetTopPostList(ctx *gin.Context) {
	query := dto.PostListQuery{Page: 1, PageSize: 10, Order: "click_count DESC"}
	if err := ctx.ShouldBind(&query); err != nil {
		config.ValidateError(ctx, err)
		return
	}
	posts, err := logic.GetTopPostList(query.Page, query.PageSize, query.Order)
	if err != nil {
		zap.L().Error("failed to get top post list", zap.Error(err))
		common.FailByMsg(ctx, "获取热门帖子列表失败")
		return
	}
	common.Success(ctx, posts)
}

func CreatePost(ctx *gin.Context) {
	var postDTO dto.PostDTO
	if err := ctx.ShouldBind(&postDTO); err != nil {
		config.ValidateError(ctx, err)
		return
	}
	userId, exists := ctx.Get("userId")
	if !exists {
		zap.L().Error("user not logged in")
		common.FailByMsg(ctx, "用户未登录")
		return
	}
	err := logic.CreatePost(userId.(int64), postDTO)
	if err != nil {
		zap.L().Error("failed to create post", zap.Error(err))
		common.FailByMsg(ctx, "创建帖子失败")
		return
	}
	common.Success(ctx, nil)
}

// UpdatePost updates a given post
func UpdatePost(ctx *gin.Context) {
	var postDTO dto.UpdatePostDTO
	var postIDParam dto.PostIDParam

	// Binding uri and json parameters from the request
	if err := ctx.ShouldBindUri(&postIDParam); err != nil {
		config.ValidateError(ctx, err)
		return
	}
	if err := ctx.ShouldBindJSON(&postDTO); err != nil {
		config.ValidateError(ctx, err)
		return
	}

	// Updating the post
	err := logic.UpdatePost(postIDParam.PostID, postDTO) // 修改此处传递postDTO的类型
	if err != nil {
		zap.L().Error("failed to update post", zap.Error(err))
		common.FailByMsg(ctx, "更新帖子失败")
		return
	}

	common.Success(ctx, nil)
}

// DeletePost deletes the specified post
func DeletePost(ctx *gin.Context) {
	postIDStr := ctx.Param("post_id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		zap.L().Error("failed to parse post_id", zap.Error(err))
		common.FailByMsg(ctx, "解析帖子ID失败")
		return
	}

	err = logic.DeletePost(postID)
	if err != nil {
		zap.L().Error("failed to delete post", zap.Error(err))
		common.FailByMsg(ctx, "删除帖子失败")
		return
	}

	common.Success(ctx, "帖子删除成功")
}
