package models

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestURLStruct(t *testing.T) {
	// Create a sample ObjectID
	id := primitive.NewObjectID()

	// Create a sample time
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	// Create a URL instance
	url := URL{
		ID:            id,
		ShortCode:     "testcode",
		ShortURL:      "http://short.url/testcode",
		LongURL:       "http://long.url/some/path",
		CreatedAt:     now,
		ExpiresAt:     &expiresAt,
		UserID:        &id,
		Clicks:        0,
		LastClickedAt: nil,
	}

	// Test field values
	if url.ID != id {
		t.Errorf("Expected ID %v, got %v", id, url.ID)
	}
	if url.ShortCode != "testcode" {
		t.Errorf("Expected ShortCode 'testcode', got '%s'", url.ShortCode)
	}
	if url.ShortURL != "http://short.url/testcode" {
		t.Errorf("Expected ShortURL 'http://short.url/testcode', got '%s'", url.ShortURL)
	}
	if url.LongURL != "http://long.url/some/path" {
		t.Errorf("Expected LongURL 'http://long.url/some/path', got '%s'", url.LongURL)
	}
	if !url.CreatedAt.Equal(now) {
		t.Errorf("Expected CreatedAt %v, got %v", now, url.CreatedAt)
	}
	if url.ExpiresAt == nil || !url.ExpiresAt.Equal(expiresAt) {
		t.Errorf("Expected ExpiresAt %v, got %v", expiresAt, url.ExpiresAt)
	}
	if url.UserID == nil || *url.UserID != id {
		t.Errorf("Expected UserID %v, got %v", id, url.UserID)
	}
	if url.Clicks != 0 {
		t.Errorf("Expected Clicks 0, got %v", url.Clicks)
	}
	if url.LastClickedAt != nil {
		t.Errorf("Expected LastClickedAt nil, got %v", url.LastClickedAt)
	}
}