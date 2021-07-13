package router

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/hackathon-21-spring-02/back-end/model"
	"github.com/labstack/echo/v4"
)

func convertError(err error) error {
	if errors.Is(err, model.ErrNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Not Found")
	} else if errors.Is(err, model.ErrNoChange) {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "No Change")
	} else {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Something Wrong: %w", err))
	}
}
