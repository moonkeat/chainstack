package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

func TestCreateResourceHandler(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	handler := fakeHandler(nil)

	// Should return 401 if no access token
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/resources", nil)
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
	req, err = http.NewRequest("POST", "/resources", nil)
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

	// Should return 500 if resource service error
	handler = fakeHandler(&fakeHandlerOptions{
		resourceServiceReturnError: true,
	})
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/resources", nil)
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

	// Should return 200 with all resources belong to the user
	handler = fakeHandler(nil)
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/resources", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer correcttoken")

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
	// createdAt returned from fake resource service
	createdAt := time.Now().Truncate(24 * time.Hour).Format(time.RFC3339Nano)
	expected = `{"key":"resource1","created_at":"` + createdAt + `"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetResourceHandler(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	handler := fakeHandler(nil)

	// Should return 401 if no access token
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/resources/resource1", nil)
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
	req, err = http.NewRequest("GET", "/resources/resource1", nil)
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

	// Should return 500 if resource service error
	handler = fakeHandler(&fakeHandlerOptions{
		resourceServiceReturnError: true,
	})
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/resources/resource1", nil)
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

	// Should return 200 with all resources belong to the user
	handler = fakeHandler(nil)
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/resources/resource1", nil)
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
	expected = `{"key":"resource1","created_at":"` + createdAt + `"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestDeleteResourceHandler(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	handler := fakeHandler(nil)

	// Should return 401 if no access token
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/resources/resource1", nil)
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
	req, err = http.NewRequest("DELETE", "/resources/resource1", nil)
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

	// Should return 500 if resource service error
	handler = fakeHandler(&fakeHandlerOptions{
		resourceServiceReturnError: true,
	})
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("DELETE", "/resources/resource1", nil)
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
	req, err = http.NewRequest("DELETE", "/resources/resource2", nil)
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

	// Should return 200 with all resources belong to the user
	handler = fakeHandler(nil)
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("DELETE", "/resources/resource1", nil)
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

func TestListResourcesHandler(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	handler := fakeHandler(nil)

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

	// Should return 500 if resource service error
	handler = fakeHandler(&fakeHandlerOptions{
		resourceServiceReturnError: true,
	})
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/resources", nil)
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

	// Should return 200 with all resources belong to the user
	handler = fakeHandler(nil)
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
