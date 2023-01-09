package service

import (
	"time"
	"todo/dao"
	"todo/model"
)

type TaskServiceInterface interface {
	CreateTask(task *model.Task) (*model.Task, error)
	DeleteTask(task *model.Task) error
	UpdateTask(task *model.Task) (*model.Task, error)
	FindTaskById(id string) (model.Task, error)
	GetTasks(listId string) ([]model.Task, error)
}

type taskService struct {
	taskDao dao.TaskDaoInterface
}

func TaskService(taskDao dao.TaskDaoInterface) *taskService {
	return &taskService{taskDao}
}

func (srv *taskService) CreateTask(task *model.Task) (*model.Task, error) {
	task.CreatedTS = time.Now()
	return srv.taskDao.CreateTask(task)
}

func (srv *taskService) DeleteTask(task *model.Task) error {
	return srv.taskDao.DeleteTask(task)
}

func (srv *taskService) UpdateTask(task *model.Task) (*model.Task, error) {
	task.ModifiedTS = time.Now()
	return srv.taskDao.UpdateTask(task)
}

func (srv *taskService) FindTaskById(id string) (model.Task, error) {
	return srv.taskDao.FindTaskById(id)
}

func (srv *taskService) GetTasks(listId string) ([]model.Task, error) {
	return srv.taskDao.GetTasks(listId)
}
