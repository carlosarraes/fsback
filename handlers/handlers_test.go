package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/carlosarraes/fsback/repository/dbrepo"
	"github.com/go-chi/chi/v5"
)

type ResponseBody struct {
	Message string `json:"message"`
}

func TestGetUsersHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(app.GetUsers)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := []dbrepo.User{
		{
			FirstName: "John",
			LastName:  "Doe",
			Progress:  42,
		},
	}

	var actual []dbrepo.User
	err = json.Unmarshal(rr.Body.Bytes(), &actual)
	if err != nil {
		t.Fatalf("failed to parse actual response: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

func TestDeleteUserHandler(t *testing.T) {
	r := chi.NewRouter()
	r.Delete("/user/{lastName}", app.DeleteUser)

	tests := []struct {
		name           string
		url            string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "delete existing user",
			url:            "/user/Doe",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"User Doe deleted"}`,
		},
		{
			name:           "delete non-existing user",
			url:            "/user/NonExisting",
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"message":"User not found"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", tc.url, nil)
			if err != nil {
				t.Fatalf("could not create DELETE request: %v", err)
			}

			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, req)

			if recorder.Code != tc.expectedStatus {
				t.Errorf("expected status code %d; got %d", tc.expectedStatus, recorder.Code)
				return
			}

			body := recorder.Body.String()
			if body != tc.expectedBody {
				t.Errorf("expected response body '%s'; got '%s'", tc.expectedBody, body)
				return
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name         string
		user         dbrepo.User
		body         []byte
		expectedCode int
		expectedMsg  string
	}{
		{
			name: "Success",
			user: dbrepo.User{
				FirstName: "John",
				LastName:  "Doe",
				Progress:  0.5,
			},
			body:         []byte(`{"FirstName":"John","LastName":"Doe","Progress":0.5}`),
			expectedCode: http.StatusCreated,
			expectedMsg:  "User created",
		},
		{
			name: "InvalidRequest_EmptyName",
			user: dbrepo.User{
				FirstName: "",
				LastName:  "Doe",
				Progress:  0.5,
			},
			body:         []byte(`{"FirstName":"","LastName":"Doe","Progress":0.5}`),
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Error creating user: First name, last name and progress are required",
		},
		{
			name: "SumCheck_Fail",
			user: dbrepo.User{
				FirstName: "John",
				LastName:  "Amount",
				Progress:  1.2,
			},
			body:         []byte(`{"FirstName":"John","LastName":"Amount","Progress":1.2}`),
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Error creating user: Progress sum cannot exceed 100",
		},
		{
			name:         "InvalidRequestBody_Error",
			user:         dbrepo.User{},
			body:         []byte(`invalid json`),
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Error creating user: Invalid request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(tt.body))
			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(app.CreateUser)
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("expected status %d but got %d", tt.expectedCode, rr.Code)
				return
			}

			var respBody ResponseBody
			err := json.Unmarshal(rr.Body.Bytes(), &respBody)
			if err != nil {
				t.Fatalf("failed to parse actual response body %v", err)
			}

			if respBody.Message != tt.expectedMsg {
				t.Errorf("expected response body %s but got %s", tt.expectedMsg, respBody.Message)
			}
		})
	}
}
