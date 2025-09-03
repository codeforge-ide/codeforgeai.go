use anyhow::Result;
use async_trait::async_trait;
use std::collections::HashMap;
use serde_json::Value;

#[async_trait]
pub trait Model: Send + Sync {
    async fn send_request(
        &self,
        prompt: &str,
        options: &HashMap<String, Value>,
    ) -> Result<String>;
}
