package config

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func NewCORS(v *viper.Viper) gin.HandlerFunc {
	allowOrigins := parseAllowOrigins(v)
	allowCredentials := v.GetBool("CORS_ALLOW_CREDENTIALS")
	if len(allowOrigins) == 1 && allowOrigins[0] == "*" && allowCredentials {
		allowCredentials = false
	}

	return cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept", "X-Requested-With", "Access-Control-Request-Method", "Access-Control-Request-Headers", "X-CSRF-Token", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "X-Requested-With", "X-CSRF-Token", "Authorization"},
		AllowCredentials: allowCredentials,
		MaxAge:           24 * time.Hour,
	})
}

func parseAllowOrigins(v *viper.Viper) []string {
	raw := strings.TrimSpace(v.GetString("CORS_ALLOW_ORIGINS"))
	if raw != "" {
		return splitAndClean(raw)
	}

	origins := v.GetStringSlice("CORS_ALLOW_ORIGINS")
	if len(origins) == 0 {
		return []string{"*"}
	}

	return splitAndClean(strings.Join(origins, ","))
}

func splitAndClean(input string) []string {
	parts := strings.Split(input, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	if len(result) == 0 {
		return []string{"*"}
	}
	return result
}
