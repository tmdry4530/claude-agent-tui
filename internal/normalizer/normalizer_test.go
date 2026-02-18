package normalizer

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/chamdom/omc-agent-tui/pkg/schema"
)

func TestNormalizer_Normalize_Success(t *testing.T) {
	n := New()

	rawData := map[string]interface{}{
		"ts":        "2024-01-01T12:00:00Z",
		"run_id":    "run-123",
		"provider":  "claude",
		"mode":      "autopilot",
		"agent_id":  "agent-456",
		"role":      "executor",
		"state":     "running",
		"type":      "task_spawn",
		"task_id":   "task-789",
		"payload": map[string]interface{}{
			"message": "test message",
		},
	}

	dataBytes, _ := json.Marshal(rawData)
	raw := schema.RawEvent{
		Source:   "test",
		Data:     dataBytes,
		Received: time.Now(),
	}

	event, err := n.Normalize(raw)
	if err != nil {
		t.Fatalf("Normalize failed: %v", err)
	}

	if event.RunID != "run-123" {
		t.Errorf("Expected run_id 'run-123', got '%s'", event.RunID)
	}
	if event.AgentID != "agent-456" {
		t.Errorf("Expected agent_id 'agent-456', got '%s'", event.AgentID)
	}
	if event.Provider != schema.ProviderClaude {
		t.Errorf("Expected provider 'claude', got '%s'", event.Provider)
	}
	if event.Mode != schema.ModeAutopilot {
		t.Errorf("Expected mode 'autopilot', got '%s'", event.Mode)
	}
	if event.Role != schema.RoleExecutor {
		t.Errorf("Expected role 'executor', got '%s'", event.Role)
	}
	if event.State != schema.StateRunning {
		t.Errorf("Expected state 'running', got '%s'", event.State)
	}
	if event.Type != schema.TypeTaskSpawn {
		t.Errorf("Expected type 'task_spawn', got '%s'", event.Type)
	}
}

func TestNormalizer_UnknownProvider(t *testing.T) {
	n := New()

	rawData := map[string]interface{}{
		"run_id":   "run-123",
		"agent_id": "agent-456",
		"provider": "unknown-provider",
		"role":     "executor",
		"state":    "running",
		"type":     "task_spawn",
	}

	dataBytes, _ := json.Marshal(rawData)
	raw := schema.RawEvent{
		Source:   "test",
		Data:     dataBytes,
		Received: time.Now(),
	}

	event, err := n.Normalize(raw)
	if err != nil {
		t.Fatalf("Normalize failed: %v", err)
	}

	if event.Provider != schema.ProviderSystem {
		t.Errorf("Expected provider 'system' for unknown provider, got '%s'", event.Provider)
	}
}

func TestNormalizer_RoleMapping(t *testing.T) {
	n := New()

	tests := []struct {
		name       string
		agentType  string
		expectRole schema.Role
	}{
		{"security-reviewer maps to guard", "security-reviewer", schema.RoleGuard},
		{"executor maps to executor", "executor", schema.RoleExecutor},
		{"planner maps to planner", "planner", schema.RolePlanner},
		{"unknown agent type maps to custom", "unknown-agent", schema.RoleCustom},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rawData := map[string]interface{}{
				"run_id":     "run-123",
				"agent_id":   "agent-456",
				"provider":   "claude",
				"agent_type": tt.agentType,
				"state":      "running",
				"type":       "task_spawn",
			}

			dataBytes, _ := json.Marshal(rawData)
			raw := schema.RawEvent{
				Source:   "test",
				Data:     dataBytes,
				Received: time.Now(),
			}

			event, err := n.Normalize(raw)
			if err != nil {
				t.Fatalf("Normalize failed: %v", err)
			}

			if event.Role != tt.expectRole {
				t.Errorf("Expected role '%s', got '%s'", tt.expectRole, event.Role)
			}
		})
	}
}

func TestNormalizer_PayloadRedaction(t *testing.T) {
	n := New()

	tests := []struct {
		name          string
		payload       map[string]interface{}
		expectRedacted bool
		checkKey      string
	}{
		{
			name: "api_key field redacted",
			payload: map[string]interface{}{
				"api_key": "sk-1234567890abcdefghij",
				"message": "safe data",
			},
			expectRedacted: true,
			checkKey:       "api_key",
		},
		{
			name: "sk- pattern redacted",
			payload: map[string]interface{}{
				"config": map[string]interface{}{
					"key": "sk-proj-abcdefghijklmnopqrstuvwxyz1234567890",
				},
			},
			expectRedacted: true,
			checkKey:       "config.key",
		},
		{
			name: "nested secret redacted",
			payload: map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": map[string]interface{}{
						"secret": "my-secret-value",
					},
				},
			},
			expectRedacted: true,
			checkKey:       "level1.level2.secret",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rawData := map[string]interface{}{
				"run_id":   "run-123",
				"agent_id": "agent-456",
				"provider": "claude",
				"role":     "executor",
				"state":    "running",
				"type":     "task_spawn",
				"payload":  tt.payload,
			}

			dataBytes, _ := json.Marshal(rawData)
			raw := schema.RawEvent{
				Source:   "test",
				Data:     dataBytes,
				Received: time.Now(),
			}

			event, err := n.Normalize(raw)
			if err != nil {
				t.Fatalf("Normalize failed: %v", err)
			}

			if len(event.Payload) == 0 {
				t.Fatal("Payload is empty")
			}

			var redactedPayload map[string]interface{}
			if err := json.Unmarshal(event.Payload, &redactedPayload); err != nil {
				t.Fatalf("Failed to unmarshal payload: %v", err)
			}

			// Check for redaction
			contains := containsRedacted(redactedPayload)
			if tt.expectRedacted && !contains {
				t.Errorf("Expected redaction in payload, but found none")
			}
		})
	}
}

func TestNormalizer_InvalidJSON(t *testing.T) {
	n := New()

	raw := schema.RawEvent{
		Source:   "test",
		Data:     []byte("invalid json"),
		Received: time.Now(),
	}

	_, err := n.Normalize(raw)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestNormalizer_MissingRequiredFields(t *testing.T) {
	n := New()

	tests := []struct {
		name    string
		data    map[string]interface{}
		wantErr bool
	}{
		{
			name: "missing run_id",
			data: map[string]interface{}{
				"agent_id": "agent-456",
				"provider": "claude",
				"role":     "executor",
				"state":    "running",
				"type":     "task_spawn",
			},
			wantErr: true,
		},
		{
			name: "missing agent_id",
			data: map[string]interface{}{
				"run_id":   "run-123",
				"provider": "claude",
				"role":     "executor",
				"state":    "running",
				"type":     "task_spawn",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataBytes, _ := json.Marshal(tt.data)
			raw := schema.RawEvent{
				Source:   "test",
				Data:     dataBytes,
				Received: time.Now(),
			}

			_, err := n.Normalize(raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("wantErr=%v, got err=%v", tt.wantErr, err)
			}
		})
	}
}

// Helper function to check if payload contains redacted values
func containsRedacted(data interface{}) bool {
	switch v := data.(type) {
	case map[string]interface{}:
		for _, val := range v {
			if str, ok := val.(string); ok && str == "***REDACTED***" {
				return true
			}
			if containsRedacted(val) {
				return true
			}
		}
	case []interface{}:
		for _, val := range v {
			if containsRedacted(val) {
				return true
			}
		}
	case string:
		return v == "***REDACTED***"
	}
	return false
}
