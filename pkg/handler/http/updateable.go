package http

import (
	"encoding/json"
	"github.com/tomogoma/usersms/pkg/user"
)

type JSONString struct {
	IsUpdating bool   `json:"isUpdating,omitempty"`
	NewValue   string `json:"newValue,omitempty"`
}

func (i *JSONString) UnmarshalJSON(data []byte) error {
	// If this method was called, the value was set in the JSON string.
	if string(data) == "null" { // Ignore the null literal value.
		return nil
	}

	var temp string
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	i.NewValue = temp
	i.IsUpdating = true
	return nil
}

func (i *JSONString) ToStringUpdate() user.StringUpdate {
	return user.StringUpdate{
		IsUpdating: i.IsUpdating,
		NewValue:   i.NewValue,
	}
}
