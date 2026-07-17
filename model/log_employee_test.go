package model

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestConsumeLogsCanBeFilteredAndSummedByEmployee(t *testing.T) {
	originalDB := DB
	originalLogDB := LOG_DB
	originalLogType := common.LogDatabaseType()
	t.Cleanup(func() {
		DB = originalDB
		LOG_DB = originalLogDB
		common.SetLogDatabaseType(originalLogType)
	})

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&User{}, &Log{}))
	DB = db
	LOG_DB = db
	common.SetLogDatabaseType(common.DatabaseTypeSQLite)
	require.NoError(t, db.Create(&User{
		Id:       1,
		Username: "system-owner",
		Password: "unused",
		Status:   common.UserStatusEnabled,
	}).Error)

	for _, employeeNo := range []string{"10086", "10010"} {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)
		ctx.Set("username", "system-owner")
		common.SetContextKey(ctx, constant.ContextKeyEmployeeNo, employeeNo)

		RecordConsumeLog(ctx, 1, RecordConsumeLogParams{
			ModelName:        "gpt-test",
			TokenName:        "agent-system",
			Quota:            100,
			PromptTokens:     10,
			CompletionTokens: 5,
		})
	}

	logs, total, err := GetUserLogs(1, LogTypeConsume, 0, 0, "", "10086", "", 0, 20, "", "", "")
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, logs, 1)
	assert.Equal(t, "10086", logs[0].EmployeeNo)

	stat, err := SumUsedQuota(LogTypeConsume, 0, 0, "", "system-owner", "10086", "", 0, "")
	require.NoError(t, err)
	assert.Equal(t, 100, stat.Quota)
}
