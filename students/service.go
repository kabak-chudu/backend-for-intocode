package students

import (
	"errors"
	"fmt"

	connectdatabase "github.com/kabak-chudu/backend-for-intocode/connect_database"
	"github.com/kabak-chudu/backend-for-intocode/groups"
)

func DeleteStudent(id uint) error {
	db, err := connectdatabase.Connect()
	if err != nil {
		return errors.New("не удалость кстановить соединение с БД")
	}
	result := db.Delete(&Student{}, id)
	if result.RowsAffected == 0 {
		return errors.New("не удалось удалить")
	}
	return nil
}

func CreateStudent(full_name, email, telegram string, group_id uint, tuition_total uint) (*Student, error) {
	db, err := connectdatabase.Connect()
	if err != nil {
		return nil, errors.New("не удалость кстановить соединение с БД")
	}

	student := Student{
		Full_name:      full_name,
		Email:          email,
		Telegram:       telegram,
		GroupID:        group_id,
		Tuition_total:  tuition_total,
		Tuition_paid:   0,          // начальная оплата 0
		Payment_status: "unpaid",   // начальный статус - не оплачено
		Study_status:   "learning", // начальный статус - обучается
	}
	// валидация!!!
	if full_name == "" {
		return nil, errors.New("введите имя студента")
	}

	group, err := groups.GetGroupID(group_id)
	if err != nil {
		return nil, errors.New("такой группы нет")
	} else {
		student.GroupID = group_id
		student.Group = group
	}

	res := db.Create(&student)
	if res.Error != nil {
		return nil, fmt.Errorf("не удалось создать студента: %w", res.Error)
	}

	return &student, nil
}

func UpdateStudent(id uint, req_tuition *uint, req_email, req_telegram, req_status, req_full_name *string) (*Student, error) {
	db, err := connectdatabase.Connect()
	if err != nil {
		return nil, errors.New("ошибка БД")
	}

	student, err := GetStudentByID(id)
	if err != nil {
		return nil, err
	}
	if err := db.Preload("Group").First(&student, id).Error; err != nil {
		return nil, errors.New("студент не найден")
	}
	if req_tuition != nil && *req_tuition > 0 {
		student.Tuition_paid += *req_tuition

		if student.Tuition_paid >= student.Tuition_total {
			student.Tuition_paid = student.Tuition_total
			student.Payment_status = "paid"
		} else {
			student.Payment_status = "partial"
		}
	}

	if req_email != nil && *req_email != "" {
		student.Email = *req_email
	}
	if req_telegram != nil && *req_telegram != "" {
		student.Telegram = *req_telegram
	}
	if req_full_name != nil && *req_full_name != "" {
		student.Full_name = *req_full_name
	}

	if req_status != nil {
		validStatuses := map[string]bool{"learning": true, "job_search": true, "working": true, "offer": true}
		if !validStatuses[*req_status] {
			return nil, errors.New("такого статуса нет")
		}
		student.Study_status = *req_status
	}

	if err := db.Save(&student).Error; err != nil {
		return nil, err
	}
	return student, nil
}

func GetStudentByID(id uint) (*Student, error) {
	db, err := connectdatabase.Connect()
	if err != nil {
		return nil, errors.New("не удалость кстановить соединение с БД")
	}
	student := Student{}
	if err := db.Preload("Group").First(&student, id).Error; err != nil {
		return nil, err
	}

	return &student, err
}

func GetAllStudents() ([]Student, error) {
	students := []Student{}

	db, err := connectdatabase.Connect()
	if err != nil {
		return nil, errors.New("не удалость кстановить соединение с БД")
	}
	if err := db.Preload("Group").Find(&students).Error; err != nil {
		return nil, err
	}
	if len(students) == 0 {
		return []Student{}, nil
	}

	return students, nil
}

func GetStudentsFiltered(groupID string, paymentStatus string, studyStatus string) ([]Student, error) {
	db, err := connectdatabase.Connect()
	if err != nil {
		return nil, err
	}
	students := []Student{}

	if groupID != "" {
		db = db.Where("group_id = ?", groupID)
	}
	if paymentStatus != "" {
		db = db.Where("payment_status = ?", paymentStatus)
	}
	if studyStatus != "" {
		db = db.Where("study_status = ?", studyStatus)
	}

	if err := db.Find(&students).Error; err != nil {
		return nil, err
	}

	return students, nil
}

func GetStudentsByGroupID(group_id uint) ([]Student, error) {
	db, err := connectdatabase.Connect()
	if err != nil {
		return nil, errors.New("не удалость кстановить соединение с БД")
	}
	students := []Student{}
	if err := db.Preload("Group").Where("group_id = ?", group_id).Find(&students).Error; err != nil {
		return nil, err
	}
	if len(students) == 0 {
		return []Student{}, nil
	}

	return students, nil
}

// func GetStudentsByPaymentStatus(payment_status string) ([]Student, error) {
// 	db, err := connectdatabase.Connect()
// 	if err != nil {
// 		return nil, errors.New("не удалость кстановить соединение с БД")
// 	}
// 	students := []Student{}
// 	if err := db.Preload("Group").Where("payment_status = ?", payment_status).Find(&students).Error; err != nil {
// 		return nil, err
// 	}
// 	if len(students) == 0 {
// 		return []Student{}, nil
// 	}

// 	return students, nil
// }

// func GetStudentsByStudyStatus(study_status string) ([]Student, error) {
// 	db, err := connectdatabase.Connect()
// 	if err != nil {
// 		return nil, errors.New("не удалость кстановить соединение с БД")
// 	}
// 	students := []Student{}
// 	if err := db.Preload("Group").Where("study_status = ?", study_status).Find(&students).Error; err != nil {
// 		return nil, err
// 	}
// 	if len(students) == 0 {
// 		return []Student{}, nil
// 	}

// 	return students, nil
// }
