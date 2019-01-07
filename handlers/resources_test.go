package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/moonkeat/chainstack/handlers"
	"github.com/rs/zerolog"
	"github.com/unrolled/render"
)

func TestListResourcesHandler(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	handler := handlers.NewHandler(&handlers.Env{
		Render:       render.New(),
		UserService:  &fakeUserService{},
		TokenService: &fakeTokenService{},
	})

	// Should return 400 if no request body
	rr := httptest.NewRecorder()
	params := url.Values{}
	req, err := http.NewRequest("GET", "/resources", strings.NewReader(params.Encode()))
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

	// Should return 400 if no request body
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/resources", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer correcttoken")

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected = `[]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}
