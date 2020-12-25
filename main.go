package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-event-store/eventstore"
	"github.com/go-event-store/example/api"
	_ "github.com/go-event-store/example/docs"
	todo "github.com/go-event-store/example/internal"
	"github.com/go-event-store/example/internal/middleware"
	"github.com/go-event-store/example/internal/postgres"
	"github.com/go-event-store/pg"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	if _, err := os.Stat("./.env"); err == nil {
		viper.SetConfigName(".env")
		viper.SetConfigType("env")
		viper.AddConfigPath("./")

		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error in config file: %s", err))
		}
	}

	viper.AutomaticEnv()
}

// @title Todo Example Service
// @version 1.0
// @description Todo Example App for GO EventStore
// @contact.name Frank Jogeleit
// @contact.email frank.jogeleit@web.de
func main() {
	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, viper.GetString("DB_URL"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	es := eventstore.NewEventStore(pg.NewPersistenceStrategy(pool))
	err = es.Install(ctx)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = es.CreateStream(ctx, postgres.TodoStream)
	if err != nil {
		if err, ok := err.(eventstore.StreamAlreadyExist); ok == false {
			fmt.Println(err.Error())
			return
		}
	}

	es.AppendMiddleware(eventstore.Appended, middleware.EventLogger(eventstore.Appended))

	typeRegistry := eventstore.NewTypeRegistry()
	typeRegistry.RegisterAggregate(&todo.Todo{})
	typeRegistry.RegisterEvents(
		todo.TodoWasCreated{},
		todo.TodoWasUpdated{},
		todo.TodoWasDeleted{},
		todo.TodoWasDone{},
		todo.TodoWasUndone{},
	)

	go func(ctx context.Context, es *eventstore.EventStore, pool *pgxpool.Pool) {
		projector := postgres.NewTodoReadModelProjector(es, pool)
		err := projector.Run(ctx, true)

		fmt.Println(err.Error())
	}(ctx, es, pool)

	repo := postgres.NewTodoRepository(es)
	finder := postgres.NewTodoFinder(pool)

	todoHandler := api.NewTodoHandler(
		todo.NewCommandHandler(repo),
		todo.NewQueryHandler(finder),
	)

	r := gin.New()
	r.Use(gin.Recovery())

	url := ginSwagger.URL(viper.GetString("SWAGGER_DOC_JSON")) // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	r.GET("/todo", todoHandler.ListHandler)
	r.GET("/todo/:id", todoHandler.GetHandler)

	r.POST("/create-todo", todoHandler.CreateHandler)
	r.POST("/update-todo", todoHandler.UpdateHandler)
	r.POST("/delete-todo", todoHandler.DeleteHandler)
	r.POST("/do-todo", todoHandler.DoHandler)
	r.POST("/undo-todo", todoHandler.UndoHandler)
	r.Run()
}
