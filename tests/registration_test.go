package tests

import (
	"bytes"
	"encoding/json"
	transport "go-api-template/gen/oapi"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegistration(t *testing.T) {
	app := setupApp(t)

	tests := []struct {
		name           string
		body           transport.RegistrationRequest
		expectedStatus int
	}{
		{
			name: "Success",
			body: transport.RegistrationRequest{
				Email:    "test@test.ru",
				Password: "Qwerty123",
				Username: "username",
			},
			expectedStatus: 200,
		},
		{
			name: "FailInvalidEmail",
			body: transport.RegistrationRequest{
				Email:    "invalid email",
				Password: "Qwerty123",
				Username: "username",
			},
			expectedStatus: 400,
		},
		{
			name:           "FailEmptyRequestBody",
			body:           transport.RegistrationRequest{},
			expectedStatus: 400,
		},
	}

	is := require.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, err := json.Marshal(&tt.body)
			is.NoError(err)
			reader := bytes.NewReader(jsonBody)

			req, err := http.NewRequest("POST", "http://localhost:8080/registration", reader)
			req.Header.Add("Content-Type", "application/json")
			is.NoError(err)

			resp, err := app.Test(req)
			is.NoError(err)

			is.Equal(tt.expectedStatus, resp.StatusCode)
		})
	}
}
