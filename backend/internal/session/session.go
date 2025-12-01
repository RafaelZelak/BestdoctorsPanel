package session

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store struct {
	client *redis.Client
	prefix string
	ttl    time.Duration
}

type SessionData struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func NewStore(redisURL string) (*Store, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("invalid redis URL: %w", err)
	}

	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return &Store{
		client: client,
		prefix: "sess:",
		ttl:    24 * time.Hour,
	}, nil
}

func (s *Store) GenerateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate session ID: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (s *Store) Create(ctx context.Context, data SessionData) (string, error) {
	sessionID, err := s.GenerateSessionID()
	if err != nil {
		return "", err
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal session data: %w", err)
	}

	key := s.prefix + sessionID
	err = s.client.Set(ctx, key, jsonData, s.ttl).Err()
	if err != nil {
		return "", fmt.Errorf("failed to store session: %w", err)
	}

	return sessionID, nil
}

func (s *Store) Get(ctx context.Context, sessionID string) (*SessionData, error) {
	if sessionID == "" {
		return nil, fmt.Errorf("empty session ID")
	}

	key := s.prefix + sessionID
	val, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("session not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	var data SessionData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session data: %w", err)
	}

	s.client.Expire(ctx, key, s.ttl)

	return &data, nil
}

func (s *Store) Delete(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return nil
	}

	key := s.prefix + sessionID
	return s.client.Del(ctx, key).Err()
}

func (s *Store) DeleteByUserID(ctx context.Context, userID int) error {
	iter := s.client.Scan(ctx, 0, s.prefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		val, err := s.client.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var data SessionData
		if err := json.Unmarshal([]byte(val), &data); err != nil {
			continue
		}

		if data.UserID == userID {
			s.client.Del(ctx, key)
		}
	}
	return iter.Err()
}
