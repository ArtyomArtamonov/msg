package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ArtyomArtamonov/msg/internal/model"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RoomStore interface {
	Add(room *model.Room) error
	AddAndSendMessage(room *model.Room, message *model.Message) (*model.Room, error)
	SendMessage(roomId uuid.UUID, message *model.Message) (error)
	Get(id uuid.UUID) (*model.Room, error)
	FindDialogRoom(userId1, userId2 uuid.UUID) (*model.Room, error)
	UsersInRoom(id uuid.UUID) ([]uuid.UUID, error)
	FindByIds(ids ...uuid.UUID) ([]model.User, error)
	ListRooms(userId uuid.UUID, lastMessageDate time.Time, pageSize int) ([]model.Room, error)
	ListRoomsFirst(userId uuid.UUID, pageSize int) ([]model.Room, error)
}

type PostgresRoomStore struct {
	db *sql.DB
}

func NewPostgresRoomStore(db *sql.DB) *PostgresRoomStore {
	return &PostgresRoomStore{
		db: db,
	}
}

func (s *PostgresRoomStore) Add(room *model.Room) error {
	tx, err := s.db.BeginTx(context.TODO(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.QueryRow("INSERT INTO rooms(id, name, created_at, dialog_room, last_message_time) VALUES($1, $2, $3, $4, $5) RETURNING id",
		room.Id, room.Name, room.CreatedAt, room.DialogRoom, room.LastMessageTime).Scan(&room.Id)
	if err != nil {
		return err
	}

	for _, userId := range room.UserIds {
		var room_id uuid.UUID
		err = tx.QueryRow("INSERT INTO user_in_room(room_id, user_id) VALUES($1, $2) RETURNING room_id",
			room.Id, userId).Scan(&room_id)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *PostgresRoomStore) AddAndSendMessage(room *model.Room, message *model.Message) (*model.Room, error) {
	// Check if dialog room is already exists
	r, err := s.FindDialogRoom(room.UserIds[0], room.UserIds[1])
	if status.Code(err) == codes.AlreadyExists {
		message.RoomId = r.Id
		return r, s.SendMessage(r.Id, message)
	}

	tx, err := s.db.BeginTx(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = tx.QueryRow("INSERT INTO rooms(id, name, created_at, dialog_room, last_message_time) VALUES($1, $2, $3, $4, $5) RETURNING id",
		room.Id, room.Name, room.CreatedAt, room.DialogRoom, room.LastMessageTime).Scan(&room.Id)
	if err != nil {
		return nil, err
	}

	for _, userId := range room.UserIds {
		var room_id uuid.UUID
		err = tx.QueryRow("INSERT INTO user_in_room(room_id, user_id) VALUES($1, $2) RETURNING room_id",
			room.Id, userId).Scan(&room_id)
		if err != nil {
			return nil, err
		}
	}

	message.RoomId = room.Id
	var messageId uuid.UUID
	err = tx.QueryRow("INSERT INTO messages(id, room_id, user_id, text, created_at) VALUES($1, $2, $3, $4, $5) RETURNING id",
	message.Id, message.RoomId, message.UserId, message.Text, message.CreatedAt).Scan(&messageId)
	if err != nil {
		return nil, err
	}

	message.Id = messageId

	return room, tx.Commit()
}

func (s *PostgresRoomStore) SendMessage(roomId uuid.UUID, message *model.Message) (error) {
	var messageId uuid.UUID
	err := s.db.QueryRow("INSERT INTO messages(id, room_id, user_id, text, created_at) VALUES($1, $2, $3, $4, $5) RETURNING id",
	message.Id, message.RoomId, message.UserId, message.Text, message.CreatedAt).Scan(&messageId)
	message.Id = messageId

	return err
}

func (s *PostgresRoomStore) Get(id uuid.UUID) (*model.Room, error) {
	var room *model.Room
	err := s.db.QueryRow("SELECT id, name, created_at, dialog_room, last_message_time FROM rooms WHERE id=$1",
		id).Scan(&room.Id, &room.Name, &room.CreatedAt, &room.DialogRoom, &room.LastMessageTime)

	return room, err
}

func (s *PostgresRoomStore) FindDialogRoom(userId1, userId2 uuid.UUID) (*model.Room, error) {
	var room model.Room
	err := s.db.QueryRow(
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

func (s *PostgresRoomStore) UsersInRoom(id uuid.UUID) ([]uuid.UUID, error) {
	rows, err := s.db.Query("SELECT user_id, room_id FROM user_in_room WHERE room_id=$1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []uuid.UUID
	for rows.Next() {
		var userId string
		var roomId string

		err = rows.Scan(&userId, &roomId)
		if err != nil {
			return nil, err
		}

		id, err := uuid.Parse(userId)
		if err != nil {
			logrus.Errorf("could not parse uuid: %v", err)
		}
		users = append(users, id)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *PostgresRoomStore) FindByIds(ids ...uuid.UUID) ([]model.User, error) {
	q := "SELECT id, username, password_hash, role FROM users WHERE"
	for _, name := range ids {
		q += fmt.Sprintf(" id='%s' OR", name)
	}
	q = q[:len(q)-3] // remove last OR
	rows, err := s.db.Query(q)
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

func (s *PostgresRoomStore) ListRooms(userId uuid.UUID, lastMessageDate time.Time, pageSize int) ([]model.Room, error) {
	rows, err := s.db.Query(
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

func (s *PostgresRoomStore) ListRoomsFirst(userId uuid.UUID, pageSize int) ([]model.Room, error) {
	rows, err := s.db.Query(
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
