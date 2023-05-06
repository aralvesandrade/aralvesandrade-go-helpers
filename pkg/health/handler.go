package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler is a Health Handler implementation
type Handler struct {
	CompositeChecker
}

// NewHandler returns a new Handler
func NewHandler() Handler {
	return Handler{}
}

// Health returns a json encoded Health
// set the status to http.StatusServiceUnavailable if the check is down
func (h Handler) Health(c *gin.Context) {
	health := h.CompositeChecker.Check()

	if health.IsDown() {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	c.JSON(http.StatusOK, health)
}
