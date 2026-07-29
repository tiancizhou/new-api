package controller

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/stretchr/testify/assert"
)

func TestEmployeeTokenUsageResponseUsesUSD(t *testing.T) {
	stat := model.EmployeeTokenUsageStat{
		RequestCount:     2,
		Quota:            int64(common.QuotaPerUnit) + 60535,
		PromptTokens:     3200,
		CompletionTokens: 1800,
		TotalTokens:      5000,
	}

	response := newEmployeeTokenUsageResponse(stat)

	assert.Equal(t, 2, int(response.RequestCount))
	assert.InDelta(t, 1.12107, response.Quota, 0.000000001)
	assert.Equal(t, int64(3200), response.PromptTokens)
	assert.Equal(t, int64(1800), response.CompletionTokens)
	assert.Equal(t, int64(5000), response.TotalTokens)
}
