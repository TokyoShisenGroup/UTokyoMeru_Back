package db

import (
	"backend/internal/model"
)


type MessageCRUD struct{}

func (crud MessageCRUD) CreateByObject(m *model.Message) error {
	db, err := GetDatabaseInstance()
	if err != nil {
		return err
	}
	return db.Create(m).Error
}

func (crud MessageCRUD) FindAll() ([]model.Message, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}
	var messages []model.Message
	result := db.Preload("Sender").Preload("Receiver").Find(&messages)
	return messages, result.Error
}

func (crud MessageCRUD) FindAllOrdered(orderBy string, order string) ([]model.Message, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}
	var messages []model.Message
	result := db.Preload("Sender").Preload("Receiver").Order(orderBy + " " + order).Find(&messages)
	return messages, result.Error
}

func (crud MessageCRUD) FindById(id uint) (*model.Message, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}
	var message model.Message
	result := db.Preload("Sender").Preload("Receiver").First(&message, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &message, nil
}

func (crud MessageCRUD) UpdateByObject(m *model.Message) error {
	db, err := GetDatabaseInstance()
	if err != nil {
		return err
	}
	return db.Save(m).Error
}

func (crud MessageCRUD) UpdateByField(field string, value interface{}, m *model.Message) error {
	db, err := GetDatabaseInstance()
	if err != nil {
		return err
	}
	return db.Model(m).Update(field, value).Error
}

func (crud MessageCRUD) DeleteById(id uint) error {
	db, err := GetDatabaseInstance()
	if err != nil {
		return err
	}
	return db.Delete(&model.Message{}, id).Error
}

func (crud MessageCRUD) FindAllByField(fieldName string, value interface{}, orderBy string, order string) ([]model.Message, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}
	var messages []model.Message
	result := db.Preload("Sender").Preload("Receiver").Where(fieldName+" = ?", value).Order(orderBy + " " + order).Find(&messages)
	return messages, result.Error
}

func (crud MessageCRUD) FindOneByUniqueField(fieldName string, value interface{}) (*model.Message, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}
	var message model.Message
	result := db.Preload("Sender").Preload("Receiver").Where(fieldName+" = ?", value).First(&message)
	return &message, result.Error
}

// 待实现
func (crud MessageCRUD) Search(ops ...searchOption) ([]model.Message, error) {
	return nil, nil
}




//这部分是消息队列的记录
type MessageQueueCRUD struct{}

func (crud MessageQueueCRUD) CreateByObject(m *model.MessageQueue) error {
	db, err := GetDatabaseInstance()
	if err != nil {
		return err
	}
	return db.Create(m).Error
}

func (crud MessageQueueCRUD) FindAll() ([]model.MessageQueue, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}
	var messages []model.MessageQueue
	result := db.Preload("Message").Find(&messages)
	return messages, result.Error
}

func (crud MessageQueueCRUD) FindByID(id uint) (*model.MessageQueue, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}
	var message model.MessageQueue
	result := db.Preload("Message").First(&message, id)
	return &message, result.Error
}

func (crud MessageQueueCRUD) FindAllByField(fieldName string, value interface{}, orderBy string, order string) ([]model.MessageQueue, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}
	var messages []model.MessageQueue
	result := db.Preload("Message").Where(fieldName+" = ?", value).Order(orderBy + " " + order).Find(&messages)
	return messages, result.Error
}

func (crud MessageQueueCRUD) FindOneByUniqueField(fieldName string, value interface{}) (*model.MessageQueue, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}
	var message model.MessageQueue
	result := db.Preload("Message").Where(fieldName+" = ?", value).First(&message)
	return &message, result.Error
}

func (crud MessageQueueCRUD) UpdateByField(field string, value interface{}, m *model.MessageQueue) error {
	db, err := GetDatabaseInstance()
	if err != nil {
		return err
	}
	return db.Model(m).Update(field, value).Error
}

func (crud MessageQueueCRUD) UpdateByObject(m *model.MessageQueue) error {
	db, err := GetDatabaseInstance()
	if err != nil {
		return err
	}
	return db.Save(m).Error
}

func (crud MessageQueueCRUD) DeleteById(id uint) error {
	db, err := GetDatabaseInstance()
	if err != nil {
		return err
	}
	return db.Delete(&model.MessageQueue{}, id).Error
}

// 待实现
func (crud MessageQueueCRUD) Search(ops ...searchOption) ([]model.MessageQueue, error) {
	return nil, nil
}
