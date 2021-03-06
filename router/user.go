package router

import (
	"net/http"

	"github.com/hackathon-21-spring-02/back-end/model"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// getUsersHandler GET /users
func getUsersHandler(c echo.Context) error {
	ctx := c.Request().Context()
	users, err := model.GetUsers(ctx)
	if err != nil {
		return generateEchoError(err)
	}

	return echo.NewHTTPError(http.StatusOK, users)
}

// getUserHandler GET /users/:userID
func getUserHandler(c echo.Context) error {
	ctx := c.Request().Context()
	userID := c.Param("userID")
	sess, err := session.Get("sessions", c)
	if err != nil {
		return errSessionNotFound(err)
	}
	accessToken := sess.Values["accessToken"].(string)
	user, err := model.GetUser(ctx, accessToken, userID)
	if err != nil {
		return generateEchoError(err)
	}

	return echo.NewHTTPError(http.StatusOK, user)
}

// getUsersMeHandler GET /users/me
func getUsersMeHandler(c echo.Context) error {
	ctx := c.Request().Context()
	sess, err := session.Get("sessions", c)
	if err != nil {
		return errSessionNotFound(err)
	}
	accessToken := sess.Values["accessToken"].(string)
	myUserID := sess.Values["id"].(string)
	res, err := model.GetUsersMe(ctx, accessToken, myUserID)
	if err != nil {
		return generateEchoError(err)
	}

	return echo.NewHTTPError(http.StatusOK, res)
}

// getUsersMeFavoritesHandler /users/me/favorites
func getUsersMeFavoritesHandler(c echo.Context) error {
	ctx := c.Request().Context()
	sess, err := session.Get("sessions", c)
	if err != nil {
		return errSessionNotFound(err)
	}
	accessToken := sess.Values["accessToken"].(string)
	userID := sess.Values["id"].(string)
	res, err := model.GetUsersMeFavorites(ctx, accessToken, userID)
	if err != nil {
		return generateEchoError(err)
	}

	return echo.NewHTTPError(http.StatusOK, res)
}
