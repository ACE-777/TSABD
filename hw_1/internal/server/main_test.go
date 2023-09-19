package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCheckReplaceHandlerUnsuccessfulMethod(t *testing.T) {
	requestBody := []byte(`{"body":"568"}`)
	req, err := http.NewRequest(http.MethodGet, "/replace", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	Replace(w, req)

	if w.Code != http.StatusBadGateway {
		t.Errorf("Expected status code %d, got %d", http.StatusBadGateway, w.Code)
	}

	expectedResponseBody := "method must be PUT"

	if w.Body.String() != expectedResponseBody {
		t.Errorf("Expected response body '%s', got '%s'", expectedResponseBody, w.Body.String())
	}
}

func TestCheckGetHandlerUnsuccessfulMethod(t *testing.T) {
	req, err := http.NewRequest(http.MethodPut, "/get", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	Get(w, req)

	if w.Code != http.StatusBadGateway {
		t.Errorf("Expected status code %d, got %d", http.StatusBadGateway, w.Code)
	}

	expectedResponseBody := "method must be GET"

	if w.Body.String() != expectedResponseBody {
		t.Errorf("Expected response body '%s', got '%s'", expectedResponseBody, w.Body.String())
	}
}

func TestCheckBothHandler(t *testing.T) {
	err := os.Chdir("../..")
	if err != nil {
		return
	}

	requestBody := []byte(`{"body":"568gtgfggfgt676"}`)
	reqReplace, err := http.NewRequest(http.MethodPut, "/replace", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	wReplace := httptest.NewRecorder()
	Replace(wReplace, reqReplace)

	if wReplace.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, wReplace.Code)
	}

	expectedResponseBodyReplace := "Successfully save body"

	if wReplace.Body.String() != expectedResponseBodyReplace {
		t.Errorf("Expected response body '%s', got '%s'", expectedResponseBodyReplace, wReplace.Body.String())
	}

	reqGet, err := http.NewRequest(http.MethodGet, "/replace", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	wGet := httptest.NewRecorder()
	Get(wGet, reqGet)

	if wGet.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, wGet.Code)
	}

	expectedResponseBodyGet := "568gtgfggfgt676"

	if wGet.Body.String() != expectedResponseBodyGet {
		t.Errorf("Expected response body '%s', got '%s'", expectedResponseBodyGet, wGet.Body.String())
	}
}
