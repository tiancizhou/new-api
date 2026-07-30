package middleware

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGlobalAPIRateLimitExcludesEmployeeUsagePath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	originalEnabled := common.GlobalApiRateLimitEnable
	originalNum := common.GlobalApiRateLimitNum
	originalDuration := common.GlobalApiRateLimitDuration
	originalRedisEnabled := common.RedisEnabled
	t.Cleanup(func() {
		common.GlobalApiRateLimitEnable = originalEnabled
		common.GlobalApiRateLimitNum = originalNum
		common.GlobalApiRateLimitDuration = originalDuration
		common.RedisEnabled = originalRedisEnabled
	})
	common.GlobalApiRateLimitEnable = true
	common.GlobalApiRateLimitNum = 1
	common.GlobalApiRateLimitDuration = 60
	common.RedisEnabled = false

	router := gin.New()
	router.Use(GlobalAPIRateLimit("/api/usage/token/employee"))
	router.GET("/api/usage/token/employee", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	for range 2 {
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/usage/token/employee", nil))
		require.Equal(t, http.StatusOK, response.Code)
	}
}

func TestEmployeeUsageRateLimitUsesTokenID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	originalEnabled := common.EmployeeUsageRateLimitEnable
	originalNum := common.EmployeeUsageRateLimitNum
	originalDuration := common.EmployeeUsageRateLimitDuration
	originalRedisEnabled := common.RedisEnabled
	t.Cleanup(func() {
		common.EmployeeUsageRateLimitEnable = originalEnabled
		common.EmployeeUsageRateLimitNum = originalNum
		common.EmployeeUsageRateLimitDuration = originalDuration
		common.RedisEnabled = originalRedisEnabled
	})
	common.EmployeeUsageRateLimitEnable = true
	common.EmployeeUsageRateLimitNum = 1
	common.EmployeeUsageRateLimitDuration = 60
	common.RedisEnabled = false

	tokenID := int(time.Now().UnixNano() % 1_000_000_000)
	router := gin.New()
	router.GET("/api/usage/token/employee", func(c *gin.Context) {
		requestTokenID, err := strconv.Atoi(c.Query("token_id"))
		require.NoError(t, err)
		c.Set("token_id", requestTokenID)
		c.Next()
	}, EmployeeUsageRateLimit(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	request := func(id int) *httptest.ResponseRecorder {
		response := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/usage/token/employee?token_id="+strconv.Itoa(id), nil)
		router.ServeHTTP(response, req)
		return response
	}

	first := request(tokenID)
	require.Equal(t, http.StatusOK, first.Code)
	second := request(tokenID)
	assert.Equal(t, http.StatusTooManyRequests, second.Code)
	third := request(tokenID + 1)
	assert.Equal(t, http.StatusOK, third.Code)
}
