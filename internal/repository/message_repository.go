package repository

import (
	"time"

	"github.com/ArtyomArtamonov/msg/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	sq "github.com/Masterminds/squirrel"
)

type MessageStore interface {
	SendMessage(message *model.Message) error
	ListMessages(id uuid.UUID, createdAt time.Time, pageSize int) ([]model.Message, error)
	ListMessagesFirst(id uuid.UUID, pageSize int) ([]model.Message, error)
}

type PostgresMessageStore struct {
	db *sqlx.DB
}

func NewPostgresMessageStore(db *sqlx.DB) *PostgresMessageStore {
	return &PostgresMessageStore{
		db: db,
	}
}

func (s *PostgresMessageStore) ListMessages(id uuid.UUID, createdAt time.Time, pageSize int) ([]model.Message, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, _, err := psql.
		Select("id, room_id, user_id, text, created_at").
		From("messages").
		Where(sq.And{
			sq.Eq{"user_id": id},
			sq.Lt{"created_at": createdAt},
		}).
		OrderBy("created_at DESC").
		Limit(uint64(pageSize)).
		ToSql()
	if err != nil {
		return nil, err
	}

	messages := []model.Message{}
	err = s.db.Select(&messages, sql, id, createdAt)
	if err != nil {
		return nil, err
	}

	return messages, err
}

func (s *PostgresMessageStore) ListMessagesFirst(id uuid.UUID, pageSize int) ([]model.Message, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, _, err := psql.
		Select("id, room_id, user_id, text, created_at").
		From("messages").
		Where(sq.Eq{"user_id": id}).
		OrderBy("created_at DESC").
		Limit(uint64(pageSize)).
		ToSql()
	if err != nil {
		return nil, err
	}

	messages := []model.Message{}
	err = s.db.Select(&messages, sql, id)
	if err != nil {
		return nil, err
	}

	return messages, err
}

func (s *PostgresMessageStore) SendMessage(message *model.Message) error {
	var messageId uuid.UUID
	err := s.db.Get(&messageId, "INSERT INTO messages(id, room_id, user_id, text, created_at) VALUES($1, $2, $3, $4, $5) RETURNING id",
		message.Id, message.RoomId, message.UserId, message.Text, message.CreatedAt)
	message.Id = messageId

	return err
}
