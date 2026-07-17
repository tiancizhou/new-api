package service

import (
	"testing"

	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetEmployeeDirectory(t *testing.T) {
	t.Helper()
	require.NoError(t, model.DB.Exec("DELETE FROM employees").Error)
	require.NoError(t, model.DB.Exec("DELETE FROM departments").Error)
	t.Cleanup(func() {
		model.DB.Exec("DELETE FROM employees")
		model.DB.Exec("DELETE FROM departments")
	})
}

func TestMasterDataSyncBuildsReadOnlyEmployeeDirectory(t *testing.T) {
	resetEmployeeDirectory(t)

	root, err := SyncMdmDepartment(dto.MdmDeptInfoRequest{
		Id:       dto.MdmID("company"),
		DeptName: "集团总部",
	})
	require.NoError(t, err)
	assert.Equal(t, "created", root.Action)

	department, err := SyncMdmDepartment(dto.MdmDeptInfoRequest{
		Id:       dto.MdmID("technology"),
		ParentId: dto.MdmID("company"),
		DeptName: "技术中心",
	})
	require.NoError(t, err)
	assert.Equal(t, "created", department.Action)

	employee, err := SyncMdmEmployee(dto.MdmUserInfoRequest{
		Id:       dto.MdmID("64003"),
		Account:  "ZT64003",
		RealName: "张三",
		DeptId:   dto.MdmID("technology"),
	})
	require.NoError(t, err)
	assert.Equal(t, "created", employee.Action)
	assert.Equal(t, "zt64003", employee.EmployeeNo)
	assert.Equal(t, department.DepartmentId, employee.DepartmentId)

	tree, err := GetDepartmentTree()
	require.NoError(t, err)
	require.Len(t, tree, 1)
	require.Len(t, tree[0].Children, 1)
	assert.EqualValues(t, 1, tree[0].TotalEmployees)
	assert.EqualValues(t, 1, tree[0].Children[0].EmployeeCount)

	rows, total, err := model.ListEmployees([]int{department.DepartmentId}, "64003", model.DirectoryStatusEnabled, 0, 10)
	require.NoError(t, err)
	assert.EqualValues(t, 1, total)
	require.Len(t, rows, 1)
	assert.Equal(t, "技术中心", rows[0].DepartmentName)
}

func TestMasterDataSyncRejectsInvalidEmployeeNumber(t *testing.T) {
	resetEmployeeDirectory(t)

	_, err := SyncMdmEmployee(dto.MdmUserInfoRequest{
		Id:       dto.MdmID("64003"),
		Account:  "64003",
		RealName: "张三",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "zt64003")
}

func TestMasterDataSyncAcceptsAlphanumericEmployeeNumber(t *testing.T) {
	resetEmployeeDirectory(t)

	result, err := SyncMdmEmployee(dto.MdmUserInfoRequest{
		Id:       dto.MdmID("C3543"),
		Account:  "ztC3543",
		RealName: "程湛恩",
	})
	require.NoError(t, err)
	assert.Equal(t, "ztc3543", result.EmployeeNo)
}

func TestEmployeeDirectoryReturnsEmptyArrayForDepartmentWithoutEmployees(t *testing.T) {
	resetEmployeeDirectory(t)

	department, err := SyncMdmDepartment(dto.MdmDeptInfoRequest{
		Id:       dto.MdmID("empty-department"),
		DeptName: "空部门",
	})
	require.NoError(t, err)

	rows, total, err := model.ListEmployees(
		[]int{department.DepartmentId},
		"",
		model.DirectoryStatusEnabled,
		0,
		20,
	)
	require.NoError(t, err)
	assert.NotNil(t, rows)
	assert.Empty(t, rows)
	assert.Zero(t, total)
}
