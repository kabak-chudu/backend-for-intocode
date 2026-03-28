package students

import (
	"github.com/kabak-chudu/backend-for-intocode/groups"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Full_name      string        `json:"full_name"`
	Email          string        `json:"email"`
	Telegram       string        `json:"telegram"`
	GroupID        uint          `json:"group_id"`                        // FK
	Group          *groups.Group `json:"group" gorm:"foreignKey:GroupID"` // FK
	Tuition_total  uint          `json:"tuition_total"`                   // полная стоимость обучения
	Tuition_paid   uint          `json:"tuition_paid"`                    // сколько оплачено
	Payment_status string        `json:"payment_status"`
	// Paid - оплата
	// Unpaid - нет оплаты
	// Partial - неполная оплата
	Study_status string `json:"study_status"`
	// learning
	// job_search
	// offer
	// working

}
