package main

// go get -u github.com/swaggo/swag/cmd/swag
// go install github.com/swaggo/swag/cmd/swag
import (
	"backend/config"
	githubRoute "backend/internal/pkg/github/route/http"
	userRoute "backend/internal/pkg/user/route/http"

	"context"
	"fmt"
	"net/http"
	"time"

	_ "backend/docs"

	swagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
)

func NewEcho() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	return e

}

func NewApp() *fx.App {
	return fx.New(
		fx.Provide(
			config.New,
			NewEcho,
		),
		fx.Invoke(
			userRoute.NewUserHandler,
			githubRoute.NewGithubHandler,
			serve,
		),
	)
}

func serve(lifecycle fx.Lifecycle, echo *echo.Echo, cfg config.Config) {
	echo.GET("/swagger/*", swagger.WrapHandler)

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println(cfg.Listen)

			go func() {
				if err := echo.Start(cfg.Listen); err != nil {
					panic(fmt.Sprintf("echo.Start() err=%v", err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Stopping monitoring server.")
			return nil
		},
	})
}

// @title worklist Sample Swagger API
// @version 1.0
// @host localhost:1323
// @BasePath /
func main() {
	app := NewApp()

	startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.Start(startCtx); err != nil {
		fmt.Println(err)
	}

	<-app.Done()

}

// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": "Server is up and running",
	})
}
