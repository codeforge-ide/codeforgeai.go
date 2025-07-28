package openai

import "github.com/codeforge-ide/codeforgeai.go/models"

import "github.com/codeforge-ide/codeforgeai.go/modeliface"

type OpenAIModel struct {
	modeliface.BaseAgent
	// ...fields...
}

func (o *OpenAIModel) SendRequest(prompt string, config interface{}) (string, error) {
	// TODO: Implement OpenAI API call
	return "", nil
}

var _ models.Model = (*OpenAIModel)(nil)
var _ modeliface.Agent = (*OpenAIModel)(nil)

// NewOpenAIModel creates a new OpenAIModel with the given name.
func NewOpenAIModel(name string) *OpenAIModel {
	return &OpenAIModel{
		BaseAgent: modeliface.BaseAgent{NameStr: name},
	}
}
