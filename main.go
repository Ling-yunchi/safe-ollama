package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/samber/slog-gin"
	"io/fs"
	"net/http"
	"safe-ollama/config"
	"safe-ollama/handler"
	"safe-ollama/model"
	"safe-ollama/utils"
)

//go:embed dist/*
var dist embed.FS

func main() {
	config.ReadConfig()
	logger := utils.InitLogger()

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(sloggin.New(logger))
	r.Use(ServerStatic("dist", dist))

	db := model.InitDB()

	handler.UserHandler(r, db)
	handler.AuthHandler(r, db)
	handler.OllamaHandler(r, db)
	handler.OllamaTokenHandler(r, db)
	handler.OllamaTokenUsageHandler(r, db)

	r.NoRoute(func(c *gin.Context) {
		fsys, err := fs.Sub(dist, "dist")
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			return
		}
		http.ServeFileFS(c.Writer, c.Request, fsys, "index.html")
	})

	err := r.Run(config.ServerAddr)
	if err != nil {
		panic(err)
	} else {
		logger.Info("Server started a", "addr", config.ServerAddr)
	}
}

func ServerStatic(prefix string, embedFs embed.FS) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 去掉前缀
		fsys, err := fs.Sub(embedFs, prefix)
		if err != nil {
			panic(err)
		}
		fs2 := http.FS(fsys)
		f, err := fs2.Open(c.Request.URL.Path)
		if err != nil {
			c.Next()
			return
		}
		defer func(f http.File) {
			_ = f.Close()
		}(f)
		http.FileServer(fs2).ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}
