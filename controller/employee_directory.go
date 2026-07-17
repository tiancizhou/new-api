package controller

import (
	"strconv"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/service"

	"github.com/gin-gonic/gin"
)

func GetDepartmentDirectoryTree(c *gin.Context) {
	tree, err := service.GetDepartmentTree()
	if err != nil {
		common.ApiError(c, err)
		return
	}
	common.ApiSuccess(c, tree)
}

func GetEmployeeDirectory(c *gin.Context) {
	pageInfo := common.GetPageQuery(c)
	departmentId, _ := strconv.Atoi(c.Query("department_id"))
	includeSubDepartments := c.DefaultQuery("include_subdepartments", "true") != "false"
	status, _ := strconv.Atoi(c.Query("status"))
	keyword := strings.TrimSpace(c.Query("keyword"))

	var departmentIds []int
	if departmentId == -1 {
		departmentIds = []int{0}
	} else if departmentId > 0 {
		departmentIds = []int{departmentId}
		if includeSubDepartments {
			var err error
			departmentIds, err = model.GetDepartmentSubtreeIds(departmentId)
			if err != nil {
				common.ApiError(c, err)
				return
			}
		}
	}
	employees, total, err := model.ListEmployees(departmentIds, keyword, status, pageInfo.GetStartIdx(), pageInfo.GetPageSize())
	if err != nil {
		common.ApiError(c, err)
		return
	}
	pageInfo.SetItems(employees)
	pageInfo.SetTotal(int(total))
	common.ApiSuccess(c, pageInfo)
}

func SyncMdmEmployeeInfo(c *gin.Context) {
	var req dto.MdmUserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ApiError(c, err)
		return
	}
	result, err := service.SyncMdmEmployee(req)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	common.ApiSuccess(c, result)
}

func SyncMdmDepartmentInfo(c *gin.Context) {
	var req dto.MdmDeptInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ApiError(c, err)
		return
	}
	result, err := service.SyncMdmDepartment(req)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	common.ApiSuccess(c, result)
}
