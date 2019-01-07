package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/moonkeat/chainstack/handlers"
	"github.com/rs/zerolog"
	"github.com/unrolled/render"
)

func TestListResourcesHandler(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	handler := handlers.NewHandler(&handlers.Env{
		Render:          render.New(),
		UserService:     &fakeUserService{},
		TokenService:    &fakeTokenService{},
		ResourceService: &fakeResourceService{},
	})

	// Should return 401 if no access token
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

	// Should return 500 if token associate with invalid user id
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/resources", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer tokenwithinvaliduserid")

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

	// Should return 500 if failed to get tokens
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/resources", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer tokenserviceerror")

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

	// Should return 200 with all resources belong to the user
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
	// createdAt returned from fake resource service
	createdAt := time.Now().Truncate(24 * time.Hour).Format(time.RFC3339Nano)
	expected = `[{"key":"resource1","created_at":"` + createdAt + `"}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}
