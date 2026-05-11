package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourname/quiz-platform/internal/domain"
)

type PostgresRoomRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRoomRepository(pool *pgxpool.Pool) *PostgresRoomRepository {
	return &PostgresRoomRepository{
		pool: pool,
	}
}

func (p *PostgresRoomRepository) Create(ctx context.Context, room *domain.Room) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	queryCreateRoom := `INSERT INTO rooms (id, host_id, code, status, max_plauer, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.Exec(ctx, queryCreateRoom, &room.ID, &room.HostID, &room.Code, &room.Status, &room.MaxPlauer, &room.CreatedAt)
	if err != nil {
		return err
	}

	newRoomPlauer := *domain.NewRoomPlauer(room.ID, room.HostID)
	queryAddNewRoomPlauer := `INSERT INTO room_players (room_id, plauer_id, joined_at) VALUES ($1, $2, $3)`
	_, err = tx.Exec(ctx, queryAddNewRoomPlauer, newRoomPlauer.RoomId, newRoomPlauer.PlauerId, newRoomPlauer.JoinedAt)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresRoomRepository) GetRoomById(ctx context.Context, roomId uuid.UUID) (*domain.Room, error) {
	query := `SELECT id, host_id, code, status, max_plauer, created_at FROM room_players WHERE id = $1`
	room := &domain.Room{}
	err := p.pool.QueryRow(ctx, query, roomId).Scan(&room.ID, &room.HostID, &room.Code, &room.Status, &room.MaxPlauer, &room.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return room, nil
}
