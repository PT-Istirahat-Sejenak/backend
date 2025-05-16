package usecase

import (
	"backend/internal/entity"
	"backend/internal/repository"
	"context"
)

type messageUseCase struct {
	messageRepo repository.MessageRepository
}

func NewMessageUseCase(messageRepo repository.MessageRepository) MessageUseCase {
	return &messageUseCase{
		messageRepo: messageRepo,
	}
}

// SendMessage implements MessageUseCase.
func (m *messageUseCase) SaveMessage(ctx context.Context, senderID uint, receiverID uint, content string) error {
	message := &entity.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
	}
	return m.messageRepo.SaveMessage(ctx, message)
}

// GetMessagesByUserID implements MessageUseCase.
func (m *messageUseCase) GetMessagesByUserID(ctx context.Context, userID1 uint, userID2 uint, limit int, offset int) ([]entity.Message, error) {
	return m.messageRepo.GetMessagesByUserID(ctx, userID1, userID2, limit, offset)
}

// GetUndeliveredMessages implements MessageUseCase.
func (m *messageUseCase) GetUndeliveredMessages(ctx context.Context, receiverID uint) ([]entity.Message, error) {
	return m.messageRepo.GetUndeliveredMessages(ctx, receiverID)
}

// MarkMessageAsDelivered implements MessageUseCase.
func (m *messageUseCase) MarkMessageAsDelivered(ctx context.Context, messageID uint) error {
	return m.messageRepo.MarkMessageAsDelivered(ctx, messageID)
}
