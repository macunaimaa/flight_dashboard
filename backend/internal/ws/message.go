package ws

import (
	"encoding/json"

	"github.com/macunaimaa/dashboard/backend/internal/domain"
)

type WSMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func NewAircraftUpdateMsg(states []domain.AircraftLiveState) []byte {
	payload, _ := json.Marshal(states)
	msg := WSMessage{
		Type:    "aircraft.update",
		Payload: payload,
	}
	data, _ := json.Marshal(msg)
	return data
}

func NewAircraftRemoveMsg(icao24s []string) []byte {
	payload, _ := json.Marshal(icao24s)
	msg := WSMessage{
		Type:    "aircraft.remove",
		Payload: payload,
	}
	data, _ := json.Marshal(msg)
	return data
}

func NewAlertTriggeredMsg(event *domain.AlertEvent) []byte {
	payload, _ := json.Marshal(event)
	msg := WSMessage{
		Type:    "alert.triggered",
		Payload: payload,
	}
	data, _ := json.Marshal(msg)
	return data
}

func NewSystemInfoMsg(message string) []byte {
	payload, _ := json.Marshal(map[string]string{"message": message})
	msg := WSMessage{
		Type:    "system.info",
		Payload: payload,
	}
	data, _ := json.Marshal(msg)
	return data
}
