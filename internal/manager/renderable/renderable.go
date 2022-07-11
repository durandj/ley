package renderable

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/render"
)

// Time makes time values renderable in API responses.
type Time time.Time

// MarshalJSON converts a time value into an ISO-8601 string.
func (renderableTime *Time) MarshalJSON() ([]byte, error) {
	stringTime := time.Time(*renderableTime).Format(time.RFC3339)

	return []byte("\"" + stringTime + "\""), nil
}

// UnmarshalJSON converts an ISO-8601 time string into a time instance.
func (renderableTime *Time) UnmarshalJSON(rawBytes []byte) error {
	cleanedString := strings.Trim(string(rawBytes), "\"")

	timeValue, err := time.Parse(time.RFC3339, cleanedString)
	if err != nil {
		return fmt.Errorf("Unable to parse as ISO-8601 time string: %w", err)
	}

	*renderableTime = Time(timeValue)

	return nil
}

var _ json.Marshaler = (*Time)(nil)
var _ json.Unmarshaler = (*Time)(nil)

// ErrorResponse is used to send some kind of failure response back
// to the requester.
type ErrorResponse struct {
	Message string `json:"message"`
}

// Render provides a hook to customize the render process.
func (errorResponse *ErrorResponse) Render(
	response http.ResponseWriter,
	request *http.Request,
) error {
	return nil
}

var _ render.Renderer = (*ErrorResponse)(nil)
