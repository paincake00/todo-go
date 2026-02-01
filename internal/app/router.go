package app

import (
	"github.com/gin-gonic/gin"
	docs "github.com/paincake00/todo-go/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (app *App) InitRouter() *gin.Engine {
	router := gin.Default()

	// Swagger setup
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	v1 := router.Group("/api/v1")
	{
		tasks := v1.Group("/tasks")
		{
			tasks.POST("", app.taskHandler.CreateTask)
			tasks.GET("", app.taskHandler.GetAllTasks)
			tasks.GET("/:id", app.taskHandler.GetTaskById)
			tasks.PUT("/:id", app.taskHandler.UpdateTaskById)
			tasks.DELETE("/:id", app.taskHandler.DeleteTaskById)
		}
	}

	return router
}
