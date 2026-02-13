package responses

import (
	"encoding/json"
	"net/http"
)

// Envelope is a generic wrapper for API responses.
// Using a consistent top-level key like "data" or "error" makes it easier
// for clients to parse responses.
type Envelope map[string]interface{}

// WriteJSON is a helper function for sending JSON responses.
// It sets the Content-Type header, the given status code, and writes the
// JSON-encoded body to the http.ResponseWriter.
func WriteJSON(w http.ResponseWriter, r *http.Request, status int, data interface{}, headers http.Header) error {
	// Set any provided headers.
	for key, value := range headers {
		w.Header()[key] = value
	}

	// Set the Content-Type header to application/json.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Create an envelope for the response data.
	// This provides a consistent structure for all successful responses.
	env := Envelope{"data": data}

	// Use json.Encoder to encode the envelope and write it to the response writer.
	// This is more efficient than json.Marshal for writing directly to an io.Writer.
	err := json.NewEncoder(w).Encode(env)
	if err != nil {
		// If encoding fails, return the error.
		// The caller (usually a handler) should handle this by logging and
		// potentially sending a plain-text error response.
		return err
	}

	return nil
}
