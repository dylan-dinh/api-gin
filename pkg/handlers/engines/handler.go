package engines

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dylan-dinh/api-gin/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/modl"
	"net/http"
	"strconv"
)

type Handler struct {
	Db *modl.DbMap
}

type Config struct {
	Ngine *gin.Engine
	Db    *modl.DbMap
}

func NewHandler(c *Config) {
	h := &Handler{
		Db: c.Db,
	}

	e := c.Ngine.Group("/api/engines")

	e.GET("/:id", h.GetById)
	e.GET("", h.GetAll)
	e.POST("", h.Post)
	e.DELETE("/:id", h.Delete)
	e.PUT("/:id", h.Put)
}

func (h *Handler) GetAll(c *gin.Context) {
	var ngines []model.Engine

	err := h.Db.Select(&ngines, `SELECT * FROM engines`)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "no engines found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"engines": ngines,
	})
}

func (h *Handler) GetById(c *gin.Context) {
	var ngine model.Engine
	var id string

	if id = c.Param("id"); id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "'id' of the ressource not provided",
		})
		return
	}

	err := h.Db.SelectOne(&ngine, `SELECT * FROM engines WHERE id = $1`, id)
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

func (h *Handler) Put(c *gin.Context) {
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
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	if err := c.Bind(&ngine); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"engine": err.Error(),
		})
		return
	}

	ngine.Id = atoi

	rows, err := h.Db.Update(&ngine)
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

func (h *Handler) Post(c *gin.Context) {
	var body model.Engine

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"engine": err.Error(),
		})
		return
	}

	err := h.Db.Insert(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"engine": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"engine": body,
	})
}

func (h *Handler) Delete(c *gin.Context) {
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

	rows, err := h.Db.Delete(&ngine)
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
