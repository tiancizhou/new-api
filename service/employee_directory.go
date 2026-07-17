package service

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/model"

	"gorm.io/gorm"
)

var employeeNumberPattern = regexp.MustCompile(`^zt[a-z0-9]{4,7}$`)

type DepartmentTreeNode struct {
	Id             int                  `json:"id"`
	Name           string               `json:"name"`
	ParentId       int                  `json:"parent_id"`
	Level          int                  `json:"level"`
	Status         int                  `json:"status"`
	ExternalDeptId string               `json:"external_dept_id"`
	EmployeeCount  int64                `json:"employee_count"`
	TotalEmployees int64                `json:"total_employees"`
	Children       []DepartmentTreeNode `json:"children,omitempty"`
}

type MdmDepartmentSyncResult struct {
	Action         string `json:"action"`
	DepartmentId   int    `json:"department_id,omitempty"`
	ExternalDeptId string `json:"external_dept_id,omitempty"`
}

type MdmEmployeeSyncResult struct {
	Action             string `json:"action"`
	EmployeeId         int    `json:"employee_id,omitempty"`
	EmployeeNo         string `json:"employee_no,omitempty"`
	ExternalEmployeeId string `json:"external_employee_id,omitempty"`
	DepartmentId       int    `json:"department_id,omitempty"`
}

func GetDepartmentTree() ([]DepartmentTreeNode, error) {
	var departments []model.Department
	if err := model.DB.
		Where("status = ?", model.DirectoryStatusEnabled).
		Order("external_sort ASC, id ASC").
		Find(&departments).Error; err != nil {
		return nil, err
	}
	counts, err := model.CountEmployeesByDepartment()
	if err != nil {
		return nil, err
	}

	children := make(map[int][]model.Department)
	for _, department := range departments {
		children[department.ParentId] = append(children[department.ParentId], department)
	}
	var build func(int) []DepartmentTreeNode
	build = func(parentId int) []DepartmentTreeNode {
		departments := children[parentId]
		sort.SliceStable(departments, func(i, j int) bool {
			if departments[i].ExternalSort == departments[j].ExternalSort {
				return departments[i].Id < departments[j].Id
			}
			return departments[i].ExternalSort < departments[j].ExternalSort
		})
		nodes := make([]DepartmentTreeNode, 0, len(departments))
		for _, department := range departments {
			node := DepartmentTreeNode{
				Id:             department.Id,
				Name:           department.Name,
				ParentId:       department.ParentId,
				Level:          department.Level,
				Status:         department.Status,
				ExternalDeptId: department.ExternalDeptId,
				EmployeeCount:  counts[department.Id],
			}
			node.Children = build(department.Id)
			node.TotalEmployees = node.EmployeeCount
			for _, child := range node.Children {
				node.TotalEmployees += child.TotalEmployees
			}
			nodes = append(nodes, node)
		}
		return nodes
	}
	roots := build(0)
	if counts[0] > 0 {
		roots = append(roots, DepartmentTreeNode{
			Id:             -1,
			Name:           "Unassigned",
			EmployeeCount:  counts[0],
			TotalEmployees: counts[0],
			Status:         model.DirectoryStatusEnabled,
		})
	}
	return roots, nil
}

