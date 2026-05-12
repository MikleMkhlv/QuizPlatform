package postgres

import (
	"context"
	"errors"
	"fmt"

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

	queryCreateRoom := `
						INSERT INTO rooms (id, host_id, code, status, max_players, created_at) 
						VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = tx.Exec(ctx, queryCreateRoom, room.ID, room.HostID, room.Code, room.Status, room.MaxPlayer, room.CreatedAt)
	if err != nil {
		return err
	}

	newPlayer := domain.NewRoomPlayer(room.ID, room.HostID)
	queryAddNewRoomPlauer := `
							INSERT INTO room_players (room_id, user_id, joined_at) 
							VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(ctx, queryAddNewRoomPlauer, newPlayer.RoomID, newPlayer.PlayerID, newPlayer.JoinedAt)
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
	query := `
			SELECT id, host_id, code, status, max_players, created_at 
			FROM rooms 
			WHERE id = $1
	`
	room := &domain.Room{}
	err := p.pool.QueryRow(ctx, query, roomId).Scan(&room.ID, &room.HostID, &room.Code, &room.Status, &room.MaxPlayer, &room.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("not found room with id: %s", roomId.String())
	}
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (p *PostgresRoomRepository) GetRoomByCode(ctx context.Context, roomCode string) (*domain.Room, error) {
	query := `
			SELECT id, host_id, code, status, max_players, created_at 
			FROM rooms 
			WHERE code = $1`
	room := &domain.Room{}
	err := p.pool.QueryRow(ctx, query, roomCode).Scan(&room.ID, &room.HostID, &room.Code, &room.Status, &room.MaxPlayer, &room.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (p *PostgresRoomRepository) AddPlayer(ctx context.Context, newPlayer *domain.RoomPlayer) error {
	query := `
			INSERT INTO room_players (room_id, user_id, joined_at) 
			VALUES ($1, $2, $3)
	`
	_, err := p.pool.Exec(ctx, query, newPlayer.RoomID, newPlayer.PlayerID, newPlayer.JoinedAt)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresRoomRepository) GetPlayersFromRoom(ctx context.Context, roomId uuid.UUID) ([]*domain.User, error) {
	query := `
			SELECT u.id, u.name, u.username, u.rating, u.created_at 
			FROM users u
			JOIN room_players rp ON rp.user_id = u.id
			WHERE rp.room_id = $1
			ORDER BY rp.joined_at ASC
	`
	var users []*domain.User
	rows, err := p.pool.Query(ctx, query, roomId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Username, &user.Rating, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return users, nil
}

func (p *PostgresRoomRepository) UpdateRoomStatus(ctx context.Context, roomId uuid.UUID, status domain.RoomStatus) error {
	query := `
			UPDATE rooms 
			SET status = $1 
			WHERE id = $2
	`
	_, err := p.pool.Exec(ctx, query, status, roomId)
	if err != nil {
		return err
	}
	return nil
}
