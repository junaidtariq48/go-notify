package controllers

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
)

// InitRouter initializes the router and sets up the routes
func InitRouter(ctx context.Context, channel *amqp.Channel, db *mongo.Client) *mux.Router {
	router := mux.NewRouter()

	// Define the notification creation route
	router.HandleFunc("/api/notifications", func(w http.ResponseWriter, r *http.Request) {
		CreateNotification(w, r, ctx, channel, db)
	}).Methods("POST")

	// You can add more routes here for other functionalities (e.g., fetching notifications, updating, etc.)

	return router
}
