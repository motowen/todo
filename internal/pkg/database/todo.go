package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/logger"
	model "viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/model/db"
)

func InsertTodo(todo model.Todo) (err error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	_, err = todoCollection.InsertOne(ctx, todo)
	if err != nil {
		logger.Error.Printf("[InsertTodo] Failed: %v", err)
		return fmt.Errorf("[InsertTodo] %s", err.Error())
	}

	return
}

func GetAllTodo() (todos []model.Todo, err error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	cursor, err := todoCollection.Find(ctx, bson.M{})
	if err != nil {
		logger.Error.Printf("[GetAllTodo] Find Failed: %v", err)
		return nil, fmt.Errorf("[GetAllTodo] %s", err.Error())
	}

	err = cursor.All(ctx, &todos)
	if err != nil {
		logger.Error.Printf("[GetAllTodo] All Failed: %v", err)
		return nil, fmt.Errorf("[GetAllTodo] %s", err.Error())
	}

	return
}

func GetTodo(id string) (todo model.Todo, err error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	err = todoCollection.FindOne(ctx, bson.M{"id": id}).Decode(&todo)
	if err != nil {
		logger.Error.Printf("[GetTodo] FindOne Failed: %v", err)
		return model.Todo{}, fmt.Errorf("[GetTodo] %s", err.Error())
	}

	return
}

func UpdateTodo(id string, todo model.Todo) (err error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	_, err = todoCollection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": todo})
	if err != nil {
		logger.Error.Printf("[UpdateTodo] UpdateOne Failed: %v", err)
		return fmt.Errorf("[UpdateTodo] %s", err.Error())
	}

	return
}

func DeleteTodo(id string) (err error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	_, err = todoCollection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		logger.Error.Printf("[DeleteTodo] DeleteOne Failed: %v", err)
		return fmt.Errorf("[DeleteTodo] %s", err.Error())
	}

	return
}
