package tslab

import (
	"os"
)

// EventWriter make our own file type
type EventWriter *os.File

// NewEventWriter creates if not exists, appends and sets 0644 perms.
func NewEventWriter(fName string) (*os.File, error) {

	f, err := os.OpenFile(fName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return nil, err
	}
	return f, nil
}
