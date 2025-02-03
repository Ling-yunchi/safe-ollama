package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"safe-ollama/model"
	"strings"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func OllamaTokenCount(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = blw
		c.Next()
		if c.Writer.Status() != http.StatusOK {
			return
		}

		isStream := strings.Contains(c.Writer.Header().Get("Content-Type"), "application/x-ndjson")
		var data struct {
			Model           string `json:"model"`
			PromptEvalCount int    `json:"prompt_eval_count"`
			EvalCount       int    `json:"eval_count"`
		}

		if isStream {
			chunks := bytes.Split(blw.body.Bytes(), []byte("\n"))
			if len(chunks) < 2 {
				slog.Error("[Ollama Token] chunks length less than 2")
				return
			}
			chunk := chunks[len(chunks)-2]
			slog.Debug("[Ollama Token]", "body", string(chunk))
			if err := json.Unmarshal(chunk, &data); err != nil {
				slog.Error("[Ollama Token] fail to parse response body", "error", err)
				return
			}
		} else {
			if err := json.Unmarshal(blw.body.Bytes(), &data); err != nil {
				slog.Error("[Ollama Token] fail to parse response body", "error", err)
				return
			}
		}

		obj, ok := c.Get("ollamaToken")
		if !ok {
			slog.Error("[Ollama Token] fail to get ollama token")
			return
		}
		token := obj.(model.OllamaToken)

		if data.Model != "" && (data.PromptEvalCount > 0 || data.EvalCount > 0) {
			go func() {
				tokenUsage := model.TokenUsage{
					UserId:          token.UserId,
					OllamaModel:     data.Model,
					PromptEvalCount: data.PromptEvalCount,
					EvalCount:       data.EvalCount,
				}
				if err := db.Create(&tokenUsage).Error; err != nil {
					slog.Error("[Ollama Token] fail to create token usage", "error", err)
				}
			}()
		}
	}
}
