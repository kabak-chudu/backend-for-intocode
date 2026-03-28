package notes

import (
	"github.com/kabak-chudu/backend-for-intocode/students"
	"gorm.io/gorm"
)

type Note struct {
	gorm.Model
	StudentID uint             `json:"student_id"`                          // FK
	Student   students.Student `json:"student" gorm:"foreignKey:StudentID"` //FK
	Author    string           `json:"author"`
	Text      string           `json:"text"`
}