func SyncMdmDepartment(req dto.MdmDeptInfoRequest) (*MdmDepartmentSyncResult, error) {
	externalId := req.ExternalId()
	if externalId == "" {
		return nil, errors.New("主数据部门ID不能为空")
	}
	result := &MdmDepartmentSyncResult{ExternalDeptId: externalId}

	err := model.DB.Transaction(func(tx *gorm.DB) error {
		if err := ensureMdmCompanyDepartment(tx, req); err != nil {
			return err
		}

		department, err := model.FindDepartmentByExternalId(tx, externalId)
		if req.IsDelete() {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				result.Action = "noop"
				return nil
			}
			if err != nil {
				return err
			}
			result.DepartmentId = department.Id
			result.Action = "disabled"
			return tx.Model(department).Update("status", model.DirectoryStatusDisabled).Error
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		parentId, level, parentPath, err := resolveMdmDepartmentParent(tx, req.ParentExternalId())
		if err != nil {
			return err
		}
		name := req.Name()
		if name == "" {
			name = externalId
		}
		fullName := req.FullDisplayName()
		if fullName == "" {
			fullName = name
		}

		if errors.Is(err, gorm.ErrRecordNotFound) || department == nil {
			department = &model.Department{
				Name:              name,
				ParentId:          parentId,
				Level:             level,
				Path:              parentPath,
				ExternalDeptId:    externalId,
				ExternalParentId:  req.ParentExternalId(),
				ExternalAncestors: req.ExternalAncestors(),
				ExternalFullName:  fullName,
				ExternalCategory:  req.Category(),
				ExternalSort:      req.SortValue(),
				Status:            model.DirectoryStatusEnabled,
				Remark:            strings.TrimSpace(req.Remark),
			}
			if err := tx.Create(department).Error; err != nil {
				return err
			}
			department.Path = buildDepartmentPath(parentPath, department.Id)
			if err := tx.Save(department).Error; err != nil {
				return err
			}
			result.Action = "created"
		} else {
			oldPath := department.Path
			department.Name = name
			department.ParentId = parentId
			department.Level = level
			department.Path = buildDepartmentPath(parentPath, department.Id)
			department.ExternalParentId = req.ParentExternalId()
			department.ExternalAncestors = req.ExternalAncestors()
			department.ExternalFullName = fullName
			department.ExternalCategory = req.Category()
			department.ExternalSort = req.SortValue()
			department.Status = model.DirectoryStatusEnabled
			department.Remark = strings.TrimSpace(req.Remark)
			if err := tx.Save(department).Error; err != nil {
				return err
			}
			if oldPath != "" && oldPath != department.Path {
				if err := updateDepartmentDescendantPaths(tx, department.Id, oldPath, department.Path); err != nil {
					return err
				}
			}
			result.Action = "updated"
		}
		result.DepartmentId = department.Id
		return reparentPendingDepartments(tx, department)
	})
	return result, err
}

func SyncMdmEmployee(req dto.MdmUserInfoRequest) (*MdmEmployeeSyncResult, error) {
	externalId := req.ExternalId()
	employeeNo := req.EmployeeNo()
	if externalId == "" && employeeNo == "" {
		return nil, errors.New("主数据人员ID或工号不能为空")
	}
	result := &MdmEmployeeSyncResult{
		ExternalEmployeeId: externalId,
		EmployeeNo:         employeeNo,
	}

	err := model.DB.Transaction(func(tx *gorm.DB) error {
		employee, found, err := findMdmEmployee(tx, externalId, employeeNo)
		if err != nil {
			return err
		}
		if req.IsDelete() {
			if !found {
				result.Action = "noop"
				return nil
			}
			result.EmployeeId = employee.Id
			result.DepartmentId = employee.DepartmentId
			result.EmployeeNo = employee.EmployeeNo
			result.Action = "disabled"
			return tx.Model(employee).Update("status", model.DirectoryStatusDisabled).Error
		}
		if !employeeNumberPattern.MatchString(employeeNo) {
			return errors.New("工号格式必须为zt加4至7位字母或数字，例如zt64003")
		}
		if externalId == "" {
			return errors.New("主数据人员ID不能为空")
		}
		name := req.DisplayName()
		if name == "" {
			return errors.New("主数据人员姓名不能为空")
		}

		externalDeptId := req.DeptExternalId()
		departmentId, err := resolveMdmEmployeeDepartment(tx, externalDeptId)
		if err != nil {
			return err
		}
		sex := 0
		if req.Sex != nil {
			sex = *req.Sex
		}
		if !found {
			employee = &model.Employee{
				EmployeeNo:         employeeNo,
				ExternalEmployeeId: externalId,
				Name:               name,
				Email:              strings.TrimSpace(req.Email),
				Phone:              strings.TrimSpace(req.Phone),
				DepartmentId:       departmentId,
				ExternalDeptId:     externalDeptId,
				PostId:             req.Post(),
				Sex:                sex,
				Status:             model.DirectoryStatusEnabled,
			}
			if err := tx.Create(employee).Error; err != nil {
				return err
			}
			result.Action = "created"
		} else {
			updates := map[string]interface{}{
				"employee_no":          employeeNo,
				"external_employee_id": externalId,
				"name":                 name,
				"email":                strings.TrimSpace(req.Email),
				"phone":                strings.TrimSpace(req.Phone),
				"department_id":        departmentId,
				"external_dept_id":     externalDeptId,
				"post_id":              req.Post(),
				"sex":                  sex,
				"status":               model.DirectoryStatusEnabled,
			}
			if err := tx.Model(employee).Updates(updates).Error; err != nil {
				return err
			}
			result.Action = "updated"
		}
		result.EmployeeId = employee.Id
		result.DepartmentId = departmentId
		return nil
	})
	return result, err
}

