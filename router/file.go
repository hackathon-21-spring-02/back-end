package router

import (
	"fmt"
	"net/http"

	"github.com/hackathon-21-spring-02/back-end/model"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// getFileDownloadHandler GET /files/:fileID/download
func getFileDownloadHandler(c echo.Context) error {
	ctx := c.Request().Context()
	fileID := c.Param("fileID")

	sess, err := session.Get("sessions", c)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Errorf("failed to get session: %w", err).Error())
	}
	accessToken := sess.Values["accessToken"].(string)

	res, err := model.GetFileDownload(ctx, fileID, accessToken)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Errorf("failed to get file: %w", err).Error())
	}

	return c.Stream(http.StatusOK, res.Header.Get("Content-Type"), res.Body)
}
