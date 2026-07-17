package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupContextForEmployee(t *testing.T) {
	tests := []struct {
		name           string
		token          *model.Token
		employeeNo     string
		wantStatus     int
		wantEmployeeNo string
	}{
		{name: "ordinary token does not require employee", token: &model.Token{}, wantStatus: http.StatusOK},
		{name: "required employee is missing", token: &model.Token{RequireEmployee: true}, wantStatus: http.StatusBadRequest},
		{name: "employee identifier is recorded without user lookup", token: &model.Token{RequireEmployee: true}, employeeNo: " external-10086 ", wantStatus: http.StatusOK, wantEmployeeNo: "external-10086"},
		{name: "oversized employee identifier is rejected", token: &model.Token{RequireEmployee: true}, employeeNo: strings.Repeat("1", 65), wantStatus: http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)
			if tt.employeeNo != "" {
				ctx.Request.Header.Set(constant.EmployeeNoHeader, tt.employeeNo)
			}

			err := SetupContextForEmployee(ctx, tt.token)
			if tt.wantStatus == http.StatusOK {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
			assert.Equal(t, tt.wantStatus, recorder.Code)
			assert.Equal(t, tt.wantEmployeeNo, common.GetContextKeyString(ctx, constant.ContextKeyEmployeeNo))
			if tt.wantEmployeeNo != "" {
				assert.Empty(t, ctx.Request.Header.Get(constant.EmployeeNoHeader))
			}
		})
	}
}
