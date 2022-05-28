package repository

import (
	"context"
	"time"

	"github.com/ArtyomArtamonov/msg/internal/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type MessageStore interface {
	SendMessage(ctx context.Context, message *model.Message) error
	ListMessages(ctx context.Context, id uuid.UUID, createdAt time.Time, pageSize int) ([]model.Message, error)
	ListMessagesFirst(ctx context.Context, id uuid.UUID, pageSize int) ([]model.Message, error)
}

type PostgresMessageStore struct {
	db *sqlx.DB
}

func NewPostgresMessageStore(db *sqlx.DB) *PostgresMessageStore {
	return &PostgresMessageStore{
		db: db,
	}
}

func (s *PostgresMessageStore) ListMessages(ctx context.Context, chatId uuid.UUID, createdAt time.Time, pageSize int) ([]model.Message, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, _, err := psql.
		Select("id, room_id, user_id, text, created_at").
		From("messages").
		Where(sq.And{
			sq.Eq{"room_id": chatId},
			sq.Lt{"created_at": createdAt},
		}).
		OrderBy("created_at DESC").
		Limit(uint64(pageSize)).
		ToSql()
	if err != nil {
		return nil, err
	}

	messages := []model.Message{}
	err = s.db.SelectContext(ctx, &messages, sql, chatId, createdAt)
	if err != nil {
		return nil, err
	}

	return messages, err
}

func (s *PostgresMessageStore) ListMessagesFirst(ctx context.Context, chatId uuid.UUID, pageSize int) ([]model.Message, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, _, err := psql.
		Select("id, room_id, user_id, text, created_at").
		From("messages").
		Where(sq.Eq{"room_id": chatId}).
		OrderBy("created_at DESC").
		Limit(uint64(pageSize)).
		ToSql()
	if err != nil {
		return nil, err
	}

	messages := []model.Message{}
	err = s.db.SelectContext(ctx, &messages, sql, chatId)
	if err != nil {
		return nil, err
	}

	return messages, err
}

func (s *PostgresMessageStore) SendMessage(ctx context.Context, message *model.Message) error {
	var messageId uuid.UUID
	err := s.db.GetContext(ctx, &messageId, "INSERT INTO messages(id, room_id, user_id, text, created_at) VALUES($1, $2, $3, $4, $5) RETURNING id",
		message.Id, message.RoomId, message.UserId, message.Text, message.CreatedAt)
	message.Id = messageId

	return err
}
