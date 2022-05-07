package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ArtyomArtamonov/msg/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RoomStore interface {
	Add(ctx context.Context, room *model.Room) error
	Get(ctx context.Context, id uuid.UUID) (*model.Room, error)
	FindDialogRoom(ctx context.Context, userId1, userId2 uuid.UUID) (*model.Room, error)
	UsersInRoom(ctx context.Context, id uuid.UUID) ([]model.User, error)
	FindByIds(ctx context.Context, usernames ...string) ([]model.User, error)
	ListRooms(ctx context.Context, userId uuid.UUID, lastMessageDate time.Time, pageSize int) ([]model.Room, error)
	ListRoomsFirst(ctx context.Context, userId uuid.UUID, pageSize int) ([]model.Room, error)
}

type PostgresRoomStore struct {
	db *pgxpool.Pool
}

func NewPostgresRoomStore(db *pgxpool.Pool) *PostgresRoomStore {
	return &PostgresRoomStore{
		db: db,
	}
}

func (s *PostgresRoomStore) Add(ctx context.Context, room *model.Room) error {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return err
	}
	tx, err := conn.BeginTx(context.TODO(), pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, "INSERT INTO rooms(id, name, created_at, dialog_room, last_message_time) VALUES($1, $2, $3, $4, $5) RETURNING id",
		room.Id, room.Name, room.CreatedAt, room.DialogRoom, room.LastMessageTime).Scan(&room.Id)
	if err != nil {
		return err
	}

	users, err := s.FindByIds(ctx, room.Users...)
	if err != nil {
		return err
	}

	for _, user := range users {
		var room_id uuid.UUID
		err = tx.QueryRow(ctx, "INSERT INTO user_in_room(room_id, user_id) VALUES($1, $2) RETURNING room_id",
			room.Id, user.Id).Scan(&room_id)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (s *PostgresRoomStore) Get(ctx context.Context, id uuid.UUID) (*model.Room, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	var room *model.Room
	err = conn.QueryRow(ctx, "SELECT id, name, created_at, dialog_room, last_message_time FROM rooms WHERE id=$1",
		id).Scan(&room.Id, &room.Name, &room.CreatedAt, &room.DialogRoom, &room.LastMessageTime)

	return room, err
}

func (s *PostgresRoomStore) FindDialogRoom(ctx context.Context, userId1, userId2 uuid.UUID) (*model.Room, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	var room model.Room
	err = conn.QueryRow(ctx,
		`
		SELECT DISTINCT rooms.id, rooms.* FROM rooms
		INNER JOIN user_in_room ON rooms.id=user_in_room.room_id
		WHERE rooms.dialog_room=TRUE AND (user_in_room.user_id=$1 OR user_in_room.user_id=$2)
	`, userId1, userId2).Scan(&room.Id, &room.Id, &room.Name, &room.CreatedAt, &room.DialogRoom, &room.LastMessageTime)

	if err != nil {
		return nil, err
	}

	return &room, status.Error(codes.AlreadyExists, "room already exists")
}

func (s *PostgresRoomStore) UsersInRoom(ctx context.Context, id uuid.UUID) ([]model.User, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	rows, err := conn.Query(ctx, "SELECT * FROM user_in_room WHERE room_id=$1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User

		err = rows.Scan(&user.Id, &user.Username, &user.PasswordHash, &user.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *PostgresRoomStore) FindByIds(ctx context.Context, ids ...string) ([]model.User, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	q := "SELECT id, username, password_hash, role FROM users WHERE"
	for _, name := range ids {
		q += fmt.Sprintf(" id='%s' OR", name)
	}
	q = q[:len(q)-3] // remove last OR
	rows, err := conn.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []model.User{}
	for rows.Next() {
		var user model.User
		err = rows.Scan(&user.Id, &user.Username, &user.PasswordHash, &user.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(users) != len(ids) {
		return nil, status.Error(codes.InvalidArgument, "user is unknown")
	}

	return users, nil
}

func (s *PostgresRoomStore) ListRooms(ctx context.Context, userId uuid.UUID, lastMessageDate time.Time, pageSize int) ([]model.Room, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	rows, err := conn.Query(ctx,
		`
		SELECT DISTINCT rooms.id, rooms.* FROM rooms
		INNER JOIN user_in_room ON rooms.id=user_in_room.room_id
		WHERE user_in_room.user_id=$1 AND rooms.last_message_time<=$2
		ORDER BY rooms.last_message_time DESC
		LIMIT $3
		`, userId, lastMessageDate, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []model.Room
	for rows.Next() {
		var room model.Room
		err = rows.Scan(&room.Id, &room.Id, &room.Name, &room.CreatedAt, &room.DialogRoom, &room.LastMessageTime)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (s *PostgresRoomStore) ListRoomsFirst(ctx context.Context, userId uuid.UUID, pageSize int) ([]model.Room, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	rows, err := conn.Query(ctx,
		`
		SELECT DISTINCT rooms.id, rooms.* FROM rooms
		INNER JOIN user_in_room ON rooms.id=user_in_room.room_id
		WHERE user_in_room.user_id=$1
		ORDER BY rooms.last_message_time DESC
		LIMIT $2
		`, userId, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []model.Room
	for rows.Next() {
		var room model.Room
		err = rows.Scan(&room.Id, &room.Id, &room.Name, &room.CreatedAt, &room.DialogRoom, &room.LastMessageTime)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}
