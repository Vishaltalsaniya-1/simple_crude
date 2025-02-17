package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"simple_crude/db"
	"simple_crude/models"
	"simple_crude/producer"
	"simple_crude/request"
	"simple_crude/response"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserMgr struct {
}

func NewUserManager() *UserMgr {
	return &UserMgr{}
}
func (um *UserMgr) CreateUser(req request.StdRequest) (response.StdResponse, error) {

	mongoCollection := db.GetDB().Database("Courses").Collection("students")
	var existingUser models.Std
	err := mongoCollection.FindOne(context.Background(), bson.M{"name": req.Name}).Decode(&existingUser)
	if err == nil {
		return response.StdResponse{}, fmt.Errorf("user with name '%s' already exists", req.Name)
	} else if err != mongo.ErrNoDocuments {
		log.Printf("MongoDB error checking user existence: %v\n", err)
		return response.StdResponse{}, fmt.Errorf("failed to check user existence: %v", err)
	}

	if req.CreatedAt == nil {
		now := time.Now()
		req.CreatedAt = &now
	}

	if req.UpdatedAt == nil {
		now := time.Now()
		req.UpdatedAt = &now
	}

	user := models.Std{
		Id:          primitive.NewObjectID(),
		Name:        req.Name,
		AuthID:      req.AuthID,
		Description: req.Description,
		Tag:         req.Tag,
		Student:     req.Student,
		CreatedAt:   req.CreatedAt,
		UpdatedAt:   req.UpdatedAt,
		DeletedAt:   nil,
	}

	_, err = mongoCollection.InsertOne(context.Background(), user)
	if err != nil {
		log.Printf("MongoDB insertion error: %v\n", err)
		return response.StdResponse{}, fmt.Errorf("failed to insert new user: %v", err)
	}
	// userData, err := json.Marshal(user)
	// if err != nil {
	// 	logrus.Errorf("Error marshalling user data: %v", err)
	// 	return response.StdResponse{}, err
	// }

	// if err := um.Producer.Publish(userData, "UserCreatedTask"); err != nil {
	// 	logrus.Errorf("Error publishing user creation event: %v", err)
	// }

	jsonData, err := json.Marshal(user)
	if err != nil {
		return response.StdResponse{}, fmt.Errorf("failed to marshal user data: %v", err)
	}

	rmp := producer.NewProducer()
	producerService := producer.NewProducerService(rmp)

	if err := producerService.Initialize(); err != nil {
		log.Println("Failed to initialize producer service:", err)
		return response.StdResponse{}, err
	}

	if err := producerService.Publish(jsonData, "UserCreatedTask"); err != nil {
		log.Println("Error publishing user data:", err)
		return response.StdResponse{}, err
	}

	userResponse := response.StdResponse{
		ID:          user.Id,
		Name:        user.Name,
		AuthID:      user.AuthID,
		Description: user.Description,
		Tag:         user.Tag,
		Student:     user.Student,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	return userResponse, nil
}

func (um *UserMgr) UpdateUser(id primitive.ObjectID, req request.StdRequest) (response.StdResponse, error) {
	mongoCollection := db.GetDB().Database("Courses").Collection("students")
	var existingUser models.Std
	err := mongoCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&existingUser)
	if err != nil {
		log.Printf("Failed to find user with ID %v: %v", id, err)
		return response.StdResponse{}, fmt.Errorf("failed to find user in MongoDB")
	}
	update := bson.M{
		"$set": bson.M{
			"auth_id":     req.AuthID,
			"name":        req.Name,
			"description": req.Description,
			"tag":         req.Tag,
			"student":     req.Student,
			"updated_at":  time.Now(),
		},
		"$setOnInsert": bson.M{
			"created_at": existingUser.CreatedAt,
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedUser models.Std
	err = mongoCollection.FindOneAndUpdate(context.Background(), bson.M{"_id": id}, update, opts).Decode(&updatedUser)
	if err != nil {
		log.Printf("Failed to update user with ID %v: %v", id, err)
		return response.StdResponse{}, fmt.Errorf("failed to update user in MongoDB")
	}

	log.Println("User successfully updated in MongoDB")
	return response.StdResponse{
		ID:          updatedUser.Id,
		Name:        updatedUser.Name,
		AuthID:      updatedUser.AuthID,
		Description: updatedUser.Description,
		Tag:         updatedUser.Tag,
		Student:     updatedUser.Student,
		CreatedAt:   existingUser.CreatedAt,
		UpdatedAt:   updatedUser.UpdatedAt,
		DeletedAt:   nil,
	}, nil
}

func (um *UserMgr) Getall(req request.StdRequest) ([]response.StdResponse, error) {

	mongoCollection := db.GetDB().Database("Courses").Collection("students")
	filter := bson.M{}

	if req.Name != "" {
		filter["name"] = req.Name
	}
	if len(req.Tag) > 0 {
		filter["tag"] = bson.M{"$in": req.Tag}
	}
	if req.Student {
		filter["student"] = req.Student
	}
	cursor, err := mongoCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %v", err)
	}
	defer cursor.Close(context.Background())

	// totalDocuments, err := mongoCollection.CountDocuments(context.Background())
	// 	if err != nil {
	// 		log.Printf("MongoDB count error: %v\n", err)
	// 		return nil,  fmt.Errorf("failed to count users: %v", err)
	// 	}

	var newUsers []models.Std
	if err = cursor.All(context.Background(), &newUsers); err != nil {
		return nil, fmt.Errorf("failed to decode users: %v", err)
	}

	userResponses := make([]response.StdResponse, len(newUsers))
	for i, user := range newUsers {
		userResponses[i] = response.StdResponse{
			ID:          user.Id,
			AuthID:      user.AuthID,
			Name:        user.Name,
			Description: user.Description,
			Tag:         user.Tag,
			Student:     user.Student,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			DeletedAt:   user.DeletedAt,
		}
	}

	return userResponses, nil
}
func (um *UserMgr) DeleteUser(id primitive.ObjectID) error {
	mongoCollection := db.GetDB().Database("Courses").Collection("students")
	filter := bson.M{"_id": id}

	result, err := mongoCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Printf("MongoDB deletion error: %v\n", err)
		return fmt.Errorf("failed to delete user from MongoDB")
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no user found with id %s", id.Hex())
	}

	log.Println("User successfully deleted from MongoDB")
	return nil
}
