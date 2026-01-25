package app

import "github.com/gin-gonic/gin"

func (app *App) InitRouter() *gin.Engine {
	router := gin.Default()

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
