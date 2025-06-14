package server

import (
	"github.com/G4L1L10/admin-dashboard-backend/config"
	"github.com/G4L1L10/admin-dashboard-backend/internal/app/handler"
	"github.com/G4L1L10/admin-dashboard-backend/internal/app/repository"
	"github.com/G4L1L10/admin-dashboard-backend/internal/app/service"
	"github.com/G4L1L10/admin-dashboard-backend/internal/db"
	"github.com/G4L1L10/admin-dashboard-backend/internal/router"
)

func Start() error {
	// 1. Load configuration
	config.LoadConfig()

	// 2. Connect to database
	db.ConnectDatabase()

	// 3. Initialize repositories
	courseRepo := repository.NewCourseRepository(db.DB)
	lessonRepo := repository.NewLessonRepository(db.DB)
	questionRepo := repository.NewQuestionRepository(db.DB)
	optionRepo := repository.NewOptionRepository(db.DB)
	tagRepo := repository.NewTagRepository(db.DB)
	questionTagRepo := repository.NewQuestionTagRepository(db.DB)
	statsRepo := repository.NewStatsRepository(db.DB)
	userProgressRepo := repository.NewUserProgressRepository(db.DB)

	// 4. Initialize services
	courseService := service.NewCourseService(courseRepo)
	lessonService := service.NewLessonService(lessonRepo, questionRepo)
	questionService := service.NewQuestionService(questionRepo, optionRepo, tagRepo, questionTagRepo)
	optionService := service.NewOptionService(optionRepo)
	tagService := service.NewTagService(tagRepo)
	statsService := service.NewStatsService(statsRepo)
	userProgressService := service.NewUserProgressService(userProgressRepo)
	gamificationService := service.NewGamificationService(db.DB) // ✅ returns interface by value

	// 5. Initialize handlers
	courseHandler := handler.NewCourseHandler(courseService)
	lessonHandler := handler.NewLessonHandler(lessonService)
	questionHandler := handler.NewQuestionHandler(questionService, optionService, tagService)
	optionHandler := handler.NewOptionHandler(optionService)
	tagHandler := handler.NewTagHandler(tagService)
	statsHandler := handler.NewStatsHandler(statsService)

	// ✅ Create the router
	appRouter := router.NewRouter(
		courseHandler,
		lessonHandler,
		questionHandler,
		optionHandler,
		tagHandler,
		statsHandler,
		userProgressService,
		gamificationService, // ✅ pass interface, not pointer
	)

	server := appRouter.SetupRouter()

	// 6. Start server
	port := config.AppConfig.ServerPort
	if port == "" {
		port = "8080"
	}

	return server.Run(":" + port)
}

