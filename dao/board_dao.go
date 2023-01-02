package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
	"todo/data"
	"todo/model"
)

type boardDao struct {
	databaseProvider data.MongoDBProviderInterface
}

type BoardDaoInterface interface {
	CreateBoard(user *model.Board) (*model.Board, error)
	DeleteBoard(user *model.Board) error
	UpdateBoard(board *model.Board) (*model.Board, error)
	FindBoardById(id string) (model.Board, error)
	FindBoardByUserId(username string) (model.Board, error)
	GetBoards(userId string) ([]model.Board, error)
}

func BoardDao(databaseProvider data.MongoDBProviderInterface) *boardDao {
	return &boardDao{databaseProvider}
}

func (dao *boardDao) CreateBoard(board *model.Board) (*model.Board, error) {
	insertResult, err := dao.databaseProvider.GetBoardsCollection().InsertOne(dao.databaseProvider.GetContext(), board)
	result, _ := dao.FindBoardById(insertResult.InsertedID.(primitive.ObjectID).Hex())
	return &result, err
}

func (dao *boardDao) DeleteBoard(board *model.Board) error {
	_, err := dao.databaseProvider.GetBoardsCollection().DeleteOne(dao.databaseProvider.GetContext(), bson.M{"_id": board.ID})
	return err
}

func (dao *boardDao) UpdateBoard(board *model.Board) (*model.Board, error) {
	_, err := dao.databaseProvider.GetBoardsCollection().ReplaceOne(context.Background(), bson.M{"_id": board.ID}, board)
	result, _ := dao.FindBoardById(board.ID.Hex())
	return &result, err
}

func (dao *boardDao) FindBoardById(id string) (model.Board, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}

	result := dao.databaseProvider.GetBoardsCollection().FindOne(context.Background(), bson.M{"_id": objectId})
	resultBoard := model.Board{}
	err = result.Decode(&resultBoard)
	if err != nil {
		fmt.Println(err)
		return resultBoard, fmt.Errorf("an error occurred while decoding record : %v", err)
	}
	return resultBoard, nil
}

func (dao *boardDao) FindBoardByUserId(userId string) (model.Board, error) {
	result := dao.databaseProvider.GetBoardsCollection().FindOne(context.Background(), bson.M{"owner_id": userId})
	resultBoard := model.Board{}
	err := result.Decode(&resultBoard)
	if err != nil {
		fmt.Println(err)
		return resultBoard, fmt.Errorf("an error occurred while decoding record : %v", err)
	}
	return resultBoard, nil
}

func (dao *boardDao) GetBoards(userId string) ([]model.Board, error) {
	fmt.Println("Finding boards owned by " + userId)
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Println("Invalid id")
	}

	var results []model.Board
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	cursor, err := dao.databaseProvider.GetBoardsCollection().Find(context.Background(), bson.M{"owner_id": objectId})
	if err != nil {
		fmt.Println("Finding all boards ERROR:", err)
		defer cursor.Close(ctx)
	}

	err = cursor.All(ctx, &results)
	if err != nil {
		defer cursor.Close(ctx)
		return results, err
	}

	return results, nil
}
