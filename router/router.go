package router

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"

	sess "github.com/hackathon-21-spring-02/back-end/session"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	s            sess.Session
	clientID     = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
)

// パッケージの初期化
func init() {
	if clientID == "" {
		panic(errors.New("clientID should not be empty."))
	}
	if clientSecret == "" {
		panic(errors.New("clientSecret should not be empty."))
	}
}

func SetRouting(sess sess.Session) {
	s = sess

	e := echo.New()

	e.Use(session.Middleware(sess.Store()))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	proxyConfig := middleware.DefaultProxyConfig
	clientURL, err := url.Parse("https://hackathon21_spring_02.trap.show/front-end/")
	if err != nil {
		panic(err)
	}
	proxyConfig.Balancer = middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
		{
			URL: clientURL,
		},
	})

	// if env == "development" || env == "mock" {
	// 	e.Pre(middleware.Rewrite(map[string]string{
	// 		"/customtheme-server/*": "/$1",
	// 	}))
	// }
	proxyConfig.Skipper = func(c echo.Context) bool {
		if strings.HasPrefix(c.Path(), "/api/") || strings.HasPrefix(c.Path(), "/openapi/") {
			return true
		}
		c.Request().Host = "hackathon21_spring_02.trap.show"
		return false
	}
	proxyConfig.ModifyResponse = func(res *http.Response) error {
		res.Header.Set("Cache-Control", "max-age=3600")
		return nil
	}
	proxyConfig.Rewrite = map[string]string{
		"/users*": "/",
		"/files":    "/",
		"/favorite": "/",
	}

	e.Use(middleware.ProxyWithConfig(proxyConfig))

	e.Static("/openapi", "docs/swagger")

	api := e.Group("/api")
	{
		apiPing := api.Group("/ping")
		{
			apiPing.GET("", func(c echo.Context) error {
				return echo.NewHTTPError(http.StatusOK, "pong!")
			})
		}

		// ユーザー情報
		apiUsers := api.Group("/users")
		{
			apiUsers.GET("", getUsersHandler, userAuthMiddleware)
		}

		apiFiles := api.Group("/files")
		{
			apiFiles.GET("", getFilesHandler, userAuthMiddleware)
			apiFiles.GET("/:fileID", getFileHandler, userAuthMiddleware)
			apiFiles.GET("/:fileID/download", getFileDownloadHandler, userAuthMiddleware)
		}

		// OAuth関連
		apiOAuth := api.Group("/oauth")
		{
			apiOAuth.GET("/callback", callbackHandler)
			apiOAuth.GET("/generate/code", postGenerateCodeHandler)
			apiOAuth.POST("/logout", postLogoutHandler, userAuthMiddleware)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	}
	err := e.Start(port)
	if err != nil {
		panic(err)
	}
}
