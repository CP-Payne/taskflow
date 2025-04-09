package events

import "encoding/json"

const (
	ChannelTaskAssigned = "events:task:assigned"
)

type TaskAssignedEvent struct {
	TaskID string `json:"taskId"`
	UserID string `json:"userId"`
}

// Marshal encodes the event into JSON bytes.
func (e *TaskAssignedEvent) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

// UnmarshalTaskAssignedEvent decodes JSON bytes into an event.
func UnmarshalTaskAssignedEvent(data []byte) (*TaskAssignedEvent, error) {
	var event TaskAssignedEvent
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}
