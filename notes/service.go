package notes

import (
	"errors"

	connectdatabase "github.com/kabak-chudu/backend-for-intocode/connect_database"
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
