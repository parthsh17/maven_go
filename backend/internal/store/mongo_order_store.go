package store

import (
	"context"
	"fmt"
	"maven/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoOrderStore struct {
	collection *mongo.Collection
}

func NewMongoOrderStore(db *mongo.Database) OrderStore {
	return &mongoOrderStore{
		collection: db.Collection("orders"),
	}
}

func (s *mongoOrderStore) AddOrder(order *models.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := s.collection.InsertOne(ctx, order)
	return err
}

func (s *mongoOrderStore) GetOrder(id string) (*models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var order models.Order
	err := s.collection.FindOne(ctx, bson.M{"id": id}).Decode(&order)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}
	return &order, nil
}

func (s *mongoOrderStore) GetAllOrders() []*models.Order {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil
	}
	defer cursor.Close(ctx)

	var orders []*models.Order
	if err = cursor.All(ctx, &orders); err != nil {
		return nil
	}
	return orders
}

func (s *mongoOrderStore) UpdateState(id, state, msg string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	order, err := s.GetOrder(id)
	if err != nil {
		return err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	event := models.OrderEvent{
		OrderID:       id,
		PreviousState: order.State,
		NewState:      state,
		Timestamp:     now,
		Message:       msg,
	}

	update := bson.M{
		"$set": bson.M{
			"state":      state,
			"updated_at": now,
		},
		"$push": bson.M{"events": event},
	}

	_, err = s.collection.UpdateOne(ctx, bson.M{"id": id}, update)
	return err
}

func (s *mongoOrderStore) UpdateSlippage(id string, slippage float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": bson.M{"slippage": slippage}})
	return err
}

func (s *mongoOrderStore) IncrementRetry(id string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, _ = s.collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$inc": bson.M{"retry_count": 1}})
}

func (s *mongoOrderStore) GetEvents(id string) ([]models.OrderEvent, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var order models.Order
	err := s.collection.FindOne(ctx, bson.M{"id": id}).Decode(&order)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}
	return order.Events, nil
}
