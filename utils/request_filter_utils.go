package utils

import (
	"fmt"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/gin-gonic/gin"
)

func EvaluateFilter(c *gin.Context, filter string) (bool, error) {
	adjustedFilter := strings.ReplaceAll(filter, "@request.", "request.")

	variables := buildVariables(c)

	program, err := expr.Compile(adjustedFilter, expr.Env(variables))
	if err != nil {
		return false, fmt.Errorf("failed to compile filter: %v", err)
	}

	output, err := expr.Run(program, variables)
	if err != nil {
		return false, fmt.Errorf("failed to evaluate filter: %v", err)
	}

	result, ok := output.(bool)
	if !ok {
		return false, fmt.Errorf("filter result is not a boolean")
	}

	return result, nil
}

func buildVariables(c *gin.Context) map[string]interface{} {
	contextValue, exists := c.Get("request.context")

	if !exists {
		contextValue = "default"
	}
	headers := make(map[string]string)
	for key, values := range c.Request.Header {
		normalizedKey := strings.ToLower(strings.ReplaceAll(key, "-", "_"))
		if len(values) > 0 {
			headers[normalizedKey] = values[0]
		}
	}
	query := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			query[key] = values[0]
		}
	}
	authData, exists := c.Get("auth")
	if !exists {
		authData = make(map[string]interface{})
	}
	body := extractBodyData(c)
	return map[string]interface{}{
		"request": map[string]interface{}{
			"context": contextValue,
			"method":  c.Request.Method,
			"headers": headers,
			"query":   query,
			"auth":    authData,
			"body":    body,
		},
	}
}

func extractBodyData(c *gin.Context) map[string]interface{} {
	if body, exists := c.Get("body"); exists {
		if data, ok := body.(map[string]interface{}); ok {
			return data
		}
	}
	body := make(map[string]interface{})
	if c.Request != nil && c.Request.PostForm != nil {
		for key, values := range c.Request.PostForm {
			if len(values) > 0 {
				body[key] = values[0]
			}
		}
	}
	return body
}
