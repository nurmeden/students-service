package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis"
	"github.com/nurmeden/students-service/internal/app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type StudentRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
	cache      *redis.Client
}

func NewStudentRepository(client *mongo.Client, dbName string, collectionName string, cache *redis.Client) (*StudentRepository, error) {
	r := &StudentRepository{
		client: client,
		cache:  cache,
	}

	collection := client.Database(dbName).Collection(collectionName)
	r.collection = collection

	return r, nil
}

func (r *StudentRepository) Create(ctx context.Context, student *model.Student) (*model.Student, error) {
	err := r.client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}
	_, err = r.collection.InsertOne(ctx, student)
	if err != nil {
		fmt.Println("osinda")
		return nil, fmt.Errorf("failed to create student: %v", err)
	}
	fmt.Println("bari jaksi")
	return student, nil
}

func (r *StudentRepository) Read(ctx context.Context, id string) (*model.Student, error) {
	cachedResult, err := r.cache.Get(id).Result()
	if err == nil {
		student := &model.Student{}
		err = json.Unmarshal([]byte(cachedResult), student)
		if err != nil {
			return nil, err
		}
		return student, nil
	}
	fmt.Printf("cachedResult: %v\n", cachedResult)

	var student model.Student
	studentId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}
	filter := bson.M{"_id": studentId}
	err = r.collection.FindOne(ctx, filter).Decode(&student)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Если студент не найден, возвращаем nil и ошибку nil
		}
		return nil, fmt.Errorf("failed to read student: %v", err)
	}

	studentJSON, err := json.Marshal(student)
	if err != nil {
		return nil, err
	}

	err = r.cache.Set(id, studentJSON, 0).Err()
	if err != nil {
		return nil, err
	}

	return &student, nil
}

func (r *StudentRepository) GetStudentByCoursesID(ctx context.Context, id string) (*model.Student, error) {
	cachedResult, err := r.cache.Get(id).Result()
	if err == nil {
		student := &model.Student{}
		err = json.Unmarshal([]byte(cachedResult), student)
		if err != nil {
			return nil, err
		}
		return student, nil
	}

	var student model.Student

	filter := bson.M{"courses": id}
	err = r.collection.FindOne(ctx, filter).Decode(&student)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Если студент не найден, возвращаем nil и ошибку nil
		}
		return nil, fmt.Errorf("failed to read student: %v", err)
	}

	studentJSON, err := json.Marshal(student)
	if err != nil {
		return nil, err
	}

	err = r.cache.Set(id, studentJSON, 0).Err()
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *StudentRepository) Update(ctx context.Context, student *model.Student) (*model.Student, error) {
	filter := bson.M{"_id": student.ID}
	update := bson.M{"$set": bson.M{
		"firstName": student.FirstName,
		"lastName":  student.LastName,
		"age":       student.Age,
	}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update student: %v", err)
	}
	return student, nil
}

// Delete - удаление студента по ID
func (r *StudentRepository) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete student: %v", err)
	}
	return nil
}
