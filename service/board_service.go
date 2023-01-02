package service

import (
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
	boardDao dao.BoardDaoInterface
}

func BoardService(boardDao dao.BoardDaoInterface) *boardService {
	return &boardService{boardDao}
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
	return srv.boardDao.DeleteBoard(board)
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
