package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type IntegrationsConfig struct {
	Ollama        IntegrationEntry `json:"ollama"`
	GithubModels  IntegrationEntry `json:"githubmodels"`
	OpenAPI       IntegrationEntry `json:"openapi"`
	GithubCopilot IntegrationEntry `json:"githubcopilot"`
	Default       string           `json:"default"`
}

type IntegrationEntry struct {
	Enabled bool `json:"enabled"`
}

type Config struct {
	GeneralModel                  string             `json:"general_model"`
	GeneralPrompt                 string             `json:"general_prompt"`
	CodeModel                     string             `json:"code_model"`
	CodePrompt                    string             `json:"code_prompt"`
	GeneralModelGithub            string             `json:"general_model_github"`
	CodeModelGithub               string             `json:"code_model_github"`
	DirectoryClassificationPrompt string             `json:"directory_classification_prompt"`
	Debug                         bool               `json:"debug"`
	FormatLineSeparator           int                `json:"format_line_separator"`
	GitmojiPrompt                 string             `json:"gitmoji_prompt"`
	CommitMessagePrompt           string             `json:"commit_message_prompt"`
	EditFinetunePrompt            string             `json:"edit_finetune_prompt"`
	CodeOrCommand                 string             `json:"code_or_command"`
	CommandAgentPrompt            string             `json:"command_agent_prompt"`
	PromptFinetunePrompt          string             `json:"prompt_finetune_prompt"`
	LanguageClassificationPrompt  string             `json:"language_classification_prompt"`
	ReadmeSummaryPrompt           string             `json:"readme_summary_prompt"`
	SpecificFileClassification    string             `json:"specific_file_classification"`
	ImproveCodePrompt             string             `json:"improve_code_prompt"`
	ExplainCodePrompt             string             `json:"explain_code_prompt"`
	SuggestionPrompt              string             `json:"suggestion_prompt"`
	ExtractCodeBlocksPrompt       string             `json:"extract_code_blocks_prompt"`
	FormatCodePrompt              string             `json:"format_code_prompt"`
	Integrations                  IntegrationsConfig `json:"integrations"`
	GithubModelsList              string             `json:"github_models_list"`
	// Optionally add GithubToken string `json:"github_token"` to Config struct if you want to support it from config.
}

func DefaultConfig() Config {
	return Config{
		GeneralModel:                  "gemma3:1b",
		GeneralPrompt:                 "based on the below prompt and without returning anything else, restructure it so that it is strictly understandable to a coding ai agent with json output for file changes:",
		CodeModel:                     "qwen2.5-coder:1.5b",
		CodePrompt:                    "in very clear, concise manner, solve the below request:",
		GeneralModelGithub:            "gpt-4o-mini",
		CodeModelGithub:               "gpt-4o-mini",
		DirectoryClassificationPrompt: "Given the complete tree structure below as valid JSON, recursively process every single file and directory (based on its relative path) that is present. For each node, assign exactly one classification: 'useful' for files and directories that developers interact with, 'useless' for build, template, or temporary files and directories, and 'source' for source control or related files. For every node, return an object with the keys: 'type' (either 'file' or 'directory'), 'name', 'contents' (an array of child entries for directories, or file details for files), and a new key 'classification' that holds one of 'useful', 'useless', or 'source'. Ensure every file and directory from the input is included exactly once with one classification. Return only valid JSON with this structure and nothing else.",
		Debug:                         false,
		FormatLineSeparator:           5,
		GitmojiPrompt:                 "reply only with a single emoji character that best fits the below commit message, and nothing else.",
		CommitMessagePrompt:           "Generate a very short and very concise, one sentence commit message for these code changes, and nothng else. ",
		EditFinetunePrompt:            "edit this code according to the below prompt and return nothing but the edited code",
		CodeOrCommand:                 "reply with either code or command only; is the below request best satisfied with a code response or command response:",
		CommandAgentPrompt:            "one for each line and nothing else, return a list of commands that can be executed to achieve the below request, and nothing else:",
		PromptFinetunePrompt:          "in a clear and concise manner, rephrase the following prompt to be more understandable to a coding ai agent, return the rephrased prompt and nothing else",
		LanguageClassificationPrompt:  "in one word only, what programming language is used in this project tree structure",
		ReadmeSummaryPrompt:           "in one short sentence only, generate a concise summary of this text below, and nothing else",
		SpecificFileClassification:    "taking the path and content of this file and classify it into either only user code file or project code file or source control file",
		ImproveCodePrompt:             "given this block of code, improve the code generally and return nothing but the improved code:",
		ExplainCodePrompt:             "explain the following code in a clear and concise manner",
		SuggestionPrompt:              "provide a helpful code suggestion for the following code context:",
		ExtractCodeBlocksPrompt:       "extract all code blocks from the following text and return them in a structured format:",
		FormatCodePrompt:              "format the following code for better readability while preserving functionality:",
		Integrations: IntegrationsConfig{
			Ollama:        IntegrationEntry{Enabled: true},
			GithubModels:  IntegrationEntry{Enabled: false},
			OpenAPI:       IntegrationEntry{Enabled: false},
			GithubCopilot: IntegrationEntry{Enabled: false},
			Default:       "ollama",
		},
		GithubModelsList: "",
	}
}

func configFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "./.codeforgeai.json"
	}
	return filepath.Join(home, ".codeforgeai.json")
}

func DataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "./.codeforgeai"
	}
	dir := filepath.Join(home, ".codeforgeai")
	os.MkdirAll(dir, 0700)
	return dir
}

func LoadConfig(path string) (Config, error) {
	if path == "" {
		path = configFilePath()
	}
	f, err := os.Open(path)
	if err != nil {
		cfg := DefaultConfig()
		SaveConfig(path, cfg)
		return cfg, nil
	}
	defer f.Close()
	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return DefaultConfig(), err
	}
	return cfg, nil
}

func SaveConfig(path string, cfg Config) error {
	if path == "" {
		path = configFilePath()
	}	return json.NewEncoder(f).Encode(cfg)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errgPrompts(path string) (Config, error) {
	}Config(path)
	f, err := os.Create(path)f err != nil {
	if err != nil {
		return err
	}	changed := false
	defer f.Close()

func EnsureConfigPrompts(path string) (Config, error) {ery field in Config)
	cfg, err := LoadConfig(path)odel == "" {
	if err != nil {cfg.GeneralModel = def.GeneralModel
		return cfg, err
	}
	changed := falserompt == "" {
	def := DefaultConfig()cfg.GeneralPrompt = def.GeneralPrompt

	// Ensure all fields are set (for every field in Config)
	if cfg.GeneralModel == "" {l == "" {
		cfg.GeneralModel = def.GeneralModelcfg.CodeModel = def.CodeModel
		changed = true
	}
	if cfg.GeneralPrompt == "" {pt == "" {
		cfg.GeneralPrompt = def.GeneralPromptcfg.CodePrompt = def.CodePrompt
		changed = true
	}
	if cfg.CodeModel == "" {odelGithub == "" {
		cfg.CodeModel = def.CodeModelcfg.GeneralModelGithub = def.GeneralModelGithub
		changed = true
	}
	if cfg.CodePrompt == "" {lGithub == "" {
		cfg.CodePrompt = def.CodePromptcfg.CodeModelGithub = def.CodeModelGithub
		changed = true
	}
	if cfg.GeneralModelGithub == "" {yClassificationPrompt == "" {
		cfg.GeneralModelGithub = def.GeneralModelGithubcfg.DirectoryClassificationPrompt = def.DirectoryClassificationPrompt
		changed = true
	}
	if cfg.CodeModelGithub == "" {neSeparator == 0 {
		cfg.CodeModelGithub = def.CodeModelGithubcfg.FormatLineSeparator = def.FormatLineSeparator
		changed = true
	}
	if cfg.DirectoryClassificationPrompt == "" {rompt == "" {
		cfg.DirectoryClassificationPrompt = def.DirectoryClassificationPromptcfg.GitmojiPrompt = def.GitmojiPrompt
		changed = true
	}
	if cfg.FormatLineSeparator == 0 {ssagePrompt == "" {
		cfg.FormatLineSeparator = def.FormatLineSeparatorcfg.CommitMessagePrompt = def.CommitMessagePrompt
		changed = true
	}
	if cfg.GitmojiPrompt == "" {tunePrompt == "" {
		cfg.GitmojiPrompt = def.GitmojiPromptcfg.EditFinetunePrompt = def.EditFinetunePrompt
		changed = true
	}
	if cfg.CommitMessagePrompt == "" {mmand == "" {
		cfg.CommitMessagePrompt = def.CommitMessagePromptcfg.CodeOrCommand = def.CodeOrCommand
		changed = true
	}
	if cfg.EditFinetunePrompt == "" {gentPrompt == "" {
		cfg.EditFinetunePrompt = def.EditFinetunePromptcfg.CommandAgentPrompt = def.CommandAgentPrompt
		changed = true
	}
	if cfg.CodeOrCommand == "" {netunePrompt == "" {
		cfg.CodeOrCommand = def.CodeOrCommandcfg.PromptFinetunePrompt = def.PromptFinetunePrompt
		changed = true
	}
	if cfg.CommandAgentPrompt == "" {ClassificationPrompt == "" {
		cfg.CommandAgentPrompt = def.CommandAgentPromptcfg.LanguageClassificationPrompt = def.LanguageClassificationPrompt
		changed = true
	}
	if cfg.PromptFinetunePrompt == "" {mmaryPrompt == "" {
		cfg.PromptFinetunePrompt = def.PromptFinetunePromptcfg.ReadmeSummaryPrompt = def.ReadmeSummaryPrompt
		changed = true
	}
	if cfg.LanguageClassificationPrompt == "" {FileClassification == "" {
		cfg.LanguageClassificationPrompt = def.LanguageClassificationPromptcfg.SpecificFileClassification = def.SpecificFileClassification
		changed = true
	}
	if cfg.ReadmeSummaryPrompt == "" {odePrompt == "" {
		cfg.ReadmeSummaryPrompt = def.ReadmeSummaryPromptcfg.ImproveCodePrompt = def.ImproveCodePrompt
		changed = true
	}
	if cfg.SpecificFileClassification == "" {odePrompt == "" {
		cfg.SpecificFileClassification = def.SpecificFileClassificationcfg.ExplainCodePrompt = def.ExplainCodePrompt
		changed = true
	}
	if cfg.ImproveCodePrompt == "" {onPrompt == "" {
		cfg.ImproveCodePrompt = def.ImproveCodePromptcfg.SuggestionPrompt = def.SuggestionPrompt
		changed = true
	}
	if cfg.ExplainCodePrompt == "" {odeBlocksPrompt == "" {
		cfg.ExplainCodePrompt = def.ExplainCodePromptcfg.ExtractCodeBlocksPrompt = def.ExtractCodeBlocksPrompt
		changed = true
	}
	if cfg.SuggestionPrompt == "" {dePrompt == "" {
		cfg.SuggestionPrompt = def.SuggestionPromptcfg.FormatCodePrompt = def.FormatCodePrompt
		changed = true
	}
	if cfg.ExtractCodeBlocksPrompt == "" {sent and complete
		cfg.ExtractCodeBlocksPrompt = def.ExtractCodeBlocksPromptions.Default == "" {
		changed = truecfg.Integrations = def.Integrations
	}
	if cfg.FormatCodePrompt == "" {
		cfg.FormatCodePrompt = def.FormatCodePrompt
		changed = truetions.Ollama == IntegrationEntry{}) {
	}cfg.Integrations.Ollama = def.Integrations.Ollama
	// Ensure integrations config is present and complete
	if cfg.Integrations.Default == "" {
		cfg.Integrations = def.Integrationstions.GithubModels == IntegrationEntry{}) {
		changed = truecfg.Integrations.GithubModels = def.Integrations.GithubModels
	}
	// Ensure all sub-integrations are present
	if (cfg.Integrations.Ollama == IntegrationEntry{}) {tions.OpenAPI == IntegrationEntry{}) {
		cfg.Integrations.Ollama = def.Integrations.Ollamacfg.Integrations.OpenAPI = def.Integrations.OpenAPI
		changed = true
	}
	if (cfg.Integrations.GithubModels == IntegrationEntry{}) {tions.GithubCopilot == IntegrationEntry{}) {
		cfg.Integrations.GithubModels = def.Integrations.GithubModelscfg.Integrations.GithubCopilot = def.Integrations.GithubCopilot
		changed = true
	}	}
	if (cfg.Integrations.OpenAPI == IntegrationEntry{}) {bool, so no need to check for empty string
		cfg.Integrations.OpenAPI = def.Integrations.OpenAPI
		changed = truef cfg.GithubModelsList == "" {
	}lsList = def.GithubModelsList
	if (cfg.Integrations.GithubCopilot == IntegrationEntry{}) {	changed = true
		cfg.Integrations.GithubCopilot = def.Integrations.GithubCopilot	}
		changed = true
	}
	// Debug is bool, so no need to check for empty string
}
	if changed {	return cfg, nil










}	fmt.Println(string(b))	b, _ := json.MarshalIndent(cfg, "", "  ")func PrintConfig(cfg Config) {}	return cfg, nil	}		SaveConfig(path, cfg)}

func PrintConfig(cfg Config) {
	b, _ := json.MarshalIndent(cfg, "", "  ")
	fmt.Println(string(b))
}
