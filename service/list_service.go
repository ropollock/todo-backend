package service

import (
	"fmt"
	"time"
	"todo/dao"
	"todo/model"
)

type ListServiceInterface interface {
	CreateList(boardList *model.BoardList) (*model.BoardList, error)
	DeleteList(boardList *model.BoardList) error
	UpdateList(boardList *model.BoardList) (*model.BoardList, error)
	FindListById(id string) (model.BoardList, error)
	GetLists(boardId string) ([]model.BoardList, error)
}

type listService struct {
	listDao     dao.ListDaoInterface
	taskService TaskServiceInterface
}

func ListService(listDao dao.ListDaoInterface, taskService TaskServiceInterface) *listService {
	return &listService{listDao, taskService}
}

func (srv *listService) CreateList(boardList *model.BoardList) (*model.BoardList, error) {
	boardList.CreatedTS = time.Now()
	return srv.listDao.CreateList(boardList)
}

func (srv *listService) UpdateList(boardList *model.BoardList) (*model.BoardList, error) {
	boardList.ModifiedTS = time.Now()
	return srv.listDao.UpdateList(boardList)
}

func (srv *listService) DeleteList(boardList *model.BoardList) error {
	listErr := srv.listDao.DeleteList(boardList)
	if listErr == nil {
		tasks, err := srv.taskService.GetTasks(boardList.ID.Hex())
		if err == nil {
			for _, task := range tasks {
				taskErr := srv.taskService.DeleteTask(&task)
				if taskErr != nil {
					fmt.Printf("failed to delete task. %s", listErr)
				}
			}
		}
	}

	return listErr
}

func (srv *listService) FindListById(id string) (model.BoardList, error) {
	return srv.listDao.FindListById(id)
}

func (srv *listService) GetLists(boardId string) ([]model.BoardList, error) {
	return srv.listDao.GetLists(boardId)
}
