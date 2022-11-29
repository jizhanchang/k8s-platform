package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"k8s-platfrom/logic"
)

func GetPodsHandler(c *gin.Context) {

	//1.验证参数

	params := new(struct {
		FilterName string `from:"filterName"`
		NameSpace  string `from:"namespace" binding:"required"`
		Limit      int    `from:"limit"`
		Page       int    `from:"page"`
	})
	fmt.Printf("%T %v\n", params, params)
	if err := c.ShouldBind(params); err != nil {
		zap.L().Error("GetPodsHandler with invalid params", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			zap.L().Error("invalid request params", zap.Error(err))
			ResponseError(c, CodeInvalidParams)
			return
		}

		zap.L().Error("invalid request params", zap.Error(errs))
		ResponseErrorWithMsg(c, CodeInvalidParams, removeTopStruct(errs.Translate(trans)))
		return
	}
	zap.L().Info("bind from data success", zap.Any("params", params))
	//2.交个logic 层处理

	pods, err := logic.Pod.GetPods(params.FilterName, params.NameSpace, params.Limit, params.Page)
	if err != nil {
		zap.L().Error("logic get pods failed", zap.Error(err))
		ResponseError(c, CodeSeverBusy)
		return
	}
	ResponseSuccess(c, pods)
}