func ensureMdmCompanyDepartment(tx *gorm.DB, req dto.MdmDeptInfoRequest) error {
	companyCode := req.CompanyCode()
	if companyCode == "" {
		return nil
	}
	companyName := req.CompanyName()
	if companyName == "" {
		companyName = companyCode
	}
	department, err := model.FindDepartmentByExternalId(tx, companyCode)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		department = &model.Department{
			Name:             companyName,
			Path:             "/",
			ExternalDeptId:   companyCode,
			ExternalParentId: "0",
			ExternalFullName: companyName,
			Status:           model.DirectoryStatusEnabled,
		}
		if err := tx.Create(department).Error; err != nil {
			return err
		}
		department.Path = buildDepartmentPath("/", department.Id)
		return tx.Save(department).Error
	}
	if err != nil {
		return err
	}
	return tx.Model(department).Updates(map[string]interface{}{
		"name":               companyName,
		"status":             model.DirectoryStatusEnabled,
		"external_parent_id": "0",
		"external_full_name": companyName,
	}).Error
}

func resolveMdmDepartmentParent(tx *gorm.DB, externalParentId string) (int, int, string, error) {
	externalParentId = strings.TrimSpace(externalParentId)
	if externalParentId == "" || externalParentId == "0" {
		return 0, 0, "/", nil
	}
	parent, err := model.FindDepartmentByExternalId(tx, externalParentId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, 0, "/", nil
	}
	if err != nil {
		return 0, 0, "", err
	}
	return parent.Id, parent.Level + 1, parent.Path, nil
}

func resolveMdmEmployeeDepartment(tx *gorm.DB, externalDeptId string) (int, error) {
	if strings.TrimSpace(externalDeptId) == "" {
		return 0, nil
	}
	department, err := model.FindDepartmentByExternalId(tx, externalDeptId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return department.Id, nil
}

func findMdmEmployee(tx *gorm.DB, externalId string, employeeNo string) (*model.Employee, bool, error) {
	var employee model.Employee
	if employeeNo != "" {
		err := tx.Where("employee_no = ?", employeeNo).First(&employee).Error
		if err == nil {
			return &employee, true, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, err
		}
	}
	if externalId != "" {
		err := tx.Where("external_employee_id = ?", externalId).First(&employee).Error
		if err == nil {
			return &employee, true, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, err
		}
	}
	return nil, false, nil
}

func buildDepartmentPath(parentPath string, departmentId int) string {
	if parentPath == "" {
		parentPath = "/"
	}
	if !strings.HasSuffix(parentPath, "/") {
		parentPath += "/"
	}
	return fmt.Sprintf("%s%d/", parentPath, departmentId)
}

func updateDepartmentDescendantPaths(tx *gorm.DB, departmentId int, oldPrefix string, newPrefix string) error {
	var descendants []model.Department
	if err := tx.Where("path LIKE ? AND id <> ?", oldPrefix+"%", departmentId).Find(&descendants).Error; err != nil {
		return err
	}
	for _, descendant := range descendants {
		nextPath := strings.Replace(descendant.Path, oldPrefix, newPrefix, 1)
		level := len(strings.Split(strings.Trim(nextPath, "/"), "/")) - 1
		if err := tx.Model(&model.Department{}).Where("id = ?", descendant.Id).Updates(map[string]interface{}{
			"path":  nextPath,
			"level": level,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

func reparentPendingDepartments(tx *gorm.DB, parent *model.Department) error {
	var children []model.Department
	if err := tx.Where("external_parent_id = ? AND id <> ? AND parent_id <> ?", parent.ExternalDeptId, parent.Id, parent.Id).Find(&children).Error; err != nil {
		return err
	}
	for _, child := range children {
		oldPath := child.Path
		child.ParentId = parent.Id
		child.Level = parent.Level + 1
		child.Path = buildDepartmentPath(parent.Path, child.Id)
		if err := tx.Save(&child).Error; err != nil {
			return err
		}
		if oldPath != "" && oldPath != child.Path {
			if err := updateDepartmentDescendantPaths(tx, child.Id, oldPath, child.Path); err != nil {
				return err
			}
		}
		if err := reparentPendingDepartments(tx, &child); err != nil {
			return err
		}
	}
	return nil
}
