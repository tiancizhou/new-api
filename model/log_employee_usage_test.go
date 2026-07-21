package model

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupEmployeeTokenUsageTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	common.SetDatabaseTypes(common.DatabaseTypeSQLite, common.DatabaseTypeSQLite)
	common.RedisEnabled = false

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	DB = db
	LOG_DB = db
	require.NoError(t, db.AutoMigrate(&Log{}))

	t.Cleanup(func() {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	})

	return db
}

func TestSumEmployeeTokenUsageFiltersByTokenAndEmployee(t *testing.T) {
	db := setupEmployeeTokenUsageTestDB(t)
	start := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC).Unix()

	logs := []Log{
		{CreatedAt: start + 10, Type: LogTypeConsume, TokenId: 11, EmployeeNo: "zt64003", Quota: 100, PromptTokens: 20, CompletionTokens: 30},
		{CreatedAt: start + 20, Type: LogTypeConsume, TokenId: 11, EmployeeNo: "zt64003", Quota: 200, PromptTokens: 40, CompletionTokens: 50},
		{CreatedAt: start + 30, Type: LogTypeConsume, TokenId: 11, EmployeeNo: "zt64004", Quota: 999, PromptTokens: 1, CompletionTokens: 1},
		{CreatedAt: start + 40, Type: LogTypeConsume, TokenId: 12, EmployeeNo: "zt64003", Quota: 999, PromptTokens: 1, CompletionTokens: 1},
		{CreatedAt: start + 50, Type: LogTypeError, TokenId: 11, EmployeeNo: "zt64003", Quota: 999, PromptTokens: 1, CompletionTokens: 1},
	}
	require.NoError(t, db.Create(&logs).Error)

	stat, err := SumEmployeeTokenUsage(11, "zt64003", start, start+3600)
	require.NoError(t, err)

	assert.EqualValues(t, 2, stat.RequestCount)
	assert.EqualValues(t, 300, stat.Quota)
	assert.EqualValues(t, 60, stat.PromptTokens)
	assert.EqualValues(t, 80, stat.CompletionTokens)
	assert.EqualValues(t, 140, stat.TotalTokens)
}

func TestGetEmployeeTokenUsageBucketsUsesDailyBoundaries(t *testing.T) {
	db := setupEmployeeTokenUsageTestDB(t)
	day1 := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	day2 := day1.AddDate(0, 0, 1)

	logs := []Log{
		{CreatedAt: day1.Add(1 * time.Hour).Unix(), Type: LogTypeConsume, TokenId: 11, EmployeeNo: "zt64003", Quota: 100, PromptTokens: 10, CompletionTokens: 20},
		{CreatedAt: day2.Add(1 * time.Hour).Unix(), Type: LogTypeConsume, TokenId: 11, EmployeeNo: "zt64003", Quota: 200, PromptTokens: 30, CompletionTokens: 40},
	}
	require.NoError(t, db.Create(&logs).Error)

	buckets, err := GetEmployeeTokenUsageBuckets(11, "zt64003", []time.Time{day1, day2}, day2.AddDate(0, 0, 1))
	require.NoError(t, err)
	require.Len(t, buckets, 2)

	assert.Equal(t, "2026-07-01", buckets[0].Label)
	assert.EqualValues(t, 100, buckets[0].Usage.Quota)
	assert.EqualValues(t, 30, buckets[0].Usage.TotalTokens)
	assert.Equal(t, "2026-07-02", buckets[1].Label)
	assert.EqualValues(t, 200, buckets[1].Usage.Quota)
	assert.EqualValues(t, 70, buckets[1].Usage.TotalTokens)
}
