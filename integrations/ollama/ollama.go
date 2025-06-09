package ollama

import "github.com/codeforge-ide/codeforgeai.go/models"

type OllamaModel struct {
	// ...fields...
}

func (o *OllamaModel) SendRequest(prompt string, config interface{}) (string, error) {
	// TODO: Implement Ollama API call
	return "", nil
}
var _ models.Model = (*OllamaModel)(nil)
