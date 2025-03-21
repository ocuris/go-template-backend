package https

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ocuris/go-template-backend/internals/modules/healthz"
	"github.com/ocuris/go-template-backend/internals/utils/ctx"
	"github.com/ocuris/go-template-backend/internals/utils/logger"
)

var log = logger.NewLogger()

type healthHandler struct {
	usecase healthz.Usecase
}

func NewHealthzHandler(e *echo.Echo, usecase healthz.Usecase) {
	handler := &healthHandler{
		usecase: usecase,
	}

	api := e.Group("/api/v1")
	api.GET("/healthz", handler.CheckHealth)
}

func (h *healthHandler) CheckHealth(c echo.Context) error {
	ac := c.(*ctx.CustomApplicationContext)

	// Call the health check use case
	status, err := h.usecase.CheckHealth()
	if err != nil {
		log.Errorf("Health check failed: %v", err)
		return ac.CustomResponse("error", nil, "Service is unhealthy", err.Error(), http.StatusInternalServerError, nil)
	}

	log.Info("Health check passed")
	return ac.CustomResponse("success", status, "Service is healthy", "", http.StatusOK, nil)
}
