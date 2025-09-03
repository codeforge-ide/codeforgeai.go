use crate::config::Config;
use crate::models::Model;
use anyhow::Result;
use async_trait::async_trait;
use std::collections::HashMap;
use serde_json::Value;

// A dummy model for placeholder implementation
struct DummyModel;

#[async_trait]
impl Model for DummyModel {
    async fn send_request(&self, prompt: &str, _options: &HashMap<String, Value>) -> Result<String> {
        println!("--- Dummy Model ---");
        println!("Prompt: {}", prompt);
        println!("-------------------");
        Ok("This is a dummy response.".to_string())
    }
}

pub struct Engine {
    config: Config,
}

impl Engine {
    pub fn new(config: Config) -> Self {
        Self { config }
    }

    fn get_general_model(&self) -> Box<dyn Model> {
        // Placeholder: In the future, this will look at self.config
        // to decide which model implementation to return.
        Box::new(DummyModel)
    }

    fn get_code_model(&self) -> Box<dyn Model> {
        // Placeholder
        Box::new(DummyModel)
    }

    pub async fn run_analysis(&self) -> Result<()> {
        println!("Running analysis...");
        let model = self.get_general_model();
        let prompt = &self.config.directory_classification_prompt;

        // Placeholder for directory tree logic
        let tree_json = "{ \"name\": \".\", \"type\": \"directory\" }";
        let full_prompt = format!("{}\n{}", prompt, tree_json);

        let response = model.send_request(&full_prompt, &HashMap::new()).await?;
        println!("Analysis response: {}", response);

        // Placeholder for saving results
        println!("Analysis complete. Results would be saved to .codeforge.json");
        Ok(())
    }

    pub async fn process_prompt(&self, prompt: &str) -> Result<()> {
        println!("Processing prompt: {}", prompt);
        let general_model = self.get_general_model();
        let code_model = self.get_code_model();

        // Step 1: Finetune the prompt
        let finetune_prompt = format!("{}\n{}", self.config.prompt_finetune_prompt, prompt);
        let fine_tuned_prompt = general_model.send_request(&finetune_prompt, &HashMap::new()).await?;
        println!("Finetuned prompt: {}", fine_tuned_prompt);

        // Step 2: Determine if response should be code or command
        let code_or_command_prompt = format!("{}\n{}", self.config.code_or_command, fine_tuned_prompt);
        let response_type = general_model.send_request(&code_or_command_prompt, &HashMap::new()).await?;
        println!("Response type: {}", response_type);

        // Step 3: Process with appropriate model
        let final_response = if response_type.to_lowercase().contains("command") {
            let final_prompt = format!("{}\n{}", self.config.command_agent_prompt, fine_tuned_prompt);
            general_model.send_request(&final_prompt, &HashMap::new()).await?
        } else {
            let final_prompt = format!("{}\n{}", self.config.code_prompt, fine_tuned_prompt);
            code_model.send_request(&final_prompt, &HashMap::new()).await?
        };

        println!("Final response:\n{}", final_response);
        Ok(())
    }

    pub async fn explain_code(&self, file_path: &str) -> Result<()> {
        println!("Explaining code in file: {}", file_path);
        // Placeholder for reading file content
        let content = format!("// Content of {}", file_path);

        let model = self.get_code_model();
        let prompt = format!("{}\n\nFile: {}\n\n{}", self.config.explain_code_prompt, file_path, content);

        let response = model.send_request(&prompt, &HashMap::new()).await?;
        println!("Explanation:\n{}", response);
        Ok(())
    }

    pub async fn process_commit_message(&self, diff: &str) -> Result<()> {
        println!("Processing commit message for diff...");
        let code_model = self.get_code_model();
        let general_model = self.get_general_model();

        // Step 1: Generate commit message
        let commit_prompt = format!("{}\n{}", self.config.commit_message_prompt, diff);
        let commit_msg = code_model.send_request(&commit_prompt, &HashMap::new()).await?;
        println!("Generated commit message: {}", commit_msg);

        // Step 2: Generate gitmoji
        let gitmoji_prompt = format!("{}\n{}", self.config.gitmoji_prompt, commit_msg);
        let gitmoji = general_model.send_request(&gitmoji_prompt, &HashMap::new()).await?;
        println!("Generated gitmoji: {}", gitmoji);

        println!("Final commit message: {} {}", gitmoji.trim(), commit_msg.trim());
        Ok(())
    }

    pub async fn edit_files(&self, paths: &[String], user_prompt: &str, _allow_ignore: bool) -> Result<()> {
        println!("Editing files: {:?} with prompt: {}", paths, user_prompt);
        // This would involve iterating through files, reading them, and calling the model.
        // For now, it's just a placeholder.
        println!("Edit files functionality is not fully implemented yet.");
        Ok(())
    }
}
