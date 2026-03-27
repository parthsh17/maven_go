package store

import (
	"context"
	"fmt"
	"log"
	"maven/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type mongoUserStore struct {
	collection *mongo.Collection
}

func NewMongoUserStore(db *mongo.Database) UserStore {
	return &mongoUserStore{
		collection: db.Collection("users"),
	}
}

func (s *mongoUserStore) Register(id, email, password string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if user already exists
	var existing models.User
	err := s.collection.FindOne(ctx, bson.M{"email": email}).Decode(&existing)
	if err == nil {
		log.Printf("[DB] User %s already exists", email)
		return nil, fmt.Errorf("email already registered")
	}
	log.Printf("[DB] User %s not found, proceeding with registration", email)

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	log.Printf("[DB] Registering new user: %s", email)

	user := models.NewUser(id, email, string(hash))
	_, err = s.collection.InsertOne(ctx, user)
	if err != nil {
		log.Printf("[DB] Failed to insert user %s: %v", email, err)
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	log.Printf("[DB] User %s successfully stored in MongoDB", email)
	return user, nil
}

func (s *mongoUserStore) Authenticate(email, password string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := s.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return &user, nil
}

func (s *mongoUserStore) GetUser(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := s.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}
