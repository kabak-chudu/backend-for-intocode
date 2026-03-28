package groups

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Title        string `json:"title"`
	Current_week uint    `json:"current_week"`
	Total_weeks  uint    `json:"total_weeks"`
	Is_finished  bool   `json:"is_finished"`
}
