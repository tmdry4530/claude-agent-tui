package store

import "encoding/json"

// parsePayload is a helper to unmarshal event payload.
func parsePayload(raw json.RawMessage, target any) error {
	if len(raw) == 0 {
		return nil
	}
	return json.Unmarshal(raw, target)
}
