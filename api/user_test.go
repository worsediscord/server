package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/worsediscord/server/services/user"
	"github.com/worsediscord/server/services/user/usertest"
	"github.com/worsediscord/server/util"
)

func TestServer_HandleUserList(t *testing.T) {
	svc := usertest.FakeUserService{
		ExpectedUsers: []*user.User{
			{Username: "spiderman", Nickname: "spidey", Password: "uncleben123"},
		},
		ExpectedError: nil,
	}

	s := NewServer(&svc, nil, nil, nil, util.NopLogHandler{})

	tests := map[string]struct {
		request          *http.Request
		recorder         *httptest.ResponseRecorder
		expectedStatus   int
		expectedResponse []UserResponse
	}{
		"valid": {
			httptest.NewRequest(http.MethodGet, "/api/users", nil),
			httptest.NewRecorder(),
			200,
			[]UserResponse{{Username: "spiderman", Nickname: "spidey"}},
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			s.handleUserList()(input.recorder, input.request)
			if input.recorder.Code != input.expectedStatus {
				t.Fatalf("got status %d, expected %d", input.recorder.Code, input.expectedStatus)
			}

			var response []UserResponse
			if err := json.NewDecoder(input.recorder.Body).Decode(&response); err != nil {
				t.Fatal(err)
			}

			for _, responseUser := range response {
				for _, expectedUser := range input.expectedResponse {
					if responseUser.Nickname != expectedUser.Nickname ||
						responseUser.Username != expectedUser.Username {
						t.Fatalf("got user %v, expected %v", responseUser, expectedUser)
					}
				}
			}
		})
	}
}
