package services

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChatbotService(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		mockResponse   string
		expectedError  bool
		expectedResult string
	}{
		{
			name:  "successful response",
			query: "Hello",
			mockResponse: `{
                "choices": [
                    {
                        "message": {
                            "content": "Hello, how are you!"
                        }
                    }
                ]
            }`,
			expectedError:  false,
			expectedResult: "Hello, how are you!",
		},
		{
			name:  "empty choices",
			query: "Hello",
			mockResponse: `{
                "choices": []
            }`,
			expectedError:  true,
			expectedResult: "",
		},
		{
			name:           "invalid json response",
			query:          "Hello",
			mockResponse:   `invalid json`,
			expectedError:  true,
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
				assert.Equal(t, defaultAuthToken, r.Header.Get("Authorization"))

				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			config := &ChatbotConfig{
				URL:           server.URL,
				Authorization: defaultAuthToken,
			}

			result, err := ChatbotService(tt.query, config)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}
