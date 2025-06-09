package config

type Config struct {
	GeneralModel string `json:"general_model"`
	CodeModel    string `json:"code_model"`
	// ...other fields...
}

func LoadConfig(path string) (*Config, error) {
	// TODO: Load from JSON file, create default if missing
	return &Config{}, nil
}
