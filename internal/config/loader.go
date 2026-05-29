package config

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// envVarPattern matches ${VAR_NAME} syntax in configuration values
var envVarPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

// Load loads configuration from the specified file path.
// It parses JSON and interpolates environment variables.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %q: %w", path, err)
	}

	return loadFromJSONBytes(data, "file")
}

// LoadFromJSON loads configuration from raw JSON content.
// It parses JSON and interpolates environment variables.
func LoadFromJSON(raw string) (*Config, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, fmt.Errorf("configuration JSON is empty")
	}

	return loadFromJSONBytes([]byte(raw), "environment variable")
}

func loadFromJSONBytes(data []byte, source string) (*Config, error) {
	if source == "" {
		source = "configuration"
	}

	// First unmarshal to get raw structure
	var rawConfig map[string]interface{}
	if err := json.Unmarshal(data, &rawConfig); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON from %s: %w", source, err)
	}

	// Interpolate environment variables
	interpolated := interpolateEnvVars(rawConfig)

	// Marshal back to JSON and unmarshal into typed config
	interpolatedJSON, err := json.Marshal(interpolated)
	if err != nil {
		return nil, fmt.Errorf("failed to re-marshal config from %s: %w", source, err)
	}

	var cfg Config
	if err := json.Unmarshal(interpolatedJSON, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config from %s into struct: %w", source, err)
	}

	return &cfg, nil
}

// interpolateEnvVars recursively replaces ${VAR_NAME} with environment variable values.
func interpolateEnvVars(data interface{}) interface{} {
	switch v := data.(type) {
	case string:
		return envVarPattern.ReplaceAllStringFunc(v, func(match string) string {
			// Extract variable name from ${VAR_NAME}
			varName := envVarPattern.FindStringSubmatch(match)[1]
			if value, exists := os.LookupEnv(varName); exists {
				return value
			}
			// Keep original if environment variable doesn't exist
			return match
		})

	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, val := range v {
			result[key] = interpolateEnvVars(val)
		}
		return result

	case []interface{}:
		result := make([]interface{}, len(v))
		for i, val := range v {
			result[i] = interpolateEnvVars(val)
		}
		return result

	default:
		return v
	}
}

// LoadFromSource loads configuration with priority: flag > env > default.
// - configPath: path from CLI flag (highest priority)
// - envVarName: environment variable name to check for path
// - defaultPath: fallback path if neither flag nor env are set
func LoadFromSource(configPath, envVarName, defaultPath string) (*Config, error) {
	// Determine which path to use based on priority
	path := configPath
	if path == "" {
		path = os.Getenv(envVarName)
	}
	if path == "" {
		path = defaultPath
	}

	if path == "" {
		return nil, fmt.Errorf("no configuration file specified")
	}

	return Load(path)
}

// LoadFromJSONOrSource loads configuration with priority:
// jsonEnvVar > flag > path env var > default file path.
// - configPath: path from CLI flag
// - pathEnvVarName: environment variable name to check for file path
// - jsonEnvVarName: environment variable name to check for raw JSON config
// - defaultPath: fallback path if no other source is set
func LoadFromJSONOrSource(configPath, pathEnvVarName, jsonEnvVarName, defaultPath string) (*Config, error) {
	if jsonEnvVarName != "" {
		if rawJSON := strings.TrimSpace(os.Getenv(jsonEnvVarName)); rawJSON != "" {
			return LoadFromJSON(rawJSON)
		}
	}

	return LoadFromSource(configPath, pathEnvVarName, defaultPath)
}
