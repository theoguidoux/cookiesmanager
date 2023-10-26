package cookiesmanager_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	cookiesmanager "github.com/theoguidoux/cookiesmanager"
)

func TestRemoveCookies(t *testing.T) {
	cfg := cookiesmanager.CreateConfig()
	cfg.Remover = append(cfg.Remover, cookiesmanager.CookieConfig{
		Name:  "test1",
		Value: "foo",
	})
	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := cookiesmanager.New(ctx, next, cfg, "test")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	cookies := []*http.Cookie{
		{Name: "test1", Value: "value1|foo", Path: "/"},
		{Name: "test2", Value: "value2", Path: "/"},
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	handler.ServeHTTP(recorder, req)

	cookies = req.Cookies()
	if len(cookies) != 2 {
		t.Errorf("there should be 2 cookies in the request, found %d", len(cookies))
	}

	if cookies[0].Value != "value1|" || cookies[0].Name != "test1" {
		t.Error("the expected cookie that should be kept does not match", len(cookies))
	}

	if cookies[1].Value != "value2" || cookies[1].Name != "test2" {
		t.Error("the expected cookie that should be kept does not match", len(cookies))
	}
}

func TestAddNonExistingCookie(t *testing.T) {
	cfg := cookiesmanager.CreateConfig()
	cfg.Adder = append(cfg.Remover, cookiesmanager.CookieConfig{
		Name:  "test3",
		Value: "foo",
	})
	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := cookiesmanager.New(ctx, next, cfg, "test")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	cookies := []*http.Cookie{
		{Name: "test1", Value: "value1", Path: "/"},
		{Name: "test2", Value: "value2", Path: "/"},
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	handler.ServeHTTP(recorder, req)

	cookies = req.Cookies()
	if len(cookies) != 3 {
		t.Errorf("there should be 2 cookies in the request, found %d", len(cookies))
	}

	if cookies[0].Value != "value1" || cookies[0].Name != "test1" {
		t.Error("the expected cookie that should be kept does not match", len(cookies))
	}

	if cookies[1].Value != "value2" || cookies[1].Name != "test2" {
		t.Error("the expected cookie that should be kept does not match", len(cookies))
	}

	if cookies[2].Value != "foo" || cookies[2].Name != "test3" {
		t.Error("the expected cookie that should be kept does not match", len(cookies))
	}
}

func TestUpdateCookie(t *testing.T) {
	cfg := cookiesmanager.CreateConfig()
	cfg.Adder = append(cfg.Remover, cookiesmanager.CookieConfig{
		Name:  "test1",
		Value: "foo",
	})
	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := cookiesmanager.New(ctx, next, cfg, "test")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	cookies := []*http.Cookie{
		{Name: "test1", Value: "value1", Path: "/"},
		{Name: "test2", Value: "value2", Path: "/"},
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	handler.ServeHTTP(recorder, req)

	cookies = req.Cookies()
	if len(cookies) != 2 {
		t.Errorf("there should be 2 cookies in the request, found %d", len(cookies))
	}

	if cookies[0].Value != "value1 foo" || cookies[0].Name != "test1" {
		t.Error("the expected cookie that should be kept does not match", len(cookies))
	}

	if cookies[1].Value != "value2" || cookies[1].Name != "test2" {
		t.Error("the expected cookie that should be kept does not match", len(cookies))
	}
}
