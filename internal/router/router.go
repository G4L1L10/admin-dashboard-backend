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

	// Apply global middlewares
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.CORSMiddleware())

	// Public API
	api := router.Group("/api")
	{
		// ===== Courses =====
		courses := api.Group("/courses")
		{
			courses.POST("", r.CourseHandler.CreateCourse)
			courses.GET("", r.CourseHandler.ListCourses) // 🛠️ ADDED: List all courses
			courses.GET("/:id", r.CourseHandler.GetCourse)
			courses.PUT("/:id", r.CourseHandler.UpdateCourse)
			courses.DELETE("/:id", r.CourseHandler.DeleteCourse)
		}

		// ===== Lessons =====
		lesson := api.Group("/lessons")
		{
			lesson.POST("", r.LessonHandler.CreateLesson)
			lesson.GET("/detail/:id", r.LessonHandler.GetLesson)
			lesson.GET("/full/:id", r.LessonHandler.GetFullLesson)
			lesson.GET("/by-course/:course_id", r.LessonHandler.GetLessonsByCourseID) // ✅ NEW
			lesson.GET("/:lesson_id/questions", r.QuestionHandler.GetQuestionsByLesson)
			lesson.PUT("/:id", r.LessonHandler.UpdateLesson)
			lesson.DELETE("/:id", r.LessonHandler.DeleteLesson)
		}

		// ===== Questions =====
		questions := api.Group("/questions")
		{
			questions.POST("", r.QuestionHandler.CreateQuestion)
			questions.GET("/:id", r.QuestionHandler.GetQuestion)
			questions.GET("", r.QuestionHandler.GetQuestionsByTag)
			questions.PUT("/:id", r.QuestionHandler.UpdateQuestion)
			questions.DELETE("/:id", r.QuestionHandler.DeleteQuestion)
		}

		// ===== Options =====
		options := api.Group("/options")
		{
			options.POST("", r.OptionHandler.CreateOption)
			options.GET("/:id", r.OptionHandler.GetOption)
			options.PUT("/:id", r.OptionHandler.UpdateOption)
			options.DELETE("/:id", r.OptionHandler.DeleteOption)
		}

		// ===== Tags =====
		tags := api.Group("/tags")
		{
			tags.POST("", r.TagHandler.CreateTag)
			tags.GET("/:id", r.TagHandler.GetTag)
			tags.GET("", r.TagHandler.SearchTags)
			tags.PUT("/:id", r.TagHandler.UpdateTag)
			tags.DELETE("/:id", r.TagHandler.DeleteTag)
		}

		// ===== Stats =====
		api.GET("/stats", r.StatsHandler.GetStats)
	}

	return router
}

