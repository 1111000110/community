package message

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MessageModel = (*customMessageModel)(nil)

type (
	// MessageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMessageModel.
	MessageModel interface {
		messageModel
		withSession(session sqlx.Session) MessageModel
		DeleteBySessionIdAndMessageID(context.Context, string, int64) error
		FindAllBySessionIdAndMessageIds(context.Context, string, []int64) ([]*Message, error)
		GetMessageList(context.Context, string, int64, int64) ([]*Message, error)
	}

	customMessageModel struct {
		*defaultMessageModel
	}
)

func (m *customMessageModel) DeleteBySessionIdAndMessageID(ctx context.Context, sessionId string, messageId int64) error {
	query := fmt.Sprintf("delete from %s where `session_id`=? and `messgage_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, sessionId, messageId)
	return err
}
func (m *customMessageModel) FindAllBySessionIdAndMessageIds(ctx context.Context, sessionId string, messageId []int64) ([]*Message, error) {
	var messages []*Message
	query := fmt.Sprintf("select * from %s where `session_id` =? and `messgage_id` IN (?)", m.table)
	err := m.conn.QueryRowCtx(ctx, &messages, query, sessionId, messageId)
	switch {
	case err == nil:
		return messages, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
func (m *customMessageModel) GetMessageList(ctx context.Context, sessionId string, offset int64, limit int64) ([]*Message, error) {
	var messages []*Message
	query := fmt.Sprintf("select * from %s where `session_id` =? offset ? and limit ?", m.table)
	err := m.conn.QueryRowCtx(ctx, &messages, query, sessionId, offset, limit)
	switch {
	case err == nil:
		return messages, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// NewMessageModel returns a model for the database table.
func NewMessageModel(conn sqlx.SqlConn) MessageModel {
	return &customMessageModel{
		defaultMessageModel: newMessageModel(conn),
	}
}

func (m *customMessageModel) withSession(session sqlx.Session) MessageModel {
	return NewMessageModel(sqlx.NewSqlConnFromSession(session))
}
