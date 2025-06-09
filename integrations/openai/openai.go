package openai

import "github.com/codeforge-ide/codeforgeai.go/models"

type OpenAIModel struct {
	// ...fields...
}

func (o *OpenAIModel) SendRequest(prompt string, config interface{}) (string, error) {
	// TODO: Implement OpenAI API call
	return "", nil
}
var _ models.Model = (*OpenAIModel)(nil)
