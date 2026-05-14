package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/MikleMkhlv/QuizPlatform/internal/domain"
)

const (
	gameKeyPrefix = "game:"
	gameTTL       = 24 * time.Hour
)

type RedisGameRepository struct {
	rClient *redis.Client
}

func NewRedisRepository(rClient *redis.Client) *RedisGameRepository {
	return &RedisGameRepository{
		rClient: rClient,
	}
}

func gameKey(roomID uuid.UUID) string {
	return fmt.Sprintf("%s%s", gameKeyPrefix, roomID.String())
}

func (r *RedisGameRepository) Save(ctx context.Context, gameState *domain.GameState) error {
	data, err := gameState.ConvertToByte()
	if err != nil {
		return err
	}
	err = r.rClient.Set(ctx, gameKey(gameState.RoomID), data, gameTTL).Err()
	if err != nil {
		return fmt.Errorf("redis set: %w", err)
	}
	return nil
}

func (r *RedisGameRepository) Get(ctx context.Context, roomID uuid.UUID) (*domain.GameState, error) {
	data, err := r.rClient.Get(ctx, gameKey(roomID)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("redis get error: %w", err)
	}
	state, err := domain.GameStateFromBytes(data)
	if err != nil {
		return nil, err
	}
	return state, nil
}

func (r *RedisGameRepository) Delete(ctx context.Context, roomID uuid.UUID) error {
	deleted, err := r.rClient.Del(ctx, gameKey(roomID)).Result()
	if err != nil {
		return fmt.Errorf("redis del: %w", err)
	}
	if deleted == 0 {
		return fmt.Errorf("game state not found for room %s", roomID)
	}
	return nil
}
