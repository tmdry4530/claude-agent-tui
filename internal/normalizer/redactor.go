package normalizer

import (
	"encoding/json"
	"regexp"
	"strings"
)

const (
	maxRedactDepth = 10
	redactedValue  = "***REDACTED***"
)

// Redactor masks sensitive information in payloads.
type Redactor struct {
	sensitiveKeys   map[string]bool
	sensitivePatterns []*regexp.Regexp
}

// NewRedactor creates a new Redactor instance.
func NewRedactor() *Redactor {
	// Sensitive key names (case-insensitive)
	sensitiveKeys := map[string]bool{
		"api_key":       true,
		"token":         true,
		"secret":        true,
		"password":      true,
		"authorization": true,
		"credential":    true,
		"private_key":   true,
		"access_key":    true,
		"secret_key":    true,
		"conn_string":   true,
		"passwd":        true,
	}

	// Sensitive value patterns
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`^sk-[A-Za-z0-9_-]{20,}$`),                    // OpenAI API keys
		regexp.MustCompile(`^AKIA[A-Z0-9]{16}$`),                        // AWS access keys
		regexp.MustCompile(`^AIza[A-Za-z0-9_-]{35}$`),                   // Google API keys
		regexp.MustCompile(`^gh[ps]_[A-Za-z0-9_]{36,}$`),                // GitHub tokens
		regexp.MustCompile(`^gho_[A-Za-z0-9_]{36,}$`),                   // GitHub OAuth
		regexp.MustCompile(`^ghu_[A-Za-z0-9_]{36,}$`),                   // GitHub user tokens
		regexp.MustCompile(`^Bearer\s+[A-Za-z0-9_\-\.]+$`),              // Bearer tokens
		regexp.MustCompile(`-----BEGIN\s+.*PRIVATE\s+KEY-----`),         // Private keys
		regexp.MustCompile(`^[A-Fa-f0-9]{40,}$`),                        // Long hex strings
		regexp.MustCompile(`^[A-Za-z0-9+/]{40,}={0,2}$`),                // Long base64 strings
	}

	return &Redactor{
		sensitiveKeys:   sensitiveKeys,
		sensitivePatterns: patterns,
	}
}

// Redact masks sensitive data in a JSON payload.
func (r *Redactor) Redact(payload json.RawMessage) json.RawMessage {
	if len(payload) == 0 {
		return payload
	}

	var data interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		// If parsing fails, return original
		return payload
	}

	redacted := r.redactValue(data, 0)
	result, err := json.Marshal(redacted)
	if err != nil {
		return payload
	}

	return result
}

func (r *Redactor) redactValue(value interface{}, depth int) interface{} {
	// Depth limit to prevent infinite recursion
	if depth > maxRedactDepth {
		return value
	}

	switch v := value.(type) {
	case map[string]interface{}:
		return r.redactMap(v, depth)
	case []interface{}:
		return r.redactSlice(v, depth)
	case string:
		if r.isSensitiveValue(v) {
			return redactedValue
		}
		return v
	default:
		return v
	}
}

func (r *Redactor) redactMap(m map[string]interface{}, depth int) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		if r.isSensitiveKey(k) {
			result[k] = redactedValue
		} else {
			result[k] = r.redactValue(v, depth+1)
		}
	}
	return result
}

func (r *Redactor) redactSlice(s []interface{}, depth int) []interface{} {
	result := make([]interface{}, len(s))
	for i, v := range s {
		result[i] = r.redactValue(v, depth+1)
	}
	return result
}

func (r *Redactor) isSensitiveKey(key string) bool {
	lowerKey := strings.ToLower(key)
	return r.sensitiveKeys[lowerKey]
}

func (r *Redactor) isSensitiveValue(value string) bool {
	// Check if value matches any sensitive pattern
	for _, pattern := range r.sensitivePatterns {
		if pattern.MatchString(value) {
			return true
		}
	}
	return false
}
