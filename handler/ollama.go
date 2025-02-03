package handler

import (
	"io"
	"log/slog"
	"net/http"
	"safe-ollama/config"
	"safe-ollama/middleware"
	"safe-ollama/model"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func OllamaHandler(router *gin.Engine, db *gorm.DB) {
	r := router.Group("")
	chatRouter := r.Group("", middleware.OllamaAuth(db), middleware.OllamaTokenCount(db))
	chatRouter.POST("/api/generate", forwardRequest("/api/generate"))
	chatRouter.POST("/api/chat", forwardRequest("/api/chat"))
	chatRouter.POST("/api/chat-stream", forwardRequest("/api/chat-stream"))
	chatRouter.POST("/v1/chat/completions", forwardRequest("/v1/chat/completions"))
	chatRouter.POST("/v1/completions", forwardRequest("/v1/completions"))
	chatRouter.POST("/v1/embeddings", forwardRequest("/v1/embeddings"))

	ollamaRouter := r.Group("/api", middleware.LoginAuth(), middleware.RoleAuth([]string{model.ADMIN_ROLE}))
	ollamaRouter.POST("/create", forwardRequest("/api/create"))
	ollamaRouter.GET("/tags", forwardRequest("/api/tags"))
	ollamaRouter.POST("/show", forwardRequest("/api/show"))
	ollamaRouter.POST("/copy", forwardRequest("/api/copy"))
	ollamaRouter.POST("/pull", forwardRequest("/api/pull"))
	ollamaRouter.POST("/push", forwardRequest("/api/push"))
	ollamaRouter.POST("/embed", forwardRequest("/api/embed"))
	ollamaRouter.POST("/embeddings", forwardRequest("/api/embeddings"))
	ollamaRouter.GET("/ps", forwardRequest("/api/ps"))
	ollamaRouter.DELETE("/delete", forwardRequest("/api/delete"))
	ollamaRouter.GET("/version", forwardRequest("/api/version"))
}

var (
	httpClient = &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        200,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     time.Duration(config.OllamaTimeout) * time.Second,
		},
	}
	rateLimiter = make(chan struct{}, 200)
	mu          sync.Mutex
)

func forwardRequest(path string) func(c *gin.Context) {
	return func(c *gin.Context) {
		url := config.OllamaHost + path
		req, err := http.NewRequest(c.Request.Method, url, c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
			return
		}

		header := c.Request.Header.Clone()
		header.Del("Authorization")
		req.Header = header

		mu.Lock()
		select {
		case rateLimiter <- struct{}{}:
			defer func() {
				<-rateLimiter // 归还令牌
				mu.Unlock()
			}()
		default:
			mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			return
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to communicate with API"})
			return
		}
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)

		slog.Debug("[Ollama]", "url", req.URL, "status", resp.StatusCode, "headers", resp.Header)

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			c.JSON(resp.StatusCode, gin.H{"error": string(body)})
			return
		}

		for key, values := range resp.Header {
			for _, value := range values {
				c.Header(key, value)
			}
		}

		c.Status(resp.StatusCode)

		_, err = io.Copy(c.Writer, resp.Body)
		if err != nil {
			slog.Error("error during copying response body", "error", err)
			return
		}
	}
}
