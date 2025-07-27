package engine

import (
	"errors"
	"github.com/codeforge-ide/codeforgeai.go/modeliface"
	"sync"
)

type AgentManager struct {
	agents  map[string]modeliface.Agent
	current string
	mu      sync.RWMutex
}

func NewAgentManager() *AgentManager {
	return &AgentManager{agents: make(map[string]modeliface.Agent)}
}

func (am *AgentManager) Register(agent modeliface.Agent) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.agents[agent.Name()] = agent
}

func (am *AgentManager) Switch(name string) error {
	am.mu.Lock()
	defer am.mu.Unlock()
	if _, ok := am.agents[name]; !ok {
		return errors.New("agent not found: " + name)
	}
	am.current = name
	return nil
}

func (am *AgentManager) Current() modeliface.Agent {
	am.mu.RLock()
	defer am.mu.RUnlock()
	return am.agents[am.current]
}

func (am *AgentManager) List() []string {
	am.mu.RLock()
	defer am.mu.RUnlock()
	names := make([]string, 0, len(am.agents))
	for name := range am.agents {
		names = append(names, name)
	}
	return names
}
