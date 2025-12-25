package middleware

import (
	"go-expense-management-system/internal/messages"
	"go-expense-management-system/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/ulule/limiter/v3"
	memorystore "github.com/ulule/limiter/v3/drivers/store/memory"
)

func NewRateLimiter(cfg *viper.Viper) gin.HandlerFunc {
	rateStr := cfg.GetString("RATE_LIMIT")
	rate, err := limiter.NewRateFromFormatted(rateStr)
	if err != nil {
		panic(err)
	}

	store := memorystore.NewStore()
	limiterInstance := limiter.New(store, rate)
	excludePaths := parseExcludePaths(cfg.GetString("RATE_LIMIT_EXCLUDE_PATHS"))

	return func(c *gin.Context) {
		if isPathExcluded(c.Request.URL.Path, excludePaths) {
			c.Next()
			return
		}

		identifier := c.ClientIP()
		limitCtx, err := limiterInstance.Get(c, identifier)
		if err != nil {
			res := utils.FailedResponse(messages.InternalServerError)
			c.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}

		if limitCtx.Reached {
			res := utils.FailedResponse(messages.TooManyRequests)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, res)
			return
		}

		c.Next()
	}
}

func parseExcludePaths(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	paths := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			paths = append(paths, trimmed)
		}
	}
	return paths
}

func isPathExcluded(path string, patterns []string) bool {
	if len(patterns) == 0 {
		return false
	}

	for _, pattern := range patterns {
		if pattern == path {
			return true
		}

		if strings.HasSuffix(pattern, "/*") {
			prefix := strings.TrimSuffix(pattern, "/*")
			if path == prefix || strings.HasPrefix(path, prefix+"/") {
				return true
			}
		}
	}

	return false
}
