package message

import (
	"context"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
)

const (
	CreateTableTemplate = `CREATE TABLE IF NOT EXISTS community.messages (
		session_id bigint,
		message_id bigint,
		send_id bigint,
		recipient_id bigint,
		reply_id bigint,
		create_time bigint,
		update_time bigint,
		status bigint,
		text text,
		message_type bigint,
		addition text,
		PRIMARY KEY (session_id, message_id)
	) WITH CLUSTERING ORDER BY (message_id DESC)`

	InsertMessage = `INSERT INTO community.messages (
		session_id, message_id, send_id, recipient_id, reply_id,
		create_time, update_time, status, text, message_type, addition
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	UpdateMessage = `UPDATE community.messages SET
		send_id = ?, recipient_id = ?, reply_id = ?,
		update_time = ?, status = ?, text = ?, message_type = ?, addition = ?
		WHERE session_id = ? AND message_id = ?`

	DeleteMessage = `DELETE FROM community.messages WHERE session_id = ? AND message_id = ?`

	GetMessageList = `SELECT session_id, message_id, send_id, recipient_id, reply_id,
		create_time, update_time, status, text, message_type, addition
		FROM community.messages WHERE session_id = ? AND message_id>? LIMIT ? order by create_time DESC`

	GetMessageByIds = `SELECT session_id, message_id, send_id, recipient_id, reply_id,
		create_time, update_time, status, text, message_type, addition
		FROM community.messages WHERE session_id = ? AND message_id IN ?`
)

type Message struct {
	MessageId   int64  `db:"message_id"`   // 消息唯一ID
	SessionId   int64  `db:"session_id"`   // 会话ID（关联群聊或单聊）
	SendId      int64  `db:"send_id"`      // 发送者ID
	RecipientId int64  `db:"recipient_id"` // 接收者ID（群聊时可设为0）
	ReplyId     int64  `db:"reply_id"`     // 回复的消息ID（0表示非回复）
	CreateTime  int64  `db:"create_time"`  // 创建时间
	UpdateTime  int64  `db:"update_time"`  // 更新时间
	Status      int64  `db:"status"`       // 消息状态
	Text        string `db:"text"`         // 文本内容
	MessageType int64  `db:"message_type"` // 消息类型
	Addition    string `db:"addition"`     // 附加消息（图片/文件等，JSON格式）
}

type MessageModel interface {
	CreateMessage(ctx context.Context, message *Message) error
	UpdateMessageById(ctx context.Context, sessionId, messageId int64, message *Message) error
	DeleteMessage(ctx context.Context, sessionId, messageId int64) error
	GetMessageList(ctx context.Context, sessionId int64, req, limit int) ([]*Message, error)
	GetMessageByIds(ctx context.Context, sessionId int64, messageIds []int64) ([]*Message, error)
}

type defaultMessageModel struct {
	session *gocql.Session
}

func NewMessageModel(session *gocql.Session) MessageModel {
	model := &defaultMessageModel{
		session: session,
	}

	if err := model.initTable(); err != nil {
		panic(fmt.Sprintf("Failed to initialize table: %v", err))
	}

	return model
}

func (m *defaultMessageModel) initTable() error {
	if err := m.session.Query(CreateTableTemplate).Exec(); err != nil {
		return errors.Wrap(err, "failed to create table")
	}
	return nil
}

// CreateMessage 创建消息
func (m *defaultMessageModel) CreateMessage(ctx context.Context, message *Message) error {
	if message == nil {
		return errors.New("message cannot be nil")
	}
	err := m.session.Query(InsertMessage,
		message.SessionId,
		message.MessageId,
		message.SendId,
		message.RecipientId,
		message.ReplyId,
		message.CreateTime,
		message.UpdateTime,
		message.Status,
		message.Text,
		message.MessageType,
		message.Addition,
	).Exec()

	if err != nil {
		return errors.Wrap(err, "failed to create message")
	}

	return nil
}

// UpdateMessageById 更新消息
func (m *defaultMessageModel) UpdateMessageById(ctx context.Context, sessionId, MessageId int64, message *Message) error {
	if message == nil {
		return errors.New("message cannot be nil")
	}

	err := m.session.Query(UpdateMessage,
		message.SendId,
		message.RecipientId,
		message.ReplyId,
		message.UpdateTime,
		message.Status,
		message.Text,
		message.MessageType,
		message.Addition,
		sessionId,
		MessageId,
	).WithContext(ctx).Exec()

	if err != nil {
		return errors.Wrap(err, "failed to update message")
	}

	return nil
}

// DeleteMessage 删除消息
func (m *defaultMessageModel) DeleteMessage(ctx context.Context, sessionId, messageId int64) error {
	err := m.session.Query(DeleteMessage, sessionId, messageId).Exec()
	if err != nil {
		return errors.Wrap(err, "failed to delete message")
	}

	return nil
}

// GetMessageList 获取消息列表（按proto中的GetMessageList接口）
func (m *defaultMessageModel) GetMessageList(ctx context.Context, sessionId int64, req, limit int) ([]*Message, error) {
	if limit <= 0 {
		limit = 50 // 默认限制
	}
	if limit > 1000 {
		limit = 1000 // 最大限制
	}

	iter := m.session.Query(GetMessageList, sessionId, req, limit).WithContext(ctx).Iter()

	var messages []*Message
	var message Message
	for iter.Scan(
		&message.SessionId,
		&message.MessageId,
		&message.SendId,
		&message.RecipientId,
		&message.ReplyId,
		&message.CreateTime,
		&message.UpdateTime,
		&message.Status,
		&message.Text,
		&message.MessageType,
		&message.Addition,
	) {
		msg := message // 复制结构体
		messages = append(messages, &msg)
	}

	if err := iter.Close(); err != nil {
		return nil, errors.Wrap(err, "failed to get message list")
	}

	return messages, nil
}

func (m *defaultMessageModel) GetMessageByIds(ctx context.Context, sessionId int64, messageIds []int64) ([]*Message, error) {
	if len(messageIds) == 0 {
		return []*Message{}, nil
	}
	iter := m.session.Query(GetMessageByIds, sessionId, messageIds).WithContext(ctx).Iter()

	var messages []*Message
	var message Message
	for iter.Scan(
		&message.SessionId,
		&message.MessageId,
		&message.SendId,
		&message.RecipientId,
		&message.ReplyId,
		&message.CreateTime,
		&message.UpdateTime,
		&message.Status,
		&message.Text,
		&message.MessageType,
		&message.Addition,
	) {
		msg := message
		messages = append(messages, &msg)
	}

	if err := iter.Close(); err != nil {
		return nil, errors.Wrap(err, "failed to get messages by ids")
	}

	return messages, nil
}

func (m *defaultMessageModel) Close() error {
	if m.session != nil {
		m.session.Close()
	}
	return nil
}
