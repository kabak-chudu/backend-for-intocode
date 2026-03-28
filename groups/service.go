package groups

import (
	"errors"

	connectdatabase "github.com/kabak-chudu/backend-for-intocode/connect_database"
)

func GetGroupsFinished(is_finished bool) ([]Group, error) {
	db, err := connectdatabase.Connect()
	if err != nil {
		return nil, err
	}

	var groups []Group
	if err := db.Where("is_finished = ?", is_finished).Find(&groups).Error; err != nil {
		return nil, err
	}
	if len(groups) == 0 {
		return nil, errors.New("завершенных групп не нашлось")
	}
	return groups, nil
}

func GetGroupsByWeek(week uint) ([]Group, error) {
	db, err := connectdatabase.Connect()
	if err != nil {
		return nil, errors.New("не удалость кстановить соединение с БД")
	}
	groups := []Group{}
	if err := db.Where("current_week = ?", week).Find(&groups).Error; err != nil {
		return nil, err
	}
	if len(groups) == 0 {
		return nil, errors.New("групп на этой неделе не найдено")
	}

	return groups, nil
}

func DeleteGroup(id uint) error {
	db, err := connectdatabase.Connect()
	if err != nil {
		return errors.New("не удалость кстановить соединение с БД")
	}
	result := db.Delete(&Group{}, id)
	if result.RowsAffected == 0 {
		return errors.New("не удалось удалить")
	}
	return nil
}

func CreateGroup(title string, total_weeks uint) (*Group, error) {
	db, err := connectdatabase.Connect()
	if err != nil {
		return nil, errors.New("не удалость кстановить соединение с БД")
	}
	group := Group{
		Title:        title,
		Total_weeks:  total_weeks,
		Current_week: 1,
		Is_finished:  false,
	}
	res := db.Create(&group)
	if res.RowsAffected == 0 {
		return nil, errors.New("не удалось создать группу")
	}

	return &group, nil
}

func UpdateGroup(id, current_week uint) (*Group, error) {
	db, err := connectdatabase.Connect()
	if err != nil {
		return nil, errors.New("не удалость кстановить соединение с БД")
	}
	group, err := GetGroupID(id)
	if err != nil {
		return nil, err
	}

	if group.Is_finished {
		return nil, errors.New("группа уже завершена изменения запрещены")
	}

	if current_week >= group.Total_weeks {
		group.Current_week = group.Total_weeks
		group.Is_finished = true
	} else {
		group.Current_week = current_week
	}
	if err := db.Save(&group).Error; err != nil {
		return nil, err
	}

	return group, nil
}

func GetAllGroups() ([]Group, error) {
	groups := []Group{}

	db, err := connectdatabase.Connect()
	if err != nil {
		return nil, errors.New("не удалость кстановить соединение с БД")
	}
	if err := db.Find(&groups).Error; err != nil {
		return nil, err
	}
	if len(groups) == 0 {
		return nil, errors.New("не нашлось групп в базе данных")
	}

	return groups, nil
}

func GetGroupID(id uint) (*Group, error) {
	db, err := connectdatabase.Connect()
	if err != nil {
		return nil, errors.New("не удалость кстановить соединение с БД")
	}
	group := Group{}
	if err := db.First(&group, id).Error; err != nil {
		return nil, err
	}

	return &group, nil
}
