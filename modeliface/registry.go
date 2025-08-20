package modeliface

import (
	"github.com/spf13/cobra"
	"sync"
)

// IntegrationMetadata holds descriptive and config info for an integration.
type IntegrationMetadata struct {
	Name         string
	Description  string
	ConfigKeys   []string
	Secrets      []string
	Commands     []string // e.g. "prompt", "list-models"
	Capabilities []string // e.g. "text_completion", "defi_price"
}

// IntegrationRegistration holds metadata and a CLI command factory.
type IntegrationRegistration struct {
	Metadata       IntegrationMetadata
	CommandFactory func() *cobra.Command // returns the root command for this integration
}

// AgentRegistry manages all registered integrations.
type AgentRegistry struct {
	integrations map[string]IntegrationRegistration
	mu           sync.RWMutex
}

// GlobalAgentRegistry is the shared registry instance.
var GlobalAgentRegistry = &AgentRegistry{
	integrations: make(map[string]IntegrationRegistration),
}

// RegisterIntegration adds a new integration to the registry.
func (r *AgentRegistry) RegisterIntegration(reg IntegrationRegistration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.integrations[reg.Metadata.Name] = reg
}

// GetIntegration retrieves a registration by name.
func (r *AgentRegistry) GetIntegration(name string) (IntegrationRegistration, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	reg, ok := r.integrations[name]
	return reg, ok
}

// ListIntegrations returns all registered integrations.
func (r *AgentRegistry) ListIntegrations() []IntegrationRegistration {
	r.mu.RLock()
	defer r.mu.RUnlock()
	regs := make([]IntegrationRegistration, 0, len(r.integrations))
	for _, reg := range r.integrations {
		regs = append(regs, reg)
	}
	return regs
}
