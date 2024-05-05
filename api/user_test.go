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

func TestServer_HandleUserCreate(t *testing.T) {
	s := NewServer(nil, nil, nil, nil, util.NopLogHandler)
	validRequest := UserCreateRequest{Username: "spiderman", Password: "password123"}
	invalidRequest := UserCreateRequest{Username: "batman", Password: ""}

	tests := map[string]struct {
		request        *http.Request
		recorder       *httptest.ResponseRecorder
		userService    user.Service
		expectedStatus int
	}{
		"valid": {
			request:        httptest.NewRequest(http.MethodPost, "/api/users", util.StructToReaderOrDie(validRequest)),
			recorder:       httptest.NewRecorder(),
			userService:    &usertest.FakeUserService{ExpectedCreateError: nil, ExpectedGetUserByIdError: user.ErrNotFound},
			expectedStatus: http.StatusOK,
		},
		"invalid": {
			request:        httptest.NewRequest(http.MethodPost, "/api/users", util.StructToReaderOrDie(invalidRequest)),
			recorder:       httptest.NewRecorder(),
			userService:    &usertest.FakeUserService{},
			expectedStatus: http.StatusBadRequest,
		},
		"conflict": {
			request:        httptest.NewRequest(http.MethodPost, "/api/users", util.StructToReaderOrDie(validRequest)),
			recorder:       httptest.NewRecorder(),
			userService:    &usertest.FakeUserService{ExpectedGetUserByIdUser: &user.User{Username: "spiderman"}},
			expectedStatus: http.StatusConflict,
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			s.UserService = input.userService
			s.handleUserCreate()(input.recorder, input.request)

			if input.recorder.Code != input.expectedStatus {
				t.Fatalf("got status %d, expected %d", input.recorder.Code, input.expectedStatus)
			}
		})
	}
}

func TestServer_HandleUserList(t *testing.T) {
	s := NewServer(nil, nil, nil, nil, util.NopLogHandler)
	emptyResponse := []*user.User{}
	nonEmptyResponse := []*user.User{{Username: "spiderman", Nickname: "spidey", Password: "uncleben123"}}

	tests := map[string]struct {
		request          *http.Request
		recorder         *httptest.ResponseRecorder
		userService      user.Service
		expectedStatus   int
		expectedResponse []UserResponse
	}{
		"valid": {
			request:          httptest.NewRequest(http.MethodGet, "/api/users", nil),
			recorder:         httptest.NewRecorder(),
			userService:      &usertest.FakeUserService{ExpectedListUsers: nonEmptyResponse},
			expectedStatus:   http.StatusOK,
			expectedResponse: []UserResponse{{Username: "spiderman", Nickname: "spidey"}},
		},
		"empty": {
			request:          httptest.NewRequest(http.MethodGet, "/api/users", nil),
			recorder:         httptest.NewRecorder(),
			userService:      &usertest.FakeUserService{ExpectedListUsers: emptyResponse},
			expectedStatus:   http.StatusOK,
			expectedResponse: []UserResponse{},
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			s.UserService = input.userService
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

func TestServer_HandleUserGet(t *testing.T) {

}

func TestServer_HandleUserLogin(t *testing.T) {

}

func TestServer_HandleUserDelete(t *testing.T) {

}
