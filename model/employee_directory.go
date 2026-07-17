package model

import "gorm.io/gorm"

const (
	DirectoryStatusEnabled  = 1
	DirectoryStatusDisabled = 2
)

type Department struct {
	Id                int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name              string `json:"name" gorm:"type:varchar(128);not null"`
	ParentId          int    `json:"parent_id" gorm:"type:int;default:0;index"`
	Level             int    `json:"level" gorm:"type:int;not null;default:0"`
	Path              string `json:"path" gorm:"type:varchar(1024);default:''"`
	ExternalDeptId    string `json:"external_dept_id" gorm:"type:varchar(64);not null;uniqueIndex"`
	ExternalParentId  string `json:"external_parent_id,omitempty" gorm:"type:varchar(64);index"`
	ExternalAncestors string `json:"external_ancestors,omitempty" gorm:"type:varchar(1024)"`
	ExternalFullName  string `json:"external_full_name,omitempty" gorm:"type:varchar(255)"`
	ExternalCategory  int    `json:"external_category,omitempty" gorm:"type:int;default:0"`
	ExternalSort      int    `json:"external_sort,omitempty" gorm:"type:int;default:0"`
	Status            int    `json:"status" gorm:"type:int;default:1;index"`
	Remark            string `json:"remark,omitempty" gorm:"type:varchar(255)"`
	CreatedAt         int64  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         int64  `json:"updated_at" gorm:"autoUpdateTime"`
}

type Employee struct {
	Id                 int    `json:"id" gorm:"primaryKey;autoIncrement"`
	EmployeeNo         string `json:"employee_no" gorm:"type:varchar(32);not null;uniqueIndex"`
	ExternalEmployeeId string `json:"external_employee_id" gorm:"type:varchar(64);not null;uniqueIndex"`
	Name               string `json:"name" gorm:"type:varchar(128);not null"`
	Email              string `json:"email,omitempty" gorm:"type:varchar(128)"`
	Phone              string `json:"phone,omitempty" gorm:"type:varchar(32)"`
	DepartmentId       int    `json:"department_id" gorm:"type:int;default:0;index"`
	ExternalDeptId     string `json:"external_dept_id,omitempty" gorm:"type:varchar(64);index"`
	PostId             string `json:"post_id,omitempty" gorm:"type:varchar(64)"`
	Sex                int    `json:"sex,omitempty" gorm:"type:int;default:0"`
	Status             int    `json:"status" gorm:"type:int;default:1;index"`
	CreatedAt          int64  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          int64  `json:"updated_at" gorm:"autoUpdateTime"`
}

type EmployeeDirectoryRow struct {
	Employee
	DepartmentName string `json:"department_name"`
}

func ListEmployees(departmentIds []int, keyword string, status int, startIdx int, pageSize int) ([]EmployeeDirectoryRow, int64, error) {
	query := DB.Model(&Employee{})
	if len(departmentIds) > 0 {
		query = query.Where("employees.department_id IN ?", departmentIds)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("employees.employee_no LIKE ? OR employees.name LIKE ?", like, like)
	}
	if status > 0 {
		query = query.Where("employees.status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	employees := make([]EmployeeDirectoryRow, 0)
	err := query.
		Select("employees.*, COALESCE(departments.name, '') AS department_name").
		Joins("LEFT JOIN departments ON departments.id = employees.department_id").
		Order("employees.employee_no ASC").
		Offset(startIdx).
		Limit(pageSize).
		Scan(&employees).Error
	return employees, total, err
}

func GetDepartmentSubtreeIds(departmentId int) ([]int, error) {
	var department Department
	if err := DB.First(&department, departmentId).Error; err != nil {
		return nil, err
	}
	ids := []int{departmentId}
	var descendants []int
	if err := DB.Model(&Department{}).
		Where("path LIKE ?", department.Path+"%").
		Where("id <> ?", departmentId).
		Pluck("id", &descendants).Error; err != nil {
		return nil, err
	}
	return append(ids, descendants...), nil
}

func CountEmployeesByDepartment() (map[int]int64, error) {
	type countRow struct {
		DepartmentId int
		Count        int64
	}
	var rows []countRow
	if err := DB.Model(&Employee{}).
		Select("department_id, COUNT(*) AS count").
		Where("status = ?", DirectoryStatusEnabled).
		Group("department_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	counts := make(map[int]int64, len(rows))
	for _, row := range rows {
		counts[row.DepartmentId] = row.Count
	}
	return counts, nil
}

func FindDepartmentByExternalId(tx *gorm.DB, externalId string) (*Department, error) {
	var department Department
	if err := tx.Where("external_dept_id = ?", externalId).First(&department).Error; err != nil {
		return nil, err
	}
	return &department, nil
}
