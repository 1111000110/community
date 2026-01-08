package model

import (
	"community/service/ai/rpc/ai"
)

func MongoAgentToRpcAgent(agent *Agent) *ai.Agent {
	if agent == nil {
		return nil
	}

	result := &ai.Agent{
		AgentId:    agent.AgentId,
		ApiKey:     agent.ApiKey,
		Name:       agent.Name,
		Desc:       agent.Desc,
		Icon:       agent.Icon,
		Status:     agent.Status,
		CreateTime: agent.CreateTime,
		UpdateTime: agent.UpdateTime,
		AgentChatConfig: &ai.AgentChatConfig{
			SystemPrompt:    agent.ChatConfig.SystemPrompt,
			IsStream:        agent.ChatConfig.IsStream,
			ChatType:        agent.ChatConfig.ChatType,
			ChatRound:       agent.ChatConfig.ChatRound,
			Temperature:     agent.ChatConfig.Temperature,
			MaxTokens:       agent.ChatConfig.MaxTokens,
			TopP:            agent.ChatConfig.TopP,
			EnableTools:     agent.ChatConfig.EnableTools,
			EnableFunctions: agent.ChatConfig.EnableFunctions,
			ToolConfig:      &ai.ToolConfig{},
			LlmId:           agent.ChatConfig.LlmId,
		},
	}

	// 处理 ToolConfig.KB
	if agent.ChatConfig.ToolConfig.KB != nil {
		result.AgentChatConfig.ToolConfig.Kb = &ai.KBSearchConfig{
			KnowledgeList:  agent.ChatConfig.ToolConfig.KB.KnowledgeList,
			TopK:           agent.ChatConfig.ToolConfig.KB.TopK,
			ScoreThreshold: agent.ChatConfig.ToolConfig.KB.ScoreThreshold,
		}
	}

	// 处理 ToolConfig.Web
	if agent.ChatConfig.ToolConfig.Web != nil {
		result.AgentChatConfig.ToolConfig.Web = &ai.WebSearchConfig{
			TopK:            agent.ChatConfig.ToolConfig.Web.TopK,
			RecencyDays:     agent.ChatConfig.ToolConfig.Web.RecencyDays,
			AllowDomains:    agent.ChatConfig.ToolConfig.Web.AllowDomains,
			BlockDomains:    agent.ChatConfig.ToolConfig.Web.BlockDomains,
			MaxCallsPerTurn: agent.ChatConfig.ToolConfig.Web.MaxCallsPerTurn,
		}
	}

	return result
}

func RpcAgentToMongoAgent(agent *ai.Agent) *Agent {
	if agent == nil {
		return nil
	}

	return &Agent{
		AgentId: agent.AgentId,
		ApiKey:  agent.ApiKey,
		Name:    agent.Name,
		Desc:    agent.Desc,
		Icon:    agent.Icon,
		Status:  agent.Status,
		ChatConfig: AgentChatConfig{
			SystemPrompt:    agent.GetAgentChatConfig().GetSystemPrompt(),
			IsStream:        agent.GetAgentChatConfig().GetIsStream(),
			ChatType:        agent.GetAgentChatConfig().GetChatType(),
			ChatRound:       agent.GetAgentChatConfig().GetChatRound(),
			Temperature:     agent.GetAgentChatConfig().GetTemperature(),
			MaxTokens:       agent.GetAgentChatConfig().GetMaxTokens(),
			TopP:            agent.GetAgentChatConfig().GetTopP(),
			EnableTools:     agent.GetAgentChatConfig().GetEnableTools(),
			EnableFunctions: agent.GetAgentChatConfig().GetEnableFunctions(),
			ToolConfig: ToolConfig{
				KB: &KBSearchConfig{
					KnowledgeList:  agent.GetAgentChatConfig().GetToolConfig().GetKb().GetKnowledgeList(),
					TopK:           agent.GetAgentChatConfig().GetToolConfig().GetKb().GetTopK(),
					ScoreThreshold: agent.GetAgentChatConfig().GetToolConfig().GetKb().GetScoreThreshold(),
				},
				Web: &WebSearchConfig{
					TopK:            agent.GetAgentChatConfig().GetToolConfig().GetWeb().GetTopK(),
					RecencyDays:     agent.GetAgentChatConfig().GetToolConfig().GetWeb().GetRecencyDays(),
					AllowDomains:    agent.GetAgentChatConfig().GetToolConfig().GetWeb().GetAllowDomains(),
					BlockDomains:    agent.GetAgentChatConfig().GetToolConfig().GetWeb().GetBlockDomains(),
					MaxCallsPerTurn: agent.GetAgentChatConfig().GetToolConfig().GetWeb().GetMaxCallsPerTurn(),
				},
			},
			LlmId: agent.GetAgentChatConfig().GetLlmId(),
		},
	}
}

func MongoAgentsToRpcAgents(agents []*Agent) []*ai.Agent {
	if agents == nil {
		return nil
	}
	result := make([]*ai.Agent, 0, len(agents))
	for _, agent := range agents {
		if agent != nil {
			result = append(result, MongoAgentToRpcAgent(agent))
		}
	}
	return result
}

func RpcAgentsToMongoAgents(agents []*ai.Agent) []*Agent {
	if agents == nil {
		return nil
	}
	result := make([]*Agent, 0, len(agents))
	for _, agent := range agents {
		if agent != nil {
			result = append(result, RpcAgentToMongoAgent(agent))
		}
	}
	return result
}
