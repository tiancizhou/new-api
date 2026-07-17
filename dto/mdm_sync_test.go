package dto

import (
	"testing"

	"github.com/QuantumNous/new-api/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMdmPersonnelEnvelopeMapsEmployeeNumber(t *testing.T) {
	payload := []byte(`{
		"DATA": {
			"PSNCODE": "64003",
			"ZTPSNCODE": "zt64003",
			"PSNNAME": "张三",
			"DEPTCODE": "A10103101",
			"ORGCODE": "A101",
			"EXTRACT": "R"
		},
		"TYPE": "PSN"
	}`)

	var req MdmUserInfoRequest
	require.NoError(t, common.Unmarshal(payload, &req))
	assert.Equal(t, "64003", req.ExternalId())
	assert.Equal(t, "zt64003", req.EmployeeNo())
	assert.Equal(t, "张三", req.DisplayName())
	assert.Equal(t, "A101A10103101", req.DeptExternalId())
}

func TestMdmDepartmentEnvelopeMapsHierarchy(t *testing.T) {
	payload := []byte(`{
		"DATA": {
			"PDEPTCODE": "A101031",
			"DEPTCODE": "A10103101",
			"ORGCODE": "A101",
			"ORGNAME": "中天钢铁集团",
			"DEPTNAME": "综合科",
			"EXTRACT": "N"
		},
		"TYPE": "DEPT"
	}`)

	var req MdmDeptInfoRequest
	require.NoError(t, common.Unmarshal(payload, &req))
	assert.Equal(t, "A101A10103101", req.ExternalId())
	assert.Equal(t, "A101A101031", req.ParentExternalId())
	assert.Equal(t, "综合科", req.Name())
	assert.Equal(t, "中天钢铁集团", req.CompanyName())
}
