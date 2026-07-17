package dto

import (
	"strconv"
	"strings"

	"github.com/QuantumNous/new-api/common"
)

type MdmID string

func (id *MdmID) UnmarshalJSON(data []byte) error {
	raw := strings.TrimSpace(string(data))
	if raw == "" || raw == "null" {
		*id = ""
		return nil
	}
	if strings.HasPrefix(raw, `"`) {
		value, err := strconv.Unquote(raw)
		if err != nil {
			return err
		}
		*id = MdmID(strings.TrimSpace(value))
		return nil
	}
	*id = MdmID(raw)
	return nil
}

func (id MdmID) String() string {
	return strings.TrimSpace(string(id))
}

type MdmDeptInfoRequest struct {
	Id                MdmID  `json:"id"`
	ParentId          MdmID  `json:"parentId"`
	ParentIdSnake     MdmID  `json:"parent_id"`
	DeptName          string `json:"deptName"`
	DeptNameSnake     string `json:"dept_name"`
	FullName          string `json:"fullName"`
	FullNameSnake     string `json:"full_name"`
	Ancestors         string `json:"ancestors"`
	AncestorsUpper    string `json:"ANCESTORS"`
	DeptCategory      *int   `json:"deptCategory"`
	DeptCategorySnake *int   `json:"dept_category"`
	Sort              *int   `json:"sort"`
	Remark            string `json:"remark"`
	Extract           string `json:"extract"`
	OrgCode           MdmID  `json:"orgCode"`
	OrgCodeSnake      MdmID  `json:"org_code"`
	OrgName           string `json:"orgName"`
	OrgNameSnake      string `json:"org_name"`
}

type mdmDeptInfoRequestAlias MdmDeptInfoRequest

type mdmDeptNRDEnvelope struct {
	Type string         `json:"TYPE"`
	Data mdmDeptNRDData `json:"DATA"`
}

type mdmDeptNRDData struct {
	ParentDeptCode MdmID  `json:"PDEPTCODE"`
	DeptCode       MdmID  `json:"DEPTCODE"`
	DeptName       string `json:"DEPTNAME"`
	OrgCode        MdmID  `json:"ORGCODE"`
	OrgName        string `json:"ORGNAME"`
	Extract        string `json:"EXTRACT"`
	Remark         string `json:"OBLIGATE1"`
	Sort           *int   `json:"SORT"`
}

func (r *MdmDeptInfoRequest) UnmarshalJSON(data []byte) error {
	var flat mdmDeptInfoRequestAlias
	if err := common.Unmarshal(data, &flat); err != nil {
		return err
	}
	*r = MdmDeptInfoRequest(flat)

	var envelope mdmDeptNRDEnvelope
	if err := common.Unmarshal(data, &envelope); err != nil {
		return err
	}
	if !strings.EqualFold(strings.TrimSpace(envelope.Type), "DEPT") || envelope.Data.DeptCode.String() == "" {
		return nil
	}

	r.Id = buildMdmDeptExternalId(envelope.Data.OrgCode, envelope.Data.DeptCode)
	r.ParentId = buildMdmParentDeptExternalId(envelope.Data.OrgCode, envelope.Data.ParentDeptCode)
	r.DeptName = strings.TrimSpace(envelope.Data.DeptName)
	r.FullName = r.DeptName
	r.OrgCode = envelope.Data.OrgCode
	r.OrgName = strings.TrimSpace(envelope.Data.OrgName)
	r.Extract = strings.TrimSpace(envelope.Data.Extract)
	r.Remark = strings.TrimSpace(envelope.Data.Remark)
	r.Sort = envelope.Data.Sort
	return nil
}

func (r MdmDeptInfoRequest) ExternalId() string { return r.Id.String() }

func (r MdmDeptInfoRequest) ParentExternalId() string {
	if value := r.ParentId.String(); value != "" {
		return value
	}
	return r.ParentIdSnake.String()
}

