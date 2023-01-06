package service

import (
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
	listDao dao.ListDaoInterface
}

func ListService(listDao dao.ListDaoInterface) *listService {
	return &listService{listDao}
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
	return srv.listDao.DeleteList(boardList)
}

func (srv *listService) FindListById(id string) (model.BoardList, error) {
	return srv.listDao.FindListById(id)
}

func (srv *listService) GetLists(boardId string) ([]model.BoardList, error) {
	return srv.listDao.GetLists(boardId)
}
