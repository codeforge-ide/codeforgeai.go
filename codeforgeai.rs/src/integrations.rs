use clap::Command;
use once_cell::sync::Lazy;
use std::collections::HashMap;
use std::sync::{Arc, RwLock};

pub mod ollama;

// Metadata for an integration
pub struct IntegrationMetadata {
    pub name: &'static str,
    pub description: &'static str,
    pub capabilities: Vec<&'static str>,
}

// The core trait for any integration
pub trait Integration: Send + Sync {
    fn metadata(&self) -> IntegrationMetadata;
    fn command(&self) -> Command;
}

// The global registry for all integrations
type IntegrationRegistry = RwLock<HashMap<&'static str, Arc<dyn Integration>>>;

pub static REGISTRY: Lazy<IntegrationRegistry> = Lazy::new(|| {
    RwLock::new(HashMap::new())
});

/// Registers an integration in the global registry.
pub fn register_integration(integration: Arc<dyn Integration>) {
    let mut registry = REGISTRY.write().unwrap();
    let name = integration.metadata().name;
    if registry.contains_key(name) {
        // Handle duplicate registration if necessary
        eprintln!("Warning: Integration '{}' is already registered.", name);
    } else {
        registry.insert(name, integration);
    }
}

/// Lists all registered integrations.
pub fn list_integrations() -> Vec<Arc<dyn Integration>> {
    let registry = REGISTRY.read().unwrap();
    registry.values().cloned().collect()
}
