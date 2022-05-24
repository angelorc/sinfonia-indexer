package server

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	c "github.com/angelorc/sinfonia-indexer/config"
	"github.com/angelorc/sinfonia-indexer/dataloader"
	"github.com/angelorc/sinfonia-indexer/generated/gqlgen"
	"github.com/angelorc/sinfonia-indexer/resolver"
	"github.com/angelorc/sinfonia-indexer/utility"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func New() {
	// Setup & configure server
	// more info -> https://echo.labstack.com/
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: utility.LoggerFormat()}))
	e.HideBanner = true

	// Load routes from graphql
	// InitGraphql(e)

	// Load routes from rest
	InitRest(e)

	// Start server routes
	go func() {
		if err := e.Start(":" + "9090"); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Stop server gracefully
	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func InitRest(e *echo.Echo) {
	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "<3")
	})
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
}

var token string
var playgroundPassword string
var submissionToken string

// Get header value and add to gql resolvers
func getHeaders(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		headers := ctx.Request().Header
		token = utility.GetHeaderString("Authorization", headers)
		playgroundPassword = utility.GetHeaderString("Playground-Password", headers)

		return next(ctx)
	}
}

func InitGraphql(e *echo.Echo) {
	// Resolvers && Directives
	resolver := resolver.Resolver{Token: &token}
	config := gqlgen.Config{Resolvers: &resolver}

	/*config.Directives.Auth = func(ctx context.Context, obj interface{}, next graphql.Resolver, role []guard.Role) (interface{}, error) {
		if err := guard.Auth(role, *resolver.Token); err != nil {
			return nil, fmt.Errorf("Access denied")
		}
		return next(ctx)
	}*/
	e.Use(getHeaders)

	// new custom handler based on gqlgen version 0.11.3
	queryHandler := handler.New(gqlgen.NewExecutableSchema(config))

	// queryHandler.Use(&debug.Tracer{})
	queryHandler.AddTransport(transport.POST{})
	queryHandler.AddTransport(transport.MultipartForm{})
	queryHandler.SetQueryCache(lru.New(1000))
	queryHandler.Use(extension.AutomaticPersistedQuery{Cache: lru.New(100)})
	queryHandler.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		rc := graphql.GetOperationContext(ctx)
		if playgroundPassword == c.GetSecret("GRAPHQL_PLAYGROUND_PASS") {
			rc.DisableIntrospection = false
		} else {
			rc.DisableIntrospection = true
		}
		return next(ctx)
	})

	e.GET("/", echo.WrapHandler(playground.Handler("GraphQL Playground", c.GetSecret("GRAPHQL_ENDPOINT"))))
	e.POST("/query", echo.WrapHandler(dataloader.DataLoaderMiddleware(queryHandler)))
}
