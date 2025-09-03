use anyhow::Result;
use clap::{CommandFactory, Parser, Subcommand};
use crate::integrations;

#[derive(Parser)]
#[command(author, version, about, long_about = None)]
#[command(propagate_version = true)]
pub struct Cli {
    #[command(subcommand)]
    pub command: Commands,

    #[arg(short, long, action = clap::ArgAction::Count, global = true)]
    /// Increase verbosity level
    verbose: u8,

    #[arg(long, global = true)]
    /// Enable debug mode
    debug: bool,
}

#[derive(Subcommand)]
pub enum Commands {
    /// Analyze current working directory
    Analyze,
    /// Process a user prompt
    Prompt {
        /// The prompt to process
        prompt: Vec<String>,
    },
    /// Manage configuration
    Config,
    /// Generate a commit message for staged changes
    CommitMessage,
    /// Explain the code in a given file
    Explain {
        /// Path to the file to explain
        file_path: String,
    },
    /// Edit code in specified files or folders
    Edit {
        /// The paths to the files or directories to edit
        paths: Vec<String>,
        #[arg(long)]
        /// The user prompt describing the desired changes
        user_prompt: String,
    },
    /// Provide a code suggestion
    Suggestion {
        #[arg(long)]
        file: Option<String>,
        #[arg(long)]
        line: Option<usize>,
    },
    /// Manage integrations and models
    Integrations {
        #[command(subcommand)]
        command: IntegrationsCommands,
    },
    /// Stage, commit, and push all changes
    Push {
        #[arg(short, long)]
        /// The commit message
        message: Option<String>,
    },
}

#[derive(Subcommand)]
pub enum IntegrationsCommands {
    /// List all available integrations
    List,
    /// List available models for a provider
    ListModels { provider: String },
    /// Set the default model for a provider
    SetModel { provider: String, model: String },
}

use crate::config;
use crate::engine::Engine;

pub async fn run_cli() -> Result<()> {
    // Load configuration
    let config = config::load_config()?;

    // Create a new engine
    let engine = Engine::new(config);

    // Start with the static CLI structure derived from the structs
    let mut cmd = Cli::command();

    // Get the list of registered integrations
    let integration_list = integrations::list_integrations();

    // Dynamically add a subcommand for each integration
    for integration in &integration_list {
        cmd = cmd.subcommand(integration.command());
    }

    // Parse the command-line arguments against the fully built command
    let matches = cmd.get_matches();

    // You can handle global flags like this
    if matches.get_flag("debug") {
        println!("Debug mode enabled.");
    }

    // Handle the matched subcommand
    match matches.subcommand() {
        Some(("analyze", _)) => engine.run_analysis().await?,
        Some(("prompt", args)) => {
            let prompt_parts: Vec<String> = args.get_many("prompt").unwrap_or_default().cloned().collect();
            engine.process_prompt(&prompt_parts.join(" ")).await?;
        }
        Some(("config", _)) => println!("'config' command not implemented yet."),
        Some(("commit-message", _)) => {
            // In a real scenario, we'd get the git diff here.
            let diff = "fake git diff";
            engine.process_commit_message(diff).await?;
        }
        Some(("explain", args)) => {
            let file_path = args.get_one::<String>("file_path").unwrap();
            engine.explain_code(file_path).await?;
        }
        Some(("edit", args)) => {
            let paths: Vec<String> = args.get_many("paths").unwrap_or_default().cloned().collect();
            let user_prompt = args.get_one::<String>("user_prompt").unwrap();
            engine.edit_files(&paths, user_prompt, false).await?;
        }
        Some(("suggestion", _)) => println!("'suggestion' command not implemented yet."),
        Some(("integrations", args)) => match args.subcommand() {
            Some(("list", _)) => println!("'integrations list' not implemented yet."),
            Some(("list-models", sub_args)) => {
                let provider = sub_args.get_one::<String>("provider").unwrap();
                println!("'integrations list-models' for '{}' not implemented yet.", provider);
            }
            Some(("set-model", sub_args)) => {
                let provider = sub_args.get_one::<String>("provider").unwrap();
                let model = sub_args.get_one::<String>("model").unwrap();
                println!("'integrations set-model' for '{}' to '{}' not implemented yet.", provider, model);
            }
            _ => unreachable!("clap should have handled this"),
        },
        Some(("push", args)) => {
            let message = args.get_one::<String>("message");
            println!("'push' command with message '{:?}' not implemented yet.", message);
        }
        // Handle dynamic integration subcommands
        Some((name, _args)) => {
            if integration_list.iter().any(|i| i.metadata().name == name) {
                println!("Dynamic integration command '{}' called. Handler not implemented yet.", name);
            } else {
                println!("Unknown command '{}'", name);
            }
        }
        None => {
            unreachable!("clap should have handled no subcommand");
        }
    }

    Ok(())
}
