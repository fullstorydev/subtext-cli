package sightmap

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	sm "github.com/sightmap/sightmap/go/sightmap"

	"github.com/fullstorydev/subtext-cli/internal/fstesting"
)

func TestUpload_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fstesting.Equals(t, http.MethodPost, r.Method, "HTTP method")
		fstesting.Equals(t, "application/json", r.Header.Get("Content-Type"), "Content-Type header")

		var body sm.Corpus
		fstesting.Ok(t, json.NewDecoder(r.Body).Decode(&body), "decode request body")
		fstesting.Equals(t, []string{"hint"}, body.Memory, "memory")
		fstesting.Equals(t, "/home", body.Views[0].Route, "view route")

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok", "components": len(body.AllComponents())})
	}))
	t.Cleanup(srv.Close)

	c := &sm.Corpus{
		Memory: []string{"hint"},
		Views: []sm.ViewDef{{
			Name:       "Home",
			Route:      "/home",
			Components: []sm.ComponentDef{{Name: "Btn", Selectors: []string{".btn"}}},
		}},
	}
	n, err := Upload(context.Background(), srv.URL, c)
	fstesting.Ok(t, err, "Upload")
	fstesting.Equals(t, 1, n, "acknowledged component count")
}

// TestUpload_NoAuthHeader verifies the upload does not add an Authorization
// header — the token is already embedded in the URL's query params.
func TestUpload_NoAuthHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fstesting.Equals(t, "", r.Header.Get("Authorization"), "no Authorization header expected")
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok", "components": 0})
	}))
	t.Cleanup(srv.Close)

	_, err := Upload(context.Background(), srv.URL, &sm.Corpus{})
	fstesting.Ok(t, err, "Upload")
}

func TestUpload_ExpiredNonce(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "invalid or expired nonce", http.StatusUnauthorized)
	}))
	t.Cleanup(srv.Close)

	_, err := Upload(context.Background(), srv.URL, &sm.Corpus{})
	fstesting.Assert(t, err != nil, "expected error for 401")
	fstesting.Assert(t, strings.Contains(err.Error(), "invalid or expired nonce"), "error should include server body: %v", err)
}

func TestUpload_BodyTooLarge(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "body too large", http.StatusRequestEntityTooLarge)
	}))
	t.Cleanup(srv.Close)

	_, err := Upload(context.Background(), srv.URL, &sm.Corpus{})
	fstesting.Assert(t, err != nil, "expected error for 413")
	fstesting.Assert(t, strings.Contains(err.Error(), "413"), "error should include status code: %v", err)
}

func TestIsLocalHost(t *testing.T) {
	cases := []struct {
		url  string
		want bool
	}{
		{"https://api.example.test:8043/mcp/subtext", true},
		{"https://localhost:3000/upload", true},
		{"https://127.0.0.1:8080/upload", true},
		{"https://api.fullstory.com/mcp/subtext", false},
		{"https://api.eu1.fullstory.com/mcp/subtext", false},
		{"https://st.fullstory.com/subtext/sightmap", false},
	}
	for _, tc := range cases {
		t.Run(tc.url, func(t *testing.T) {
			fstesting.Equals(t, tc.want, isLocalHost(tc.url), "isLocalHost(%q)", tc.url)
		})
	}
}
