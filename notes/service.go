package notes

import (
	"errors"

	connectdatabase "github.com/kabak-chudu/backend-for-intocode/connect_database"
	"github.com/kabak-chudu/backend-for-intocode/students"
)

func GetAllNotes() ([]Note, error) {
	notes := []Note{}

	db, err := connectdatabase.Connect()
	if err != nil {
		return nil, errors.New("не удалость кстановить соединение с БД")
	}
	if err := db.Find(&notes).Error; err != nil {
		return nil, err
	}
	if len(notes) == 0 {
		return nil, errors.New("не нашлось заметок в базе данных")
	}

	return notes, nil
}

func GetNoteID(id uint) (*Note, error) {
	db, err := connectdatabase.Connect()
	if err != nil {
		return nil, errors.New("не удалость кстановить соединение с БД")
	}
	note := Note{}
	if err := db.First(&note, id).Error; err != nil {
		return nil, err
	}

	return &note, nil
}

func GetAllNotesForStudent(student_id uint) ([]Note, error) {
	db, err := connectdatabase.Connect()
	if err != nil {
		return nil, errors.New("не удалость кстановить соединение с БД")
	}
	notes := []Note{}
	if err := db.Preload("Student").Where("student_id = ?", student_id).Find(&notes).Error; err != nil {
		return nil, err
	}

	if len(notes) == 0 {
		return nil, errors.New("не найдено заметок по студенту")
	}

	return notes, nil
}

func CreateNote(student_id uint, author, text string) (*Note, error) {
	db, err := connectdatabase.Connect()
	if err != nil {
		return nil, errors.New("не удалость кстановить соединение с БД")
	}

	student, err := students.GetStudentByID(student_id)
	if err != nil {
		return nil, errors.New("студента по такому вйди нет")
	}

	note := Note{
		StudentID: student_id,
		Author:    author,
		Text:      text,
		Student:   *student,
	}
	result := db.Create(&note)
	if result.RowsAffected == 0 {
		return nil, errors.New("не удалось создать заметку")
	}

	return &note, nil

}

func UpdateNote(id uint, author string) (*Note, error) {
	db, err := connectdatabase.Connect()
	if err != nil {
		return nil, errors.New("не удалость кстановить соединение с БД")
	}
	note, err := GetNoteID(id)
	if err != nil {
		return nil, err
	}
	if author != "" {
		note.Author = author
	}
	if err := db.Preload("Student").Save(&note).Error; err != nil {
		return nil, err
	}

	return note, nil
}

func DeleteNote(id uint) error {
	db, err := connectdatabase.Connect()
	if err != nil {
		return errors.New("не удалость кстановить соединение с БД")
	}
	note, err := GetNoteID(id)
	if err != nil {
		return err
	}
	res := db.Delete(&note)
	if res.RowsAffected == 0 {
		return errors.New("не удалось удалить(")
	}

	return nil
}
