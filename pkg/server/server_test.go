package server

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ebuckley/bff/pkg/bff"
)

func testBff(t *testing.T) *bff.BFF {
	t.Helper()
	bffInstance := bff.New()
	err := bffInstance.RegisterAction("some-action", func(ctx context.Context, io *bff.Io) error {
		return nil
	})
	if err != nil {
		t.Fatal("sregistering some-action", err)
	}
	return bffInstance
}

func TestServer_SetPrefix(t *testing.T) {
	bffInstance := testBff(t)
	server := NewServer(bffInstance, bff.WithPrefix("/dashboard"))

	if server.handlerPrefix != "/dashboard" {
		t.Errorf("expected prefix to be /dashboard, got %s", server.handlerPrefix)
	}

	// the react app should include the prefix in the links too...
	// the static assets in the app should be prefixed. Everything needs to be prefixed!
}

func TestServer_ReturnIndexPage(t *testing.T) {
	bffInstance := testBff(t)
	server := NewServer(bffInstance, "")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %v", resp.Status)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "some-action") {
		t.Errorf("expected aome-action in body, got %s", b)
	}
}

func TestServer_GetActionPage(t *testing.T) {
	bffInstance := testBff(t)

	server := NewServer(bffInstance, "")

	req := httptest.NewRequest(http.MethodGet, "/a/some-action", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %v", resp.Status)
	}
	//	TODO test this is a react app somehow
}

func TestServer_IndexPageWithPrefix(t *testing.T) {
	bffInstance := testBff(t)
	server := NewServer(bffInstance, "/dashboard")

	req := httptest.NewRequest(http.MethodGet, "/dashboard/", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %v", resp.Status)
	}
}

func TestServer_InvalidEnvironment(t *testing.T) {
	bffInstance := &bff.BFF{}
	server := NewServer(bffInstance, "")

	req := httptest.NewRequest(http.MethodGet, "/e/invalid-env", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status NotFound, got %v", resp.Status)
	}
}
