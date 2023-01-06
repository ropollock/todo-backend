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

type listDao struct {
	databaseProvider data.MongoDBProviderInterface
}

type ListDaoInterface interface {
	CreateList(board *model.BoardList) (*model.BoardList, error)
	DeleteList(board *model.BoardList) error
	UpdateList(board *model.BoardList) (*model.BoardList, error)
	FindListById(id string) (model.BoardList, error)
	GetLists(boardId string) ([]model.BoardList, error)
}

func ListDao(databaseProvider data.MongoDBProviderInterface) *listDao {
	return &listDao{databaseProvider}
}

func (dao *listDao) CreateList(boardList *model.BoardList) (*model.BoardList, error) {
	insertResult, err := dao.databaseProvider.GetListsCollection().InsertOne(dao.databaseProvider.GetContext(), boardList)
	result, _ := dao.FindListById(insertResult.InsertedID.(primitive.ObjectID).Hex())
	return &result, err
}

func (dao *listDao) DeleteList(boardList *model.BoardList) error {
	_, err := dao.databaseProvider.GetListsCollection().DeleteOne(dao.databaseProvider.GetContext(), bson.M{"_id": boardList.ID})
	return err
}

func (dao *listDao) UpdateList(boardList *model.BoardList) (*model.BoardList, error) {
	_, err := dao.databaseProvider.GetListsCollection().ReplaceOne(context.Background(), bson.M{"_id": boardList.ID}, boardList)
	result, _ := dao.FindListById(boardList.ID.Hex())
	return &result, err
}

func (dao *listDao) FindListById(id string) (model.BoardList, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}

	result := dao.databaseProvider.GetListsCollection().FindOne(context.Background(), bson.M{"_id": objectId})
	resultList := model.BoardList{}
	err = result.Decode(&resultList)
	if err != nil {
		fmt.Println(err)
		return resultList, fmt.Errorf("an error occurred while decoding record : %v", err)
	}
	return resultList, nil
}

func (dao *listDao) GetLists(boardId string) ([]model.BoardList, error) {
	boardObjectId, err := primitive.ObjectIDFromHex(boardId)
	if err != nil {
		log.Println("Invalid board id")
	}

	var results []model.BoardList
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	cursor, err := dao.databaseProvider.GetListsCollection().Find(context.Background(), bson.M{"board_id": boardObjectId})
	if err != nil {
		fmt.Println("Finding all lists ERROR:", err)
		defer cursor.Close(ctx)
	}

	err = cursor.All(ctx, &results)
	if err != nil {
		defer cursor.Close(ctx)
		return results, err
	}

	return results, nil
}
