package router

import (
	"github.com/G4L1L10/admin-dashboard-backend/internal/app/handler"
	"github.com/G4L1L10/admin-dashboard-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	CourseHandler   *handler.CourseHandler
	LessonHandler   *handler.LessonHandler
	QuestionHandler *handler.QuestionHandler
	OptionHandler   *handler.OptionHandler
	TagHandler      *handler.TagHandler
	StatsHandler    *handler.StatsHandler
}

func NewRouter(
	courseHandler *handler.CourseHandler,
	lessonHandler *handler.LessonHandler,
	questionHandler *handler.QuestionHandler,
	optionHandler *handler.OptionHandler,
	tagHandler *handler.TagHandler,
	statsHandler *handler.StatsHandler,
) *Router {
	return &Router{
		CourseHandler:   courseHandler,
		LessonHandler:   lessonHandler,
		QuestionHandler: questionHandler,
		OptionHandler:   optionHandler,
		TagHandler:      tagHandler,
		StatsHandler:    statsHandler,
	}
}

func (r *Router) SetupRouter() *gin.Engine {
	router := gin.New()

	// Global middlewares
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.CORSMiddleware())

	api := router.Group("/api")
	{
		// ===== Courses =====
		courses := api.Group("/courses")
		{
			courses.GET("", r.CourseHandler.ListCourses)
			courses.GET("/:id", r.CourseHandler.GetCourse)

			// Protected mutations
			courses.Use(middleware.AuthMiddleware())
			{
				courses.POST("", r.CourseHandler.CreateCourse)
				courses.PUT("/:id", r.CourseHandler.UpdateCourse)
				courses.DELETE("/:id", r.CourseHandler.DeleteCourse)
			}
		}

		// ===== Lessons =====
		lesson := api.Group("/lessons")
		{
			lesson.GET("/detail/:id", r.LessonHandler.GetLesson)
			lesson.GET("/full/:id", r.LessonHandler.GetFullLesson)
			lesson.GET("/by-course/:course_id", r.LessonHandler.GetLessonsByCourseID)
			lesson.GET("/:lesson_id/questions", r.QuestionHandler.GetQuestionsByLesson)

			// Protected mutations
			lesson.Use(middleware.AuthMiddleware())
			{
				lesson.POST("", r.LessonHandler.CreateLesson)
				lesson.PUT("/:id", r.LessonHandler.UpdateLesson)
				lesson.DELETE("/:id", r.LessonHandler.DeleteLesson)
			}
		}

		// ===== Questions =====
		questions := api.Group("/questions")
		{
			questions.GET("/:id", r.QuestionHandler.GetQuestion)
			questions.GET("", r.QuestionHandler.GetQuestionsByTag)

			// Protected mutations
			questions.Use(middleware.AuthMiddleware())
			{
				questions.POST("", r.QuestionHandler.CreateQuestion)
				questions.PUT("/:id", r.QuestionHandler.UpdateQuestion)
				questions.DELETE("/:id", r.QuestionHandler.DeleteQuestion)
			}
		}

		// ===== Options =====
		options := api.Group("/options")
		{
			options.GET("/:id", r.OptionHandler.GetOption)

			// Protected mutations
			options.Use(middleware.AuthMiddleware())
			{
				options.POST("", r.OptionHandler.CreateOption)
				options.PUT("/:id", r.OptionHandler.UpdateOption)
				options.DELETE("/:id", r.OptionHandler.DeleteOption)
			}
		}

		// ===== Tags =====
		tags := api.Group("/tags")
		{
			tags.GET("/:id", r.TagHandler.GetTag)
			tags.GET("", r.TagHandler.SearchTags)

			// Protected mutations
			tags.Use(middleware.AuthMiddleware())
			{
				tags.POST("", r.TagHandler.CreateTag)
				tags.PUT("/:id", r.TagHandler.UpdateTag)
				tags.DELETE("/:id", r.TagHandler.DeleteTag)
			}
		}

		// ===== Stats =====
		api.GET("/stats", r.StatsHandler.GetStats)

		// ===== Media Upload (Protected) =====
		//		api.POST("/media/upload", middleware.AuthMiddleware(), handler.UploadMedia)
		api.GET("/media/signed-url", handler.GetSignedURL)
		api.GET("/media/upload-url", handler.GetUploadURL)

	}

	return router
}
