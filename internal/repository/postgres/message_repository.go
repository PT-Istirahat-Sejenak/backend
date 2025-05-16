package postgres

import (
	"backend/internal/entity"
	"context"
	"database/sql"
	"time"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{
		db: db,
	}
}

func (r *MessageRepository) SaveMessage(ctx context.Context, message *entity.Message) error {
	query := `
	INSERT INTO messages (sender_id, receiver_id, content, is_delivered, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id;
	`

	now := time.Now()
	message.CreatedAt = now
	message.UpdatedAt = now

	// Default is_delivered to false (message not delivered yet)
	isDelivered := false

	var id uint
	err := r.db.QueryRowContext(
		ctx, query,
		message.SenderID,
		message.ReceiverID,
		message.Content,
		isDelivered,
		message.CreatedAt,
		message.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return err
	}

	message.ID = id
	return nil
}

// GetUndeliveredMessages retrieves all undelivered messages for a specific receiver
func (r *MessageRepository) GetUndeliveredMessages(ctx context.Context, receiverID uint) ([]entity.Message, error) {
	query := `
	SELECT id, sender_id, receiver_id, content, created_at, updated_at
	FROM messages
	WHERE receiver_id = $1 AND is_delivered = false
	ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, receiverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []entity.Message
	for rows.Next() {
		var msg entity.Message
		if err := rows.Scan(
			&msg.ID,
			&msg.SenderID,
			&msg.ReceiverID,
			&msg.Content,
			&msg.CreatedAt,
			&msg.UpdatedAt,
		); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

// MarkMessageAsDelivered marks a message as delivered
func (r *MessageRepository) MarkMessageAsDelivered(ctx context.Context, messageID uint) error {
	query := `
	UPDATE messages
	SET is_delivered = true, updated_at = $1
	WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, time.Now(), messageID)
	return err
}

// GetMessagesByUserID retrieves messages between two users
func (r *MessageRepository) GetMessagesByUserID(ctx context.Context, userID1, userID2 uint, limit, offset int) ([]entity.Message, error) {
	query := `
	SELECT id, sender_id, receiver_id, content, is_delivered, created_at, updated_at
	FROM messages
	WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1)
	ORDER BY created_at DESC
	LIMIT $3 OFFSET $4
	`

	rows, err := r.db.QueryContext(ctx, query, userID1, userID2, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []entity.Message
	for rows.Next() {
		var msg entity.Message
		var isDelivered bool

		if err := rows.Scan(
			&msg.ID,
			&msg.SenderID,
			&msg.ReceiverID,
			&msg.Content,
			&isDelivered,
			&msg.CreatedAt,
			&msg.UpdatedAt,
		); err != nil {
			return nil, err
		}

		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
