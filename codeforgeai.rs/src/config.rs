use anyhow::Result;
use serde::{Deserialize, Serialize};
use std::path::PathBuf;

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct IntegrationEntry {
    pub enabled: bool,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct IntegrationsConfig {
    pub ollama: IntegrationEntry,
    pub githubmodels: IntegrationEntry,
    pub openapi: IntegrationEntry,
    pub githubcopilot: IntegrationEntry,
    pub io: IntegrationEntry,
    pub default: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
#[serde(rename_all = "camelCase")]
pub struct Config {
    #[serde(default = "default_general_model")]
    pub general_model: String,
    #[serde(default = "default_general_prompt")]
    pub general_prompt: String,
    #[serde(default = "default_code_model")]
    pub code_model: String,
    #[serde(default = "default_code_prompt")]
    pub code_prompt: String,
    #[serde(default = "default_general_model_github")]
    pub general_model_github: String,
    #[serde(default = "default_code_model_github")]
    pub code_model_github: String,
    #[serde(default = "default_directory_classification_prompt")]
    pub directory_classification_prompt: String,
    #[serde(default)]
    pub debug: bool,
    #[serde(default = "default_format_line_separator")]
    pub format_line_separator: i64,
    #[serde(default = "default_gitmoji_prompt")]
    pub gitmoji_prompt: String,
    #[serde(default = "default_commit_message_prompt")]
    pub commit_message_prompt: String,
    #[serde(default = "default_edit_finetune_prompt")]
    pub edit_finetune_prompt: String,
    #[serde(default = "default_code_or_command")]
    pub code_or_command: String,
    #[serde(default = "default_command_agent_prompt")]
    pub command_agent_prompt: String,
    #[serde(default = "default_prompt_finetune_prompt")]
    pub prompt_finetune_prompt: String,
    #[serde(default = "default_language_classification_prompt")]
    pub language_classification_prompt: String,
    #[serde(default = "default_readme_summary_prompt")]
    pub readme_summary_prompt: String,
    #[serde(default = "default_specific_file_classification")]
    pub specific_file_classification: String,
    #[serde(default = "default_improve_code_prompt")]
    pub improve_code_prompt: String,
    #[serde(default = "default_explain_code_prompt")]
    pub explain_code_prompt: String,
    #[serde(default = "default_suggestion_prompt")]
    pub suggestion_prompt: String,
    #[serde(default = "default_extract_code_blocks_prompt")]
    pub extract_code_blocks_prompt: String,
    #[serde(default = "default_format_code_prompt")]
    pub format_code_prompt: String,

    // Integration config fields
    #[serde(default = "default_ollama_model")]
    pub ollama_model: String,
    #[serde(default = "default_ollama_endpoint")]
    pub ollama_endpoint: String,
    #[serde(default = "default_github_models_model")]
    pub github_models_model: String,
    #[serde(default = "default_github_models_endpoint")]
    pub github_models_endpoint: String,
    #[serde(default)]
    pub github_token: String,
    #[serde(default)]
    pub io_net_api_key: String,
    #[serde(default)]
    pub integrations: IntegrationsConfig,
    #[serde(default)]
    pub github_models_list: String,
    #[serde(default)]
    pub copilot_token: String,
    #[serde(default)]
    pub openai_api_key: String,
}

impl Default for Config {
    fn default() -> Self {
        Self {
            general_model: default_general_model(),
            general_prompt: default_general_prompt(),
            code_model: default_code_model(),
            code_prompt: default_code_prompt(),
            general_model_github: default_general_model_github(),
            code_model_github: default_code_model_github(),
            directory_classification_prompt: default_directory_classification_prompt(),
            debug: false,
            format_line_separator: default_format_line_separator(),
            gitmoji_prompt: default_gitmoji_prompt(),
            commit_message_prompt: default_commit_message_prompt(),
            edit_finetune_prompt: default_edit_finetune_prompt(),
            code_or_command: default_code_or_command(),
            command_agent_prompt: default_command_agent_prompt(),
            prompt_finetune_prompt: default_prompt_finetune_prompt(),
            language_classification_prompt: default_language_classification_prompt(),
            readme_summary_prompt: default_readme_summary_prompt(),
            specific_file_classification: default_specific_file_classification(),
            improve_code_prompt: default_improve_code_prompt(),
            explain_code_prompt: default_explain_code_prompt(),
            suggestion_prompt: default_suggestion_prompt(),
            extract_code_blocks_prompt: default_extract_code_blocks_prompt(),
            format_code_prompt: default_format_code_prompt(),
            ollama_model: default_ollama_model(),
            ollama_endpoint: default_ollama_endpoint(),
            github_models_model: default_github_models_model(),
            github_models_endpoint: default_github_models_endpoint(),
            github_token: String::new(),
            io_net_api_key: String::new(),
            integrations: Default::default(),
            github_models_list: String::new(),
            copilot_token: String::new(),
            openai_api_key: String::new(),
        }
    }
}

impl Default for IntegrationsConfig {
    fn default() -> Self {
        Self {
            ollama: IntegrationEntry { enabled: true },
            githubmodels: IntegrationEntry { enabled: false },
            openapi: IntegrationEntry { enabled: false },
            githubcopilot: IntegrationEntry { enabled: false },
            io: IntegrationEntry { enabled: false },
            default: "ollama".to_string(),
        }
    }
}

pub fn load_config() -> Result<Config> {
    let config_path = config_file_path()?;

    let builder = config::Config::builder()
        .add_source(config::File::from(config_path).required(false))
        .add_source(config::Environment::with_prefix("CODEFORGEAI").separator("__"));

    let settings: Config = builder.build()?.try_deserialize()?;
    Ok(settings)
}

fn config_file_path() -> Result<PathBuf> {
    let home_dir = dirs::home_dir().ok_or_else(|| anyhow::anyhow!("Could not find home directory"))?;
    Ok(home_dir.join(".codeforgeai.json"))
}

// Default value functions
fn default_general_model() -> String { "gemma3:1b".to_string() }
fn default_general_prompt() -> String { "based on the below prompt and without returning anything else, restructure it so that it is strictly understandable to a coding ai agent with json output for file changes:".to_string() }
fn default_code_model() -> String { "qwen2.5-coder:1.5b".to_string() }
fn default_code_prompt() -> String { "in very clear, concise manner, solve the below request:".to_string() }
fn default_general_model_github() -> String { "gpt-4o-mini".to_string() }
fn default_code_model_github() -> String { "gpt-4o-mini".to_string() }
fn default_directory_classification_prompt() -> String { "Given the complete tree structure below as valid JSON, recursively process every single file and directory (based on its relative path) that is present. For each node, assign exactly one classification: 'useful' for files and directories that developers interact with, 'useless' for build, template, or temporary files and directories, and 'source' for source control or related files. For every node, return an object with the keys: 'type' (either 'file' or 'directory'), 'name', 'contents' (an array of child entries for directories, or file details for files), and a new key 'classification' that holds one of 'useful', 'useless', or 'source'. Ensure every file and directory from the input is included exactly once with one classification. Return only valid JSON with this structure and nothing else.".to_string() }
fn default_format_line_separator() -> i64 { 5 }
fn default_gitmoji_prompt() -> String { "reply only with a single emoji character that best fits the below commit message, and nothing else.".to_string() }
fn default_commit_message_prompt() -> String { "Generate a very short and very concise, one sentence commit message for these code changes, and nothng else. ".to_string() }
fn default_edit_finetune_prompt() -> String { "edit this code according to the below prompt and return nothing but the edited code".to_string() }
fn default_code_or_command() -> String { "reply with either code or command only; is the below request best satisfied with a code response or command response:".to_string() }
fn default_command_agent_prompt() -> String { "one for each line and nothing else, return a list of commands that can be executed to achieve the below request, and nothing else:".to_string() }
fn default_prompt_finetune_prompt() -> String { "in a clear and concise manner, rephrase the following prompt to be more understandable to a coding ai agent, return the rephrased prompt and nothing else".to_string() }
fn default_language_classification_prompt() -> String { "in one word only, what programming language is used in this project tree structure".to_string() }
fn default_readme_summary_prompt() -> String { "in one short sentence only, generate a concise summary of this text below, and nothing else".to_string() }
fn default_specific_file_classification() -> String { "taking the path and content of this file and classify it into either only user code file or project code file or source control file".to_string() }
fn default_improve_code_prompt() -> String { "given this block of code, improve the code generally and return nothing but the improved code:".to_string() }
fn default_explain_code_prompt() -> String { "explain the following code in a clear and concise manner".to_string() }
fn default_suggestion_prompt() -> String { "provide a helpful code suggestion for the following code context:".to_string() }
fn default_extract_code_blocks_prompt() -> String { "extract all code blocks from the following text and return them in a structured format:".to_string() }
fn default_format_code_prompt() -> String { "format the following code for better readability while preserving functionality:".to_string() }
fn default_ollama_model() -> String { "qwen2.5-coder:1.5b".to_string() }
fn default_ollama_endpoint() -> String { "http://localhost:11434/api/generate".to_string() }
fn default_github_models_model() -> String { "gpt-4o-mini".to_string() }
fn default_github_models_endpoint() -> String { "https://models.inference.ai.azure.com/chat/completions".to_string() }
