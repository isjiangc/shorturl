package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	nethttp "net/http"
	"nunu_ginblog/docs"
	"nunu_ginblog/internal/handler"
	"nunu_ginblog/internal/middleware"
	"nunu_ginblog/pkg/jwt"
	"nunu_ginblog/pkg/log"
	"nunu_ginblog/pkg/server/http"
)

func NewHTTPServer(
	logger *log.Logger,
	conf *viper.Viper,
	jwt *jwt.JWT,
	userHandler *handler.UserHandler,
	shorturlHandler *handler.ShorturlHandler,
) *http.Server {
	gin.SetMode(gin.DebugMode)
	s := http.NewServer(
		gin.Default(),
		logger,
		http.WithServerHost(conf.GetString("http.host")),
		http.WithServerPort(conf.GetInt("http.port")),
	)
	s.Static("static", "web/static")
	s.LoadHTMLGlob("web/*.html")
	// swagger doc
	docs.SwaggerInfo.BasePath = "/v1"
	s.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerfiles.Handler,
		//ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", conf.GetInt("app.http.port"))),
		ginSwagger.DefaultModelsExpandDepth(-1),
	))
	s.Use(
		middleware.CORSMiddleware(),
		middleware.ResponseLogMiddleware(logger),
		middleware.RequestLogMiddleware(logger),
		//middleware.SignMiddleware(log),
	)

	s.GET("/:url", shorturlHandler.ShortUrlDetail)

	s.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(nethttp.StatusNotFound, "404.html", gin.H{
			"message": "你访问的页面已失效",
			"code":    nethttp.StatusNotFound,
		})
	})

	v1 := s.Group("/v1")
	{
		// No route group has permission
		noAuthRouter := v1.Group("/")
		{
			noAuthRouter.POST("/register", userHandler.Register)
			noAuthRouter.POST("/login", userHandler.Login)
			noAuthRouter.POST("/shorturl", shorturlHandler.GenShortUrl)
			noAuthRouter.PUT("/shorturl/:url", shorturlHandler.UpdateUrlState)
			noAuthRouter.DELETE("/shorturl/:url", shorturlHandler.DeleteShortUrl)
			noAuthRouter.GET("/shorturl", shorturlHandler.GetShortUrlList)
			noAuthRouter.GET("/shorturl/:url", shorturlHandler.GetShortUrlInfo)

		}
		// Non-strict permission routing group
		noStrictAuthRouter := v1.Group("/").Use(middleware.NoStrictAuth(jwt, logger))
		{
			noStrictAuthRouter.GET("/user", userHandler.GetProfile)
		}

		// Strict permission routing group
		strictAuthRouter := v1.Group("/").Use(middleware.StrictAuth(jwt, logger))
		{
			strictAuthRouter.PUT("/user", userHandler.UpdateProfile)
		}
	}

	return s
}
