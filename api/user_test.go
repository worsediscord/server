package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/worsediscord/server/services/auth"
	"github.com/worsediscord/server/services/fake"
	"github.com/worsediscord/server/services/user"
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
			userService:    &fake.UserService{ExpectedCreateError: nil, ExpectedGetUserByIdError: user.ErrNotFound},
			expectedStatus: http.StatusOK,
		},
		"invalid": {
			request:        httptest.NewRequest(http.MethodPost, "/api/users", util.StructToReaderOrDie(invalidRequest)),
			recorder:       httptest.NewRecorder(),
			userService:    &fake.UserService{},
			expectedStatus: http.StatusBadRequest,
		},
		"conflict": {
			request:        httptest.NewRequest(http.MethodPost, "/api/users", util.StructToReaderOrDie(validRequest)),
			recorder:       httptest.NewRecorder(),
			userService:    &fake.UserService{ExpectedGetUserByIdUser: &user.User{Username: "spiderman"}},
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
	validResponse := []*user.User{{Username: "spiderman", Nickname: "spidey", Password: "uncleben123"}}
	emptyResponse := make([]*user.User, 0)

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
			userService:      &fake.UserService{ExpectedListUsers: validResponse},
			expectedStatus:   http.StatusOK,
			expectedResponse: []UserResponse{{Username: "spiderman", Nickname: "spidey"}},
		},
		"empty": {
			request:          httptest.NewRequest(http.MethodGet, "/api/users", nil),
			recorder:         httptest.NewRecorder(),
			userService:      &fake.UserService{ExpectedListUsers: emptyResponse},
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
	s := NewServer(nil, nil, nil, nil, util.NopLogHandler)
	validResponse := &user.User{Username: "spiderman", Nickname: "spidey", Password: "uncleben123"}
	emptyResponse := &user.User{}

	tests := map[string]struct {
		request          *http.Request
		recorder         *httptest.ResponseRecorder
		userService      user.Service
		expectedStatus   int
		expectedResponse UserResponse
	}{
		"valid": {
			request:          httptest.NewRequest(http.MethodGet, "/api/users/spiderman", nil),
			recorder:         httptest.NewRecorder(),
			userService:      &fake.UserService{ExpectedGetUserByIdUser: validResponse},
			expectedStatus:   http.StatusOK,
			expectedResponse: UserResponse{Username: "spiderman", Nickname: "spidey"},
		},
		"empty": {
			request:          httptest.NewRequest(http.MethodGet, "/api/users/batman", nil),
			recorder:         httptest.NewRecorder(),
			userService:      &fake.UserService{ExpectedGetUserByIdUser: emptyResponse, ExpectedGetUserByIdError: user.ErrNotFound},
			expectedStatus:   http.StatusNotFound,
			expectedResponse: UserResponse{},
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			s.UserService = input.userService
			s.handleUserGet()(input.recorder, input.request)

			if input.recorder.Code != input.expectedStatus {
				t.Fatalf("got status %d, expected %d", input.recorder.Code, input.expectedStatus)
			}

			var response UserResponse
			// 404 returns an empty body, should probably do something better somewhere
			_ = json.NewDecoder(input.recorder.Body).Decode(&response)

			if response.Nickname != input.expectedResponse.Nickname || response.Username != input.expectedResponse.Username {
				t.Fatalf("got user %v, expected %v", response, input.expectedResponse)
			}
		})
	}
}

func TestServer_HandleUserLogin(t *testing.T) {
	s := NewServer(nil, nil, nil, nil, util.NopLogHandler)

	validRequest := httptest.NewRequest(http.MethodGet, "/api/users/login", nil)
	validRequest.SetBasicAuth("spiderman", "uncleben123")
	validResponse := &user.User{Username: "spiderman", Nickname: "spidey", Password: "uncleben123"}

	invalidRequest := httptest.NewRequest(http.MethodGet, "/api/users/login", nil)
	invalidRequest.SetBasicAuth("batman", "iamthenight")

	tests := map[string]struct {
		request          *http.Request
		recorder         *httptest.ResponseRecorder
		userService      user.Service
		authService      auth.Service
		expectedStatus   int
		expectedResponse UserResponse
	}{
		"valid": {
			request:          httptest.NewRequest(http.MethodGet, "/api/users/login", nil),
			recorder:         httptest.NewRecorder(),
			userService:      &fake.UserService{ExpectedGetUserByIdUser: validResponse},
			authService:      &fake.AuthService{},
			expectedStatus:   http.StatusOK,
			expectedResponse: UserResponse{Username: "spiderman", Nickname: "spidey"},
		},
		"empty": {
			request:        httptest.NewRequest(http.MethodGet, "/api/users/login", nil),
			recorder:       httptest.NewRecorder(),
			userService:    &fake.UserService{ExpectedGetUserByIdUser: &user.User{}, ExpectedGetUserByIdError: user.ErrNotFound},
			authService:    &fake.AuthService{},
			expectedStatus: http.StatusNotFound,
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			s.UserService = input.userService
			s.handleUserGet()(input.recorder, input.request)

			if input.recorder.Code != input.expectedStatus {
				t.Fatalf("got status %d, expected %d", input.recorder.Code, input.expectedStatus)
			}

			var response UserResponse
			// 404 returns an empty body, should probably do something better somewhere
			_ = json.NewDecoder(input.recorder.Body).Decode(&response)

			if response.Nickname != input.expectedResponse.Nickname || response.Username != input.expectedResponse.Username {
				t.Fatalf("got user %v, expected %v", response, input.expectedResponse)
			}
		})
	}
}

func TestServer_HandleUserDelete(t *testing.T) {

}
