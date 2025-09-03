pub mod api;
pub mod cli;
pub mod config;
pub mod engine;
pub mod integrations;
pub mod models;
pub mod utils;

pub async fn run() {
    if let Err(e) = cli::run_cli().await {
        eprintln!("Error: {}", e);
        std::process::exit(1);
    }
}