func (r MdmDeptInfoRequest) Name() string {
	for _, value := range []string{r.DeptName, r.DeptNameSnake, r.FullName, r.FullNameSnake} {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
	return ""
}

func (r MdmDeptInfoRequest) FullDisplayName() string {
	if value := strings.TrimSpace(r.FullName); value != "" {
		return value
	}
	return strings.TrimSpace(r.FullNameSnake)
}

func (r MdmDeptInfoRequest) ExternalAncestors() string {
	if value := strings.TrimSpace(r.Ancestors); value != "" {
		return value
	}
	return strings.TrimSpace(r.AncestorsUpper)
}

func (r MdmDeptInfoRequest) Category() int {
	if r.DeptCategory != nil {
		return *r.DeptCategory
	}
	if r.DeptCategorySnake != nil {
		return *r.DeptCategorySnake
	}
	return 0
}

func (r MdmDeptInfoRequest) SortValue() int {
	if r.Sort != nil {
		return *r.Sort
	}
	return 0
}

func (r MdmDeptInfoRequest) CompanyCode() string {
	if value := r.OrgCode.String(); value != "" {
		return value
	}
	return r.OrgCodeSnake.String()
}

func (r MdmDeptInfoRequest) CompanyName() string {
	if value := strings.TrimSpace(r.OrgName); value != "" {
		return value
	}
	return strings.TrimSpace(r.OrgNameSnake)
}

func (r MdmDeptInfoRequest) IsDelete() bool {
	return strings.EqualFold(strings.TrimSpace(r.Extract), "D")
}

type MdmUserInfoRequest struct {
	Id            MdmID  `json:"id"`
	Account       string `json:"account"`
	Name          string `json:"name"`
	RealName      string `json:"realName"`
	RealNameSnake string `json:"real_name"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	DeptId        MdmID  `json:"deptId"`
	DeptIdSnake   MdmID  `json:"dept_id"`
	PostId        string `json:"postId"`
	PostIdSnake   string `json:"post_id"`
	Sex           *int   `json:"sex"`
	Extract       string `json:"extract"`
}

type mdmUserInfoRequestAlias MdmUserInfoRequest

type mdmUserNRDEnvelope struct {
	Type string         `json:"TYPE"`
	Data mdmUserNRDData `json:"DATA"`
}

type mdmUserNRDData struct {
	UserCode  MdmID  `json:"PSNCODE"`
	Account   string `json:"ZTPSNCODE"`
	Name      string `json:"PSNNAME"`
	Email     string `json:"EMAIL"`
	Phone     string `json:"TELEPHONE"`
	WorkPhone string `json:"WORKPHONE"`
	DeptCode  MdmID  `json:"DEPTCODE"`
	OrgCode   MdmID  `json:"ORGCODE"`
	PostCode  string `json:"POSTCODE"`
	Extract   string `json:"EXTRACT"`
	Sex       string `json:"SEX"`
}

func (r *MdmUserInfoRequest) UnmarshalJSON(data []byte) error {
	var flat mdmUserInfoRequestAlias
	if err := common.Unmarshal(data, &flat); err != nil {
		return err
	}
	*r = MdmUserInfoRequest(flat)

	var envelope mdmUserNRDEnvelope
	if err := common.Unmarshal(data, &envelope); err != nil {
		return err
	}
	if !strings.EqualFold(strings.TrimSpace(envelope.Type), "PSN") || envelope.Data.UserCode.String() == "" {
		return nil
	}

	r.Id = envelope.Data.UserCode
	r.Account = strings.TrimSpace(envelope.Data.Account)
	r.Name = strings.TrimSpace(envelope.Data.Name)
	r.RealName = r.Name
	r.Email = strings.TrimSpace(envelope.Data.Email)
	r.Phone = firstNonEmpty(envelope.Data.Phone, envelope.Data.WorkPhone)
	r.DeptId = buildMdmDeptExternalId(envelope.Data.OrgCode, envelope.Data.DeptCode)
	r.PostId = strings.TrimSpace(envelope.Data.PostCode)
	r.Extract = strings.TrimSpace(envelope.Data.Extract)
	r.Sex = parseMdmSex(envelope.Data.Sex)
	return nil
}

func (r MdmUserInfoRequest) ExternalId() string { return r.Id.String() }
func (r MdmUserInfoRequest) EmployeeNo() string { return strings.ToLower(strings.TrimSpace(r.Account)) }

func (r MdmUserInfoRequest) DisplayName() string {
	for _, value := range []string{r.RealName, r.RealNameSnake, r.Name} {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
	return ""
}

func (r MdmUserInfoRequest) DeptExternalId() string {
	value := r.DeptId.String()
	if value == "" {
		value = r.DeptIdSnake.String()
	}
	for _, part := range strings.Split(value, ",") {
		if part = strings.TrimSpace(part); part != "" {
			return part
		}
	}
	return ""
}

func (r MdmUserInfoRequest) Post() string {
	if value := strings.TrimSpace(r.PostId); value != "" {
		return value
	}
	return strings.TrimSpace(r.PostIdSnake)
}

func (r MdmUserInfoRequest) IsDelete() bool {
	return strings.EqualFold(strings.TrimSpace(r.Extract), "D")
}

func buildMdmDeptExternalId(orgCode MdmID, deptCode MdmID) MdmID {
	org := orgCode.String()
	dept := deptCode.String()
	if dept == "" {
		return ""
	}
	if org == "" || strings.HasPrefix(dept, org+org) {
		return MdmID(dept)
	}
	return MdmID(org + dept)
}

func buildMdmParentDeptExternalId(orgCode MdmID, parentDeptCode MdmID) MdmID {
	org := orgCode.String()
	parent := parentDeptCode.String()
	if parent == "" || parent == "0" {
		return MdmID(parent)
	}
	if org == "" || strings.HasPrefix(parent, org+org) {
		return MdmID(parent)
	}
	return MdmID(org + parent)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
	return ""
}

func parseMdmSex(value string) *int {
	switch strings.TrimSpace(value) {
	case "男", "M", "m", "1":
		sex := 1
		return &sex
	case "女", "F", "f", "2":
		sex := 2
		return &sex
	default:
		return nil
	}
}
