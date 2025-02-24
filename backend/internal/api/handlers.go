package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stephen/panopticon/internal/db"
	"github.com/stephen/panopticon/internal/models"
	"go.opentelemetry.io/otel"
)

type Handler struct {
	db *db.Database
}

func NewHandler(db *db.Database) *Handler {
	return &Handler{db: db}
}

func (h *Handler) GetTasks(c *gin.Context) {
	ctx := c.Request.Context()
	tracer := otel.Tracer("api")
	ctx, span := tracer.Start(ctx, "GetTasks")
	defer span.End()

	tasks, err := h.db.GetTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update metrics
	var active, completed int
	for _, task := range tasks {
		if task.Completed {
			completed++
		} else {
			active++
		}
	}
	activeTasksGauge.Set(float64(active))
	completedTasksGauge.Set(float64(completed))

	taskOperationsTotal.WithLabelValues("get").Inc()
	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) CreateTask(c *gin.Context) {
	ctx := c.Request.Context()
	tracer := otel.Tracer("api")
	ctx, span := tracer.Start(ctx, "CreateTask")
	defer span.End()

	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.db.CreateTask(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update metrics
	taskOperationsTotal.WithLabelValues("create").Inc()
	if !task.Completed {
		activeTasksGauge.Inc()
	} else {
		completedTasksGauge.Inc()
	}

	c.JSON(http.StatusCreated, task)
}

func (h *Handler) UpdateTask(c *gin.Context) {
	ctx := c.Request.Context()
	tracer := otel.Tracer("api")
	ctx, span := tracer.Start(ctx, "UpdateTask")
	defer span.End()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.UpdateTask(id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update metrics
	taskOperationsTotal.WithLabelValues("update").Inc()

	c.Status(http.StatusOK)
}

func (h *Handler) DeleteTask(c *gin.Context) {
	ctx := c.Request.Context()
	tracer := otel.Tracer("api")
	ctx, span := tracer.Start(ctx, "DeleteTask")
	defer span.End()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	if err := h.db.DeleteTask(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update metrics
	taskOperationsTotal.WithLabelValues("delete").Inc()

	c.Status(http.StatusOK)
}
