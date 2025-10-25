package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"api/helpers"
	"api/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client

func GetClient() *mongo.Client {
	return client
}

const (
	defaultPort = ":8081"
)

// MongoDB Variables
var (
	Client        *mongo.Client
	URLCollection *mongo.Collection
)

// URL Struct
type URL struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ShortCode     string             `bson:"short_code" json:"short_code"`
	LongURL       string             `bson:"long_url" json:"long_url"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	Clicks        int64              `bson:"clicks" json:"clicks"`
	LastClickedAt time.Time          `bson:"last_clicked_at" json:"last_clicked_at"`
}

// ErrorResponse struct
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse struct
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Get environment variable or default
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// GenerateShortCode
func generateShortCode(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}

func SetupDatabase() (*mongo.Client, error) {
	mongoURI := getEnvOrDefault("MONGO_URI", "")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	return client, nil
}

// Handle Shorten URL
func handleShortenURL(c *gin.Context) {
	var request struct {
		LongURL   string  `json:"long_url" binding:"required,url"`
		ShortCode *string `json:"short_code"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Error binding JSON in handleShortenURL: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	collection := GetClient().Database("goapp").Collection("urls")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var shortCode string
	if request.ShortCode != nil && *request.ShortCode != "" {
		// Check if custom short code already exists
		filter := bson.M{"short_code": *request.ShortCode}
		count, err := collection.CountDocuments(ctx, filter)
		if err != nil {
			log.Printf("Error checking custom short code existence: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to check custom short code"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, ErrorResponse{Error: "Custom short code already in use"})
			return
		}
		shortCode = *request.ShortCode
	} else {
		// Generate a unique short code
		for {
			shortCode, err := helpers.GenerateShortCode(6)
			if err != nil {
				log.Printf("Error generating short code: %v", err)
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to generate short code"})
				return
			}
			filter := bson.M{"short_code": shortCode}
			count, err := collection.CountDocuments(ctx, filter)
			if err != nil {
				log.Printf("Error checking short code existence: %v", err)
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to generate short code"})
				return
			}
			if count == 0 {
				break
			}
		}
	}

	// Create a new URL document
	newURL := models.URL{
		ID:        primitive.NewObjectID(),
		LongURL:   request.LongURL,
		ShortCode: shortCode,
		ShortURL:  fmt.Sprintf("http://localhost:8081/%s", shortCode),
		Clicks:    0,
		CreatedAt: time.Now(),
	}

	// Insert the new URL into the database
	_, err := collection.InsertOne(ctx, newURL)
	if err != nil {
		log.Printf("Error inserting URL: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to shorten URL"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: gin.H{"short_url": newURL.ShortURL}})
}

// Handle Redirect
func handleRedirect(c *gin.Context) {
	shortCode := c.Param("shortCode")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var url URL
	err := URLCollection.FindOne(ctx, bson.M{"short_code": shortCode}).Decode(&url)
	if err != nil {
		log.Printf("Short URL %s not found: %v", shortCode, err)
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Short URL not found"})
		return
	}

	// Update click count and last clicked at
	_, err = URLCollection.UpdateOne(
		ctx,
		bson.M{"_id": url.ID},
		bson.M{"$inc": bson.M{"clicks": 1}, "$set": bson.M{"last_clicked_at": time.Now()}},
	)
	if err != nil {
		log.Printf("Failed to update click count for %s: %v", shortCode, err)
	}

	c.Redirect(http.StatusMovedPermanently, url.LongURL)
}

// Handle Get Stats
func handleGetStats(c *gin.Context) {
	shortCode := c.Param("shortCode")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var url URL
	err := URLCollection.FindOne(ctx, bson.M{"short_code": shortCode}).Decode(&url)
	if err != nil {
		log.Printf("Short URL %s for stats not found: %v", shortCode, err)
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Short URL not found"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "URL statistics retrieved successfully",
		Data: gin.H{
			"long_url":        url.LongURL,
			"short_code":      url.ShortCode,
			"clicks":          url.Clicks,
			"created_at":      url.CreatedAt,
			"last_clicked_at": url.LastClickedAt,
		},
	})
}

// handleListURLs retrieves all shortened URLs from the database
func handleListURLs(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var urls []URL
	cursor, err := URLCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URLs"})
		return
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &urls); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode URLs"})
		return
	}

	c.JSON(http.StatusOK, urls)
}

// handleDeleteURL deletes a shortened URL from the database
func handleDeleteURL(c *gin.Context) {
	shortCode := c.Param("shortCode")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"short_code": shortCode}
	result := URLCollection.FindOneAndDelete(ctx, filter)

	if result.Err() != nil {
		fmt.Printf("Error deleting URL: %v", result.Err())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "URL deleted successfully"})
}

// handleEditURL updates an existing shortened URL in the database
func handleEditURL(c *gin.Context) {
	shortCode := c.Param("shortCode")

	var requestBody struct {
		LongURL string `json:"long_url"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if requestBody.LongURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Long URL cannot be empty"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"short_code": shortCode}
	update := bson.M{"$set": bson.M{"long_url": requestBody.LongURL}}

	result, err := URLCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update URL"})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found or no changes made"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "URL updated successfully"})
}

func main() {
	var err error
	Client, err = SetupDatabase()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := Client.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	URLCollection = Client.Database("goapp").Collection("urls")

	// Initialize Gin router
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://your-production-domain.com", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Public routes
	r.POST("/shorten", handleShortenURL)
	r.GET("/:shortCode", handleRedirect)
	r.GET("/:shortCode/stats", handleGetStats)
	r.GET("/urls", handleListURLs)
	r.DELETE("/urls/:shortCode", handleDeleteURL)
	r.PUT("/urls/:shortCode", handleEditURL)

	port := getEnvOrDefault("PORT", defaultPort)
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	fmt.Printf("Server running on port %s\n", port)
}
