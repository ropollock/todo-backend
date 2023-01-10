package service

import (
	"fmt"
	"time"
	"todo/dao"
	"todo/model"
)

type BoardServiceInterface interface {
	CreateBoard(board *model.Board) (*model.Board, error)
	DeleteBoard(board *model.Board) error
	UpdateBoard(board *model.Board) (*model.Board, error)
	FindBoardById(id string) (model.Board, error)
	FindBoardByUserId(userId string) (model.Board, error)
	GetBoards(userId string) ([]model.Board, error)
}

type boardService struct {
	boardDao    dao.BoardDaoInterface
	listService ListServiceInterface
}

func BoardService(boardDao dao.BoardDaoInterface, listService ListServiceInterface) *boardService {
	return &boardService{boardDao, listService}
}

func (srv *boardService) CreateBoard(board *model.Board) (*model.Board, error) {
	board.CreatedTS = time.Now()
	return srv.boardDao.CreateBoard(board)
}

func (srv *boardService) UpdateBoard(board *model.Board) (*model.Board, error) {
	board.ModifiedTS = time.Now()
	return srv.boardDao.UpdateBoard(board)
}

func (srv *boardService) DeleteBoard(board *model.Board) error {
	boardErr := srv.boardDao.DeleteBoard(board)
	if boardErr == nil {
		lists, err := srv.listService.GetLists(board.ID.Hex())
		if err == nil {
			for _, list := range lists {
				listErr := srv.listService.DeleteList(&list)
				if listErr != nil {
					fmt.Printf("failed to delete list. %s", listErr)
				}
			}
		}
	}

	return boardErr
}

func (srv *boardService) FindBoardById(id string) (model.Board, error) {
	return srv.boardDao.FindBoardById(id)
}

func (srv *boardService) FindBoardByUserId(userId string) (model.Board, error) {
	return srv.boardDao.FindBoardByUserId(userId)
}

func (srv *boardService) GetBoards(userId string) ([]model.Board, error) {
	return srv.boardDao.GetBoards(userId)
}
