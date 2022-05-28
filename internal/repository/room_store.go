package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ArtyomArtamonov/msg/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RoomStore interface {
	AddAndSendMessage(ctx context.Context, room *model.Room, message *model.Message) (*model.Room, error)
	Add(ctx context.Context, room *model.Room) error
	Get(ctx context.Context, id uuid.UUID) (*model.Room, error)
	FindDialogRoom(ctx context.Context, userId1, userId2 uuid.UUID) (*model.Room, error)
	UsersInRoom(ctx context.Context, id uuid.UUID) ([]uuid.UUID, error)
	FindByIds(ctx context.Context, ids ...uuid.UUID) ([]model.User, error)
	ListRooms(ctx context.Context, userId uuid.UUID, lastMessageDate time.Time, pageSize int) ([]model.Room, error)
	ListRoomsFirst(ctx context.Context, userId uuid.UUID, pageSize int) ([]model.Room, error)
}

type PostgresRoomStore struct {
	db *sqlx.DB
}

func NewPostgresRoomStore(db *sqlx.DB) *PostgresRoomStore {
	return &PostgresRoomStore{
		db: db,
	}
}

func (s *PostgresRoomStore) AddAndSendMessage(ctx context.Context, room *model.Room, message *model.Message) (*model.Room, error) {
	r, err := s.FindDialogRoom(ctx, room.UserIds[0], room.UserIds[1])
	if err != nil {
		return r, err
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = tx.GetContext(ctx, &room.Id, "INSERT INTO rooms(id, name, created_at, dialog_room, last_message_time) VALUES($1, $2, $3, $4, $5) RETURNING id",
		room.Id, room.Name, room.CreatedAt, room.DialogRoom, room.LastMessageTime)
	if err != nil {
		return nil, err
	}

	for _, userId := range room.UserIds {
		var room_id uuid.UUID
		err = tx.Get(&room_id, "INSERT INTO user_in_room(room_id, user_id) VALUES($1, $2) RETURNING room_id",
			room.Id, userId)
		if err != nil {
			return nil, err
		}
	}

	message.RoomId = room.Id
	var messageId uuid.UUID
	err = tx.GetContext(ctx, &messageId, "INSERT INTO messages(id, room_id, user_id, text, created_at) VALUES($1, $2, $3, $4, $5) RETURNING id",
		message.Id, message.RoomId, message.UserId, message.Text, message.CreatedAt)
	if err != nil {
		return nil, err
	}

	message.Id = messageId

	return room, tx.Commit()
}

func (s *PostgresRoomStore) Add(ctx context.Context, room *model.Room) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.GetContext(ctx, &room.Id, "INSERT INTO rooms(id, name, created_at, dialog_room, last_message_time) VALUES($1, $2, $3, $4, $5) RETURNING id",
		room.Id, room.Name, room.CreatedAt, room.DialogRoom, room.LastMessageTime)
	if err != nil {
		return err
	}

	for _, userId := range room.UserIds {
		var room_id uuid.UUID
		err = tx.GetContext(ctx, &room_id, "INSERT INTO user_in_room(room_id, user_id) VALUES($1, $2) RETURNING room_id",
			room.Id, userId)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *PostgresRoomStore) Get(ctx context.Context, id uuid.UUID) (*model.Room, error) {
	room := new(model.Room)
	err := s.db.GetContext(ctx, room, "SELECT id, name, created_at, dialog_room, last_message_time FROM rooms WHERE id=$1",
		id)

	return room, err
}

func (s *PostgresRoomStore) FindDialogRoom(ctx context.Context, userId1, userId2 uuid.UUID) (*model.Room, error) {
	room := new(model.Room)
	err := s.db.GetContext(
		ctx,
		room,
		`
		SELECT DISTINCT rooms.id, rooms.* FROM rooms
		INNER JOIN user_in_room ON rooms.id=user_in_room.room_id
		WHERE rooms.dialog_room=TRUE AND (user_in_room.user_id=$1 OR user_in_room.user_id=$2)
		`, userId1, userId2)

	if err != nil {
		return nil, err
	}

	return room, status.Error(codes.AlreadyExists, "room already exists")
}

func (s *PostgresRoomStore) UsersInRoom(ctx context.Context, id uuid.UUID) ([]uuid.UUID, error) {
	userIds := []string{}
	err := s.db.SelectContext(ctx, &userIds, "SELECT user_id FROM user_in_room WHERE room_id=$1", id)
	if err != nil {
		return nil, err
	}

	var userUUIDs []uuid.UUID
	for _, userId := range userIds {
		id, err := uuid.Parse(userId)
		if err != nil {
			logrus.Errorf("could not parse uuid: %v", err)
		}
		userUUIDs = append(userUUIDs, id)
	}

	return userUUIDs, nil
}

func (s *PostgresRoomStore) FindByIds(ctx context.Context, ids ...uuid.UUID) ([]model.User, error) {
	// TODO: rewrite with more safety. 
	// We can pass array with $1
	q := "SELECT * FROM users WHERE"
	for _, name := range ids {
		q += fmt.Sprintf(" id='%s' OR", name)
	}
	q = q[:len(q)-3] // remove last OR

	users := []model.User{}
	err := s.db.SelectContext(ctx, &users, q)
	if err != nil {
		return nil, err
	}

	if len(users) != len(ids) {
		return nil, status.Error(codes.InvalidArgument, "user is unknown")
	}

	return users, nil
}

func (s *PostgresRoomStore) ListRooms(ctx context.Context, userId uuid.UUID, lastMessageDate time.Time, pageSize int) ([]model.Room, error) {
	res := []model.Room{}
	// This SQL Query has a workaround in it: first and second returned rows are actually the same.
	// But removing first row somehow ruins everything.
	// We are using unsafe here only to silently not map first row
	udb := s.db.Unsafe()
	err := udb.SelectContext(
		ctx,
		&res,
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

	return res, nil
}

func (s *PostgresRoomStore) ListRoomsFirst(ctx context.Context, userId uuid.UUID, pageSize int) ([]model.Room, error) {
	res := []model.Room{}
	// This SQL Query has a workaround in it: first and second returned rows are actually the same.
	// But removing first row somehow ruins everything.
	// We are using unsafe here only to silently not map first row
	udb := s.db.Unsafe()
	err := udb.SelectContext(
		ctx,
		&res,
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

	return res, nil
}
