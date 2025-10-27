package message

import "community/service/message/rpc/message"

func RpcModelToModel(messageInfo *message.MessageDetail) *Message {
	return &Message{
		MessageId:   messageInfo.GetMessageId(),
		SessionId:   messageInfo.GetSessionId(),
		SendId:      messageInfo.GetSendId(),
		ReplyId:     messageInfo.GetReplyId(),
		CreateTime:  messageInfo.GetCreateTime(),
		UpdateTime:  messageInfo.GetUpdateTime(),
		Status:      messageInfo.GetStatus(),
		Text:        messageInfo.GetContent().GetText(),
		MessageType: messageInfo.GetContent().GetMessageType(),
		Addition:    messageInfo.GetContent().GetAddition(),
	}
}

func ModelToRpcModel(messageInfo *Message) *message.MessageDetail {
	return &message.MessageDetail{
		MessageId:  messageInfo.MessageId,
		SessionId:  messageInfo.SessionId,
		SendId:     messageInfo.SendId,
		ReplyId:    messageInfo.ReplyId,
		CreateTime: messageInfo.CreateTime,
		UpdateTime: messageInfo.UpdateTime,
		Status:     messageInfo.Status,
		Content: &message.MessageContent{
			Text:        messageInfo.Text,
			MessageType: messageInfo.MessageType,
			Addition:    messageInfo.Addition,
		},
	}
}

func ModelsToRpcModels(messageInfoList []*Message) []*message.MessageDetail {
	result := make([]*message.MessageDetail, len(messageInfoList))
	for i, messageInfo := range messageInfoList {
		result[i] = ModelToRpcModel(messageInfo)
	}
	return result
}
