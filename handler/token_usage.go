package handler

import (
	"net/http"
	"safe-ollama/middleware"
	"safe-ollama/model"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserUsageResponse struct {
	Model          string `json:"model"`
	APICallCount   int    `json:"api_call_count"`
	PromptTokens   int    `json:"prompt_tokens"`
	ResponseTokens int    `json:"response_tokens"`
}

type AdminUsageResponse struct {
	UserID        uint   `json:"user_id"`
	Model         string `json:"model"`
	TotalPrompt   int    `json:"total_prompt"`
	TotalResponse int    `json:"total_response"`
}

func OllamaTokenUsageHandler(router *gin.Engine, db *gorm.DB) {
	r := router.Group("/api/token_usage", middleware.LoginAuth())
	userRoutes := r.Group("/user", middleware.RoleAuth([]string{model.USER_ROLE, model.ADMIN_ROLE}))
	userRoutes.GET("/usage/daily", GetUserDailyUsage(db))
	userRoutes.GET("/usage", GetUserUsage(db))
	userRoutes.GET("/usage/:model", GetUserModelUsage(db))

	adminRoute := r.Group("/admin", middleware.RoleAuth([]string{model.ADMIN_ROLE}))
	adminRoute.GET("/usage", GetSystemUsage(db))
}

// GET /user/usage/daily?start=2024-01-01&end=2024-01-07
func GetUserDailyUsage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		userID := claims.(model.JwtPayload).UserId
		start, end := parseTimeRange(c)
		end = end.AddDate(0, 0, 1) // 将结束时间设置为下一天的开始时间

		var results []struct {
			Date           string `json:"date"`
			PromptTokens   int    `json:"promptTokens"`
			ResponseTokens int    `json:"responseTokens"`
		}

		err := db.Model(&model.TokenUsage{}).
			Select("DATE(time) as date, SUM(prompt_eval_count) as prompt_tokens, SUM(eval_count) as response_tokens").
			Where("user_id = ? AND time BETWEEN ? AND ?", userID, start, end).
			Group("date").
			Find(&results).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取使用数据"})
			return
		}

		c.JSON(http.StatusOK, results)
	}
}

// GetUserUsage GET /user/usage?start=2024-01-01&end=2024-01-31
func GetUserUsage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("currentUserID").(uint)
		start, end := parseTimeRange(c)

		var results []UserUsageResponse
		query := db.Model(&model.TokenUsage{}).
			Select("ollama_model as model, "+
				"count(*) as api_call_count, "+
				"sum(prompt_eval_count) as prompt_tokens, "+
				"sum(eval_count) as response_tokens").
			Where("user_id = ? AND time BETWEEN ? AND ?", userID, start, end).
			Group("ollama_model")

		if err := query.Find(&results).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取使用数据"})
			return
		}

		c.JSON(http.StatusOK, results)
	}
}

// GetUserModelUsage 用户指定模型使用详情
// GET /user/usage/llama3.2?start=2024-01-01&end=2024-01-31
func GetUserModelUsage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("currentUserID").(uint)
		_model := c.Param("model")
		start, end := parseTimeRange(c)

		var result UserUsageResponse
		err := db.Model(&model.TokenUsage{}).
			Select("ollama_model as model, "+
				"count(*) as api_call_count, "+
				"sum(prompt_eval_count) as prompt_tokens, "+
				"sum(eval_count) as response_tokens").
			Where("user_id = ? AND ollama_model = ? AND time BETWEEN ? AND ?",
				userID, _model, start, end).
			First(&result).Error

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "未找到使用记录"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// GetSystemUsage 管理员查看全系统使用统计
// GET /admin/usage?model=llama3.2&user_id=123&start=2024-01-01&end=2024-01-31
func GetSystemUsage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filter struct {
			Model  string `form:"model"`
			UserID uint   `form:"user_id"`
		}
		err := c.ShouldBindQuery(&filter)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
			return
		}
		start, end := parseTimeRange(c)

		query := db.Model(&model.TokenUsage{})
		if filter.Model != "" {
			query = query.Where("ollama_model = ?", filter.Model)
		}
		if filter.UserID > 0 {
			query = query.Where("user_id = ?", filter.UserID)
		}

		var results []AdminUsageResponse
		err = query.
			Select("user_id, ollama_model as model, "+
				"sum(prompt_eval_count) as total_prompt, "+
				"sum(eval_count) as total_response").
			Where("time BETWEEN ? AND ?", start, end).
			Group("user_id, ollama_model").
			Find(&results).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取系统使用数据"})
			return
		}

		c.JSON(http.StatusOK, results)
	}
}

// 辅助函数：解析时间范围
func parseTimeRange(c *gin.Context) (start, end time.Time) {
	defaultStart := time.Now().AddDate(0, -1, 0) // 默认最近一个月
	defaultEnd := time.Now()

	start, _ = time.Parse("2006-01-02", c.Query("start"))
	end, _ = time.Parse("2006-01-02", c.Query("end"))

	if start.IsZero() {
		start = defaultStart
	}
	if end.IsZero() {
		end = defaultEnd
	}
	return
}
