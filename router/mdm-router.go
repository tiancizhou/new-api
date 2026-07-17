package router

import (
	"github.com/QuantumNous/new-api/controller"
	"github.com/QuantumNous/new-api/middleware"

	"github.com/gin-gonic/gin"
)

func registerMdmRoutes(group *gin.RouterGroup) {
	group.Use(middleware.MdmSyncAuth())
	group.POST("/userInfo", controller.SyncMdmEmployeeInfo)
	group.POST("/deptInfo", controller.SyncMdmDepartmentInfo)
}
