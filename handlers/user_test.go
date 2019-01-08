package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func TestCreateUserHandler(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	handler := fakeHandler(nil)

	// Should return 401 if no access token
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
	expected := `{"code":401,"message":"access denied"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Should return 400 if request body is nil
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer tokenwithinvaliduserid")

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
	expected = `{"code":400,"message":"request body is nil"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Should return 400 if request body is not valid json
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/users", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer tokenwithinvaliduserid")

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
	expected = `{"code":400,"message":"failed to parse request body as json"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Should return 500 if user service error
	handler = fakeHandler(&fakeHandlerOptions{
		userServiceReturnError: true,
	})
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/users", strings.NewReader(`{
		"email": "test@test.com",
		"password": "password"
	}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer correcttoken")

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	expected = `{"code":500,"message":"internal server error"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Should return 400 if email invalid
	handler = fakeHandler(nil)
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/users", strings.NewReader(`{
		"email": "test",
		"password": "password"
	}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer correcttoken")

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
	expected = `{"code":400,"message":"invalid email: 'test' is not a valid email"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Should return 400 if password invalid
	handler = fakeHandler(nil)
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/users", strings.NewReader(`{
		"email": "test@test.com",
		"password": "invalid"
	}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer correcttoken")

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
	expected = `{"code":400,"message":"invalid password: password should be at least 8 characters"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Should return 200 with the created user
	handler = fakeHandler(nil)
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/users", strings.NewReader(`{
		"email": "test@test.com",
		"password": "password"
	}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer correcttoken")

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
	expected = `{"id":1,"email":"test@test.com","admin":false}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Should return 200 with the created admin user
	handler = fakeHandler(nil)
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/users", strings.NewReader(`{
		"email": "test@test.com",
		"password": "password",
		"admin": true
	}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer correcttoken")

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
	expected = `{"id":1,"email":"test@test.com","admin":true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetUserHandler(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	handler := fakeHandler(nil)

	// Should return 401 if no access token
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
	expected := `{"code":401,"message":"access denied"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Should return 500 if user service error
	handler = fakeHandler(&fakeHandlerOptions{
		userServiceReturnError: true,
	})
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer correcttoken")

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	expected = `{"code":500,"message":"internal server error"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Should return 403 if access unauthorize resource
	handler = fakeHandler(nil)
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/resources/resource2", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer correcttoken")

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}
	expected = `{"code":403,"message":"access denied"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Should return 404 if user not found
	handler = fakeHandler(nil)
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/users/2", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer correcttoken")

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
	expected = `{"code":404,"message":"user not found"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Should return 404 if user id invalid
	handler = fakeHandler(nil)
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/users/invalid", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer correcttoken")

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
	expected = `{"code":404,"message":"user not found"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Should return 200 with requested user
	handler = fakeHandler(nil)
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer correcttoken")

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected = `{"id":1,"email":"test@test.com","admin":false,"quota":-1}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestDeleteUserHandler(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	handler := fakeHandler(nil)

	// Should return 401 if no access token
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
	expected := `{"code":401,"message":"access denied"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Should return 500 if user service error
	handler = fakeHandler(&fakeHandlerOptions{
		userServiceReturnError: true,
	})
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("DELETE", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer correcttoken")

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	expected = `{"code":500,"message":"internal server error"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Should return 404 if user not found
	handler = fakeHandler(&fakeHandlerOptions{})
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("DELETE", "/users/2", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer correcttoken")

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
	expected = `{"code":404,"message":"user not found"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Should return 200 with no content
	handler = fakeHandler(nil)
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("DELETE", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer correcttoken")

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}
	expected = ``
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestListUsersHandler(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	handler := fakeHandler(nil)

	// Should return 401 if no access token
	rr := httptest.NewRecorder()
	params := url.Values{}
	req, err := http.NewRequest("GET", "/users", strings.NewReader(params.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
	expected := `{"code":401,"message":"access denied"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Should return 500 if user service error
	handler = fakeHandler(&fakeHandlerOptions{
		userServiceReturnError: true,
	})
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer correcttoken")

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	expected = `{"code":500,"message":"internal server error"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Should return 200 with all users
	handler = fakeHandler(nil)
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer correcttoken")

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected = `[{"id":1,"email":"test@test.com","admin":false,"quota":-1}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
