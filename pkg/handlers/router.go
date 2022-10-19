package handlers

import (
	"github.com/dylan-dinh/api-gin/pkg/database"
	"github.com/gin-gonic/gin"
)

type Router struct {
	*gin.Engine
	db *database.DB
}

func NewRouter(db *database.DB) *Router {
	r := &Router{Engine: gin.Default(), db: db}

	ngine := r.Group("/api/engines")

	// private method
	ngine.GET("/:id", r.getEngineById)
	ngine.GET("", r.getAllEngines)
	ngine.POST("", r.postEngine)
	ngine.DELETE("/:id", r.deleteEngine)
	ngine.PUT("/:id", r.putEngine)

	site := r.Group("/api/sites")

	site.GET("/:id", r.getSiteById)
	site.GET("", r.getAllSites)
	site.POST("", r.postSite)
	site.DELETE("/:id", r.deleteSite)
	site.PUT("/:id", r.putSite)

	return r
}
