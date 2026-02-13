package client

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-github/v82/github"
)

func TestNewGitHubClient(t *testing.T) {
	ctx := context.TODO()
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("authorization")
		if want := "Bearer YOUR_TOKEN"; authorization != want {
			t.Errorf("authorization wants %s but %s (%s)", want, authorization, r.URL)
		}
		switch {
		case r.Method == "POST" && r.URL.Path == "/api/graphql":
			w.Header().Set("content-type", "application/json")
			if _, err := fmt.Fprint(w, "{}"); err != nil {
				t.Errorf("error while writing body: %s", err)
			}
		case r.Method == "POST" && r.URL.Path == "/api/v3/repos/owner/repo/git/blobs":
			w.Header().Set("content-type", "application/json")
			if _, err := fmt.Fprint(w, "{}"); err != nil {
				t.Errorf("error while writing body: %s", err)
			}
		default:
			t.Logf("Not found: %s %s", r.Method, r.URL)
			http.NotFound(w, r)
		}
	}))
	defer s.Close()

	o := Option{
		Token: "YOUR_TOKEN",
		URLv3: s.URL + "/api/v3/",
	}
	c, err := New(o)
	if err != nil {
		t.Fatalf("Init returned error: %s", err)
	}

	// v4 API
	var q struct{}
	var v map[string]interface{}
	if err := c.Query(ctx, &q, v); err != nil {
		t.Errorf("Query returned error: %s", err)
	}

	// v3 API
	if _, _, err := c.CreateBlob(ctx, "owner", "repo", &github.Blob{}); err != nil {
		t.Errorf("CreateBlob returned error: %s", err)
	}
}
