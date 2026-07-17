package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/QuantumNous/new-api/common"

	"github.com/gin-gonic/gin"
)

func MdmSyncAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		expected := strings.TrimSpace(common.GetEnvOrDefaultString("MDM_SYNC_TOKEN", ""))
		if expected == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"success": false, "message": "主数据同步令牌未配置"})
			c.Abort()
			return
		}
		token := strings.TrimSpace(c.GetHeader("X-MDM-Token"))
		if token == "" {
			token = strings.TrimSpace(c.GetHeader("Authorization"))
			token = strings.TrimPrefix(token, "Bearer ")
			token = strings.TrimPrefix(token, "bearer ")
		}
		if len(token) != len(expected) || subtle.ConstantTimeCompare([]byte(token), []byte(expected)) != 1 {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "主数据同步令牌无效"})
			c.Abort()
			return
		}
		c.Next()
	}
}
