package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dylan-dinh/api-gin/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Router) getAllSites(c *gin.Context) {
	var sites []model.Site

	err := h.db.Select(&sites, `SELECT * FROM sites`)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "no sites found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sites": sites,
	})
}

func (h *Router) getSiteById(c *gin.Context) {
	var site model.Site
	var id string

	if id = c.Param("id"); id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "'id' of the ressource not provided",
		})
		return
	}

	err := h.db.SelectOne(&site, `SELECT * FROM sites WHERE id = $1`, id)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("sites with id '%s' not found", id),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sites": site,
	})
}

func (h *Router) putSite(c *gin.Context) {
	var site model.Site
	var id string

	if id = c.Param("id"); id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "'id' of the ressource not provided",
		})
		return
	}

	atoi, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	if err := c.BindJSON(&site); err != nil {
		fmt.Println("weflkihbwefgbjkwefbjkwebjkf")
		c.JSON(http.StatusBadRequest, gin.H{
			"sites": err.Error(),
		})
		return
	}

	site.Id = atoi

	rows, err := h.db.Update(&site)
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("sites with id '%s' not found", id),
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (h *Router) postSite(c *gin.Context) {
	var site model.Site

	if err := c.BindJSON(&site); err != nil {
		c.JSON(http.StatusCreated, gin.H{
			"sites": err.Error(),
		})
		return
	}

	err := h.db.Insert(&site)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"sites": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"sites": site,
	})
}

func (h *Router) deleteSite(c *gin.Context) {
	var site model.Site
	var id string

	if id = c.Param("id"); id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "'id' of the ressource not provided",
		})
		return
	}

	atoi, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	site.Id = atoi

	rows, err := h.db.Delete(&site)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("sites with id '%s' not found", id),
		})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
