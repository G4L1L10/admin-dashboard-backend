package service

import (
	"github.com/G4L1L10/admin-dashboard-backend/internal/app/model"
	"github.com/G4L1L10/admin-dashboard-backend/internal/app/repository"
)

type CourseService struct {
	courseRepo *repository.CourseRepository
}

func NewCourseService(courseRepo *repository.CourseRepository) *CourseService {
	return &CourseService{courseRepo: courseRepo}
}

// CREATE
func (s *CourseService) CreateCourse(course *model.Course) error {
	return s.courseRepo.Create(course)
}

// READ - Single Course
func (s *CourseService) GetCourseByID(id string) (*model.Course, error) {
	return s.courseRepo.GetByID(id)
}

// READ - List All Courses ✅
func (s *CourseService) ListCourses() ([]*model.Course, error) {
	return s.courseRepo.ListCourses()
}

// UPDATE
func (s *CourseService) UpdateCourse(course *model.Course) error {
	return s.courseRepo.Update(course)
}

// DELETE
func (s *CourseService) DeleteCourse(id string) error {
	return s.courseRepo.Delete(id)
}

