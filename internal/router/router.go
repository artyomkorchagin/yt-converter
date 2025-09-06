package router

import (
	"net/http"

	"github.com/artyomkorchagin/yt-converter/pkg/helpers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
}

func NewHandler(logger *zap.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) InitRouter() *gin.Engine {
	staticPath, htmlPath := helpers.GetStaticPath()
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Static("/static", staticPath)
	router.LoadHTMLGlob(htmlPath)

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	main := router.Group("/")
	{
		main.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", nil)
		})

		main.GET("/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
	}

	apiv1 := router.Group("/api/v1/")
	{
		apiv1.GET("/video/:id", h.wrap(h.getVideo))
	}
	h.logger.Info("Routes initialized")
	return router
}
