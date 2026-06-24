package config

import "os"

type Config struct {
	Port         string
	SupabaseURL  string
	SupabaseKey  string
	GeminiKey    string
	GroqKey      string
	OpenRouterKey string
	HuggingFaceKey string
}

func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "8080"),
		SupabaseURL:   getEnv("SUPABASE_URL", ""),
		SupabaseKey:   getEnv("SUPABASE_KEY", ""),
		GeminiKey:     getEnv("GEMINI_KEY", ""),
		GroqKey:       getEnv("GROQ_KEY", ""),
		OpenRouterKey: getEnv("OPENROUTER_KEY", ""),
		HuggingFaceKey: getEnv("HUGGINGFACE_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
