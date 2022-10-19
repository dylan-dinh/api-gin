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

func (h *Router) getAllEngines(c *gin.Context) {
	var ngines []model.Engine

	err := h.db.Select(&ngines, `SELECT * FROM engines`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"engines": ngines,
	})
}

func (h *Router) getEngineById(c *gin.Context) {
	var ngine model.Engine
	var id string

	if id = c.Param("id"); id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "'id' of the ressource not provided",
		})
		return
	}

	err := h.db.SelectOne(&ngine, `SELECT * FROM engines WHERE id = $1`, id)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("engine with id '%s' not found", id),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"engine": ngine,
	})
}

func (h *Router) putEngine(c *gin.Context) {
	var ngine model.Engine
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

	if err := c.BindJSON(&ngine); err != nil {
		fmt.Println("HEREFQWEFQWERFWERGQERGERG")
		c.JSON(http.StatusBadRequest, gin.H{
			"engine": err.Error(),
		})
		return
	}

	ngine.Id = atoi

	rows, err := h.db.Update(&ngine)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("engine with id '%s' not found", id),
		})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (h *Router) postEngine(c *gin.Context) {
	var ngine model.Engine

	if err := c.BindJSON(&ngine); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"engine": err.Error(),
		})
		return
	}

	err := h.db.Insert(&ngine)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"engine": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"engine": ngine,
	})
}

func (h *Router) deleteEngine(c *gin.Context) {
	var ngine model.Engine
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

	ngine.Id = atoi

	rows, err := h.db.Delete(&ngine)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("engine with id '%s' not found", id),
		})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
