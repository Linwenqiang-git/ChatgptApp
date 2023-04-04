package repo

import (
	. "github.com/lwq/internal/data"
	. "github.com/lwq/internal/data/dto"

	. "github.com/lwq/internal/data/entity"
)

type MessageRepo struct {
	dbContext DbContext
}

func ProvideMessageRepo(dbContext DbContext) MessageRepo {
	return MessageRepo{
		dbContext: dbContext,
	}
}

func (m MessageRepo) CreateUser(user User) (int, error) {
	tx := m.dbContext.GetDb().Create(&user)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return int(tx.RowsAffected), nil
}

func (m MessageRepo) GetUser(userName string) (User, error) {
	var user User
	tx := m.dbContext.GetDb().Table(User.TableName(user)).Where("Name = ?", userName).First(&user)
	if tx.Error != nil {
		return User{}, tx.Error
	}
	return user, nil
}

func (m MessageRepo) InserUserMessage(records []ChatRecord) (int, error) {
	if records == nil {
		return 0, nil
	}
	tx := m.dbContext.GetDb().Create(&records)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return int(tx.RowsAffected), nil
}
func (m MessageRepo) GetUserHistory(userName string) ([]UserHistoryDto, error) {
	dto := []UserHistoryDto{}
	var c ChatRecord
	main_table := ChatRecord.TableName(c)
	tx := m.dbContext.GetDb().Table(main_table).Select("role, Message").Joins("left join user on "+main_table+".user_name = user.name").Where("name = ?", userName).Order(main_table + ".created_at asc").Scan(&dto)
	return dto, tx.Error
}
