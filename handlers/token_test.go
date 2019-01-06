package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/unrolled/render"

	"github.com/moonkeat/chainstack/handlers"
)

type fakeUserService struct{}

func (s fakeUserService) CreateUser(email string, password string, isAdmin bool) error {
	return nil
}

func (s fakeUserService) AuthenticateUser(email string, password string) (bool, error) {
	if email == "internalerror" {
		return false, fmt.Errorf("intermal server error occurred")
	}

	if email != "correct@email.com" || password != "correctpassword" {
		return false, nil
	}

	return true, nil
}

func TestTokenHandler(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	handler := handlers.NewHandler(&handlers.Env{
		Render:      render.New(),
		UserService: &fakeUserService{},
	})

	// Should return 400 if no request body
	rr := httptest.NewRecorder()
	params := url.Values{}
	req, err := http.NewRequest("POST", "/token", strings.NewReader(params.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
	expected := `{"code":400,"message":"invalid grant type: ''"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Should return 400 if grant_type invalid
	rr = httptest.NewRecorder()
	params = url.Values{}
	params.Set("grant_type", "authorization_code")
	req, err = http.NewRequest("POST", "/token", strings.NewReader(params.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
	expected = `{"code":400,"message":"invalid grant type: 'authorization_code'"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Should return 400 if no client_id
	rr = httptest.NewRecorder()
	params = url.Values{}
	params.Set("grant_type", "client_credentials")
	params.Set("client_id", "")
	req, err = http.NewRequest("POST", "/token", strings.NewReader(params.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
	expected = `{"code":400,"message":"client_id is required"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Should return 400 if no client_secret
	rr = httptest.NewRecorder()
	params = url.Values{}
	params.Set("grant_type", "client_credentials")
	params.Set("client_id", "some_client_id")
	params.Set("client_secret", "")
	req, err = http.NewRequest("POST", "/token", strings.NewReader(params.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
	expected = `{"code":400,"message":"client_secret is required"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Should return 401 if credential invalid
	rr = httptest.NewRecorder()
	params = url.Values{}
	params.Set("grant_type", "client_credentials")
	params.Set("client_id", "wrong@email.com")
	params.Set("client_secret", "wrongpassword")
	req, err = http.NewRequest("POST", "/token", strings.NewReader(params.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
	expected = `{"code":401,"message":"invalid credentials"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Should return 500 if internal server error occurred
	rr = httptest.NewRecorder()
	params = url.Values{}
	params.Set("grant_type", "client_credentials")
	params.Set("client_id", "internalerror")
	params.Set("client_secret", "anypassword")
	req, err = http.NewRequest("POST", "/token", strings.NewReader(params.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
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

	// Should return 200 if credential valid
	rr = httptest.NewRecorder()
	params = url.Values{}
	params.Set("grant_type", "client_credentials")
	params.Set("client_id", "correct@email.com")
	params.Set("client_secret", "correctpassword")
	req, err = http.NewRequest("POST", "/token", strings.NewReader(params.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected = `{"access_token":"","token_type":"","expires_in":0,"scope":""}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
