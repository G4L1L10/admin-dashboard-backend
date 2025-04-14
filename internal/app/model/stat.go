package model

type Stats struct {
	Courses   int `json:"courses"`
	Lessons   int `json:"lessons"`
	Questions int `json:"questions"`
	Tags      int `json:"tags"`
}
