package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/unrolled/render"

	"github.com/moonkeat/chainstack/handlers"
)

func TestTokenHandler(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	handler := handlers.NewHandler(&handlers.Env{
		Render: render.New(),
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

	// Should return 400 if no grant_type invalid
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

	// Should return 200 if request valid
	rr = httptest.NewRecorder()
	params = url.Values{}
	params.Set("grant_type", "client_credentials")
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
