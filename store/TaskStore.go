package store

import (
	"time"

	"github.com/yasngleer/studentplan/models"
	"github.com/yasngleer/studentplan/utils"

	"gorm.io/gorm"
)

type TaskStore struct {
	DB *gorm.DB
}

func NewTaskStore(db *gorm.DB) *TaskStore {
	db.AutoMigrate(&models.Task{})
	return &TaskStore{DB: db}
}

func (s *TaskStore) CreateTask(task *models.Task) error {
	var count1 int64
	var count2 int64

	s.DB.Where("start_date BETWEEN ? AND ?", task.EndDate, task.StartDate).Model(models.Task{}).Count(&count1)
	s.DB.Where("end_date BETWEEN ? AND ?", task.EndDate, task.StartDate).Model(models.Task{}).Count(&count2)
	if count1 > 0 || count2 > 0 {
		return utils.NewError("There is task with same date")
	}

	return s.DB.Create(task).Error
}

func (s *TaskStore) GetTaskByID(id uint) (*models.Task, error) {
	var task models.Task
	if err := s.DB.First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}
func (s *TaskStore) GetTasksByUser(userID uint, date1 *time.Time, date2 *time.Time) ([]models.Task, error) {
	var tasks []models.Task
	if date1 == nil || date2 == nil {
		if err := s.DB.Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
			return nil, err
		}
	} else {
		if err := s.DB.Where("user_id = ?", userID).
			Where("start_date BETWEEN ? AND ?", date1, date2).
			Find(&tasks).Error; err != nil {
			return nil, err
		}
	}

	return tasks, nil
}

func (s *TaskStore) UpdateTask(task *models.Task) error {
	return s.DB.Save(task).Error
}

func (s *TaskStore) DeleteTask(id uint) error {
	return s.DB.Delete(&models.Task{}, id).Error
}
