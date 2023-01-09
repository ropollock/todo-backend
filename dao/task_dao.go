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

type taskDao struct {
	databaseProvider data.MongoDBProviderInterface
}

type TaskDaoInterface interface {
	CreateTask(task *model.Task) (*model.Task, error)
	DeleteTask(task *model.Task) error
	UpdateTask(task *model.Task) (*model.Task, error)
	FindTaskById(id string) (model.Task, error)
	GetTasks(listId string) ([]model.Task, error)
}

func TaskDao(databaseProvider data.MongoDBProviderInterface) *taskDao {
	return &taskDao{databaseProvider}
}

func (dao *taskDao) CreateTask(task *model.Task) (*model.Task, error) {
	insertResult, err := dao.databaseProvider.GetTasksCollection().InsertOne(dao.databaseProvider.GetContext(), task)
	result, _ := dao.FindTaskById(insertResult.InsertedID.(primitive.ObjectID).Hex())
	return &result, err
}

func (dao *taskDao) DeleteTask(task *model.Task) error {
	_, err := dao.databaseProvider.GetTasksCollection().DeleteOne(dao.databaseProvider.GetContext(), bson.M{"_id": task.ID})
	return err
}

func (dao *taskDao) UpdateTask(task *model.Task) (*model.Task, error) {
	_, err := dao.databaseProvider.GetTasksCollection().ReplaceOne(context.Background(), bson.M{"_id": task.ID}, task)
	result, _ := dao.FindTaskById(task.ID.Hex())
	return &result, err
}

func (dao *taskDao) FindTaskById(id string) (model.Task, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}

	result := dao.databaseProvider.GetTasksCollection().FindOne(context.Background(), bson.M{"_id": objectId})
	resultTask := model.Task{}
	err = result.Decode(&resultTask)
	if err != nil {
		fmt.Println(err)
		return resultTask, fmt.Errorf("an error occurred while decoding record : %v", err)
	}
	return resultTask, nil
}

func (dao *taskDao) GetTasks(listId string) ([]model.Task, error) {
	listObjectId, err := primitive.ObjectIDFromHex(listId)
	if err != nil {
		log.Println("Invalid list id")
	}

	var results []model.Task
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	cursor, err := dao.databaseProvider.GetTasksCollection().Find(context.Background(), bson.M{"list_id": listObjectId})
	if err != nil {
		fmt.Println("Finding all tasks ERROR:", err)
		defer cursor.Close(ctx)
	}

	err = cursor.All(ctx, &results)
	if err != nil {
		defer cursor.Close(ctx)
		return results, err
	}

	return results, nil
}
