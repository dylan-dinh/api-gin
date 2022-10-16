package engines

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/modl"
	"net/http"
)

type Handler struct {
	Db *modl.DbMap
}

type Config struct {
	Ngine *gin.Engine
	Db    *modl.DbMap
}

func NewHandler(c *Config) {
	h := &Handler{}

	e := c.Ngine.Group("/api/engine")

	e.GET("/:name", h.Get)
}

func (h *Handler) Get(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"engine": c.Param("name"),
	})
}
