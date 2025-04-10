package model

type FullLesson struct {
	Lesson    *Lesson                `json:"lesson"`
	Questions []*QuestionWithOptions `json:"questions"`
}
