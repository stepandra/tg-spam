package webapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/umputun/tg-spam/app/webapi/mocks"
	"github.com/umputun/tg-spam/lib"
)

func TestServer_Run(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := NewServer(Config{ListenAddr: ":9876", Version: "dev", Detector: &mocks.DetectorMock{},
		SpamFilter: &mocks.SpamFilterMock{}, AuthPasswd: "test"})
	done := make(chan struct{})
	go func() {
		err := srv.Run(ctx)
		assert.NoError(t, err)
		close(done)
	}()
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:9876/ping")
	assert.NoError(t, err)
	t.Log(resp)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "pong", string(body))

	assert.Contains(t, resp.Header.Get("App-Name"), "tg-spam")
	assert.Contains(t, resp.Header.Get("App-Version"), "dev")

	cancel()
	<-done
}

func TestServer_RunAuth(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mockDetector := &mocks.DetectorMock{
		CheckFunc: func(msg string, userID string) (bool, []lib.CheckResult) {
			return false, []lib.CheckResult{{Details: "not spam"}}
		},
	}
	mockSpamFilter := &mocks.SpamFilterMock{}

	srv := NewServer(Config{ListenAddr: ":9877", Version: "dev", Detector: mockDetector, SpamFilter: mockSpamFilter, AuthPasswd: "test"})
	done := make(chan struct{})
	go func() {
		err := srv.Run(ctx)
		assert.NoError(t, err)
		close(done)
	}()
	time.Sleep(100 * time.Millisecond)

	t.Run("ping", func(t *testing.T) {
		resp, err := http.Get("http://localhost:9877/ping")
		assert.NoError(t, err)
		t.Log(resp)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode) // no auth on ping
	})

	t.Run("check unauthorized, no basic auth", func(t *testing.T) {
		resp, err := http.Get("http://localhost:9877/check")
		assert.NoError(t, err)
		t.Log(resp)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("check authorized", func(t *testing.T) {
		reqBody, err := json.Marshal(map[string]string{
			"msg":     "spam example",
			"user_id": "user123",
		})
		require.NoError(t, err)
		req, err := http.NewRequest("POST", "http://localhost:9877/check", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)
		req.SetBasicAuth("tg-spam", "test")
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		t.Log(resp)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("wrong basic auth", func(t *testing.T) {
		reqBody, err := json.Marshal(map[string]string{
			"msg":     "spam example",
			"user_id": "user123",
		})
		require.NoError(t, err)
		req, err := http.NewRequest("POST", "http://localhost:9877/check", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)
		req.SetBasicAuth("tg-spam", "bad")
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		t.Log(resp)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
	cancel()
	<-done
}

func TestServer_routes(t *testing.T) {
	detectorMock := &mocks.DetectorMock{
		CheckFunc: func(msg string, userID string) (bool, []lib.CheckResult) {
			return false, []lib.CheckResult{{Details: "not spam"}}
		},
		AddApprovedUsersFunc: func(ids ...string) {
			if len(ids) == 0 {
				panic("no ids")
			}
		},
		RemoveApprovedUsersFunc: func(ids ...string) {
			if len(ids) == 0 {
				panic("no ids")
			}
		},
		ApprovedUsersFunc: func() []string {
			return []string{"user1", "user2"}
		},
	}
	spamFilterMock := &mocks.SpamFilterMock{
		UpdateHamFunc:  func(msg string) error { return nil },
		UpdateSpamFunc: func(msg string) error { return nil },
	}
	locatorMock := &mocks.LocatorMock{
		UserIDByNameFunc: func(userName string) int64 {
			if userName == "user1" {
				return 12345
			}
			return 0
		},
	}

	server := NewServer(Config{Detector: detectorMock, SpamFilter: spamFilterMock, Locator: locatorMock})
	ts := httptest.NewServer(server.routes(chi.NewRouter()))
	defer ts.Close()

	t.Run("check", func(t *testing.T) {
		detectorMock.ResetCalls()
		reqBody, err := json.Marshal(map[string]string{
			"msg":     "spam example",
			"user_id": "user123",
		})
		require.NoError(t, err)
		resp, err := http.Post(ts.URL+"/check", "application/json", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
		assert.Equal(t, 1, len(detectorMock.CheckCalls()))
		assert.Equal(t, "spam example", detectorMock.CheckCalls()[0].Msg)
		assert.Equal(t, "user123", detectorMock.CheckCalls()[0].UserID)
	})

	t.Run("update spam", func(t *testing.T) {
		detectorMock.ResetCalls()
		reqBody, err := json.Marshal(map[string]string{
			"msg": "test message",
		})
		require.NoError(t, err)
		resp, err := http.Post(ts.URL+"/update/spam", "application/json", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
		assert.Equal(t, 1, len(spamFilterMock.UpdateSpamCalls()))
		assert.Equal(t, "test message", spamFilterMock.UpdateSpamCalls()[0].Msg)
	})

	t.Run("update ham", func(t *testing.T) {
		detectorMock.ResetCalls()
		reqBody, err := json.Marshal(map[string]string{
			"msg": "test message",
		})
		require.NoError(t, err)
		resp, err := http.Post(ts.URL+"/update/ham", "application/json", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
		assert.Equal(t, 1, len(spamFilterMock.UpdateHamCalls()))
		assert.Equal(t, "test message", spamFilterMock.UpdateHamCalls()[0].Msg)
	})

	t.Run("add user", func(t *testing.T) {
		detectorMock.ResetCalls()
		req, err := http.NewRequest("POST", ts.URL+"/users/add", bytes.NewBuffer([]byte(`{"user_id" : "id1"}`)))
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
		assert.Equal(t, 1, len(detectorMock.AddApprovedUsersCalls()))
		assert.Equal(t, []string{"id1"}, detectorMock.AddApprovedUsersCalls()[0].Ids)
	})

	t.Run("add user by name", func(t *testing.T) {
		detectorMock.ResetCalls()
		locatorMock.ResetCalls()
		req, err := http.NewRequest("POST", ts.URL+"/users/add", bytes.NewBuffer([]byte(`{"user_name" : "user1"}`)))
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
		assert.Equal(t, 1, len(detectorMock.AddApprovedUsersCalls()))
		assert.Equal(t, []string{"12345"}, detectorMock.AddApprovedUsersCalls()[0].Ids)
		assert.Equal(t, 1, len(locatorMock.UserIDByNameCalls()))
		assert.Equal(t, "user1", locatorMock.UserIDByNameCalls()[0].UserName)
	})

	t.Run("add user by name, not found", func(t *testing.T) {
		detectorMock.ResetCalls()
		locatorMock.ResetCalls()
		req, err := http.NewRequest("POST", ts.URL+"/users/add", bytes.NewBuffer([]byte(`{"user_name" : "user2"}`)))
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, 1, len(locatorMock.UserIDByNameCalls()))
		assert.Equal(t, "user2", locatorMock.UserIDByNameCalls()[0].UserName)
	})

	t.Run("remove user", func(t *testing.T) {
		detectorMock.ResetCalls()
		locatorMock.ResetCalls()

		req, err := http.NewRequest("POST", ts.URL+"/users/delete", bytes.NewBuffer([]byte(`{"user_id" : "id1"}`)))
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
		assert.Equal(t, 1, len(detectorMock.RemoveApprovedUsersCalls()))
		assert.Equal(t, []string{"id1"}, detectorMock.RemoveApprovedUsersCalls()[0].Ids)
	})

	t.Run("remove user by name", func(t *testing.T) {
		detectorMock.ResetCalls()
		locatorMock.ResetCalls()
		req, err := http.NewRequest("POST", ts.URL+"/users/delete", bytes.NewBuffer([]byte(`{"user_name" : "user1"}`)))
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
		assert.Equal(t, 1, len(detectorMock.RemoveApprovedUsersCalls()))
		assert.Equal(t, []string{"12345"}, detectorMock.RemoveApprovedUsersCalls()[0].Ids)
		assert.Equal(t, 1, len(locatorMock.UserIDByNameCalls()))
		assert.Equal(t, "user1", locatorMock.UserIDByNameCalls()[0].UserName)
	})

	t.Run("get users", func(t *testing.T) {
		detectorMock.ResetCalls()
		resp, err := http.Get(ts.URL + "/users")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
		assert.Equal(t, 1, len(detectorMock.ApprovedUsersCalls()))
		respBody, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, `{"user_ids":["user1","user2"]}`+"\n", string(respBody))
	})
}

func TestServer_checkHandler(t *testing.T) {
	mockDetector := &mocks.DetectorMock{
		CheckFunc: func(msg string, userID string) (bool, []lib.CheckResult) {
			if msg == "spam example" {
				return true, []lib.CheckResult{{Spam: true, Name: "test", Details: "this was spam"}}
			}
			return false, []lib.CheckResult{{Details: "not spam"}}
		},
	}
	server := NewServer(Config{
		Detector: mockDetector,
		Version:  "1.0",
	})

	t.Run("spam", func(t *testing.T) {
		reqBody, err := json.Marshal(map[string]string{
			"msg":     "spam example",
			"user_id": "user123",
		})
		require.NoError(t, err)
		req, err := http.NewRequest("POST", "/check", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.checkHandler)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

		var response struct {
			Spam   bool              `json:"spam"`
			Checks []lib.CheckResult `json:"checks"`
		}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err, "error unmarshalling response")
		assert.True(t, response.Spam, "expected spam")
		assert.Equal(t, "test", response.Checks[0].Name, "unexpected check name")
		assert.Equal(t, "this was spam", response.Checks[0].Details, "unexpected check result")
	})

	t.Run("not spam", func(t *testing.T) {
		reqBody, err := json.Marshal(map[string]string{
			"msg":     "not spam example",
			"user_id": "user123",
		})
		require.NoError(t, err)
		req, err := http.NewRequest("POST", "/check", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.checkHandler)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

		var response struct {
			Spam   bool              `json:"spam"`
			Checks []lib.CheckResult `json:"checks"`
		}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err, "error unmarshalling response")
		assert.False(t, response.Spam, "expected not spam")
		assert.Equal(t, "not spam", response.Checks[0].Details, "unexpected check result")
	})

	t.Run("bad request", func(t *testing.T) {
		reqBody := []byte("bad request")
		req, err := http.NewRequest("POST", "/check", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)
		req.Body.Close()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.checkHandler)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code")
	})

}

func TestServer_updateSampleHandler(t *testing.T) {
	spamFilterMock := &mocks.SpamFilterMock{
		UpdateSpamFunc: func(msg string) error {
			if msg == "error" {
				return assert.AnError
			}
			return nil
		},
		UpdateHamFunc: func(msg string) error {
			if msg == "error" {
				return assert.AnError
			}
			return nil
		},
	}

	server := NewServer(Config{SpamFilter: spamFilterMock})

	t.Run("successful update ham", func(t *testing.T) {
		spamFilterMock.ResetCalls()
		reqBody, err := json.Marshal(map[string]string{
			"msg": "test message",
		})
		require.NoError(t, err)
		req, err := http.NewRequest("POST", "/update", bytes.NewBuffer(reqBody))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.updateSampleHandler(spamFilterMock.UpdateHam))
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
		var response struct {
			Updated bool   `json:"updated"`
			Msg     string `json:"msg"`
		}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Updated)
		assert.Equal(t, "test message", response.Msg)
		assert.Equal(t, 1, len(spamFilterMock.UpdateHamCalls()))
		assert.Equal(t, "test message", spamFilterMock.UpdateHamCalls()[0].Msg)
	})

	t.Run("update ham with error", func(t *testing.T) {
		spamFilterMock.ResetCalls()
		reqBody, err := json.Marshal(map[string]string{
			"msg": "error",
		})
		require.NoError(t, err)
		req, err := http.NewRequest("POST", "/update", bytes.NewBuffer(reqBody))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.updateSampleHandler(spamFilterMock.UpdateHam))
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code, "handler returned wrong status code")
		var response struct {
			Err     string `json:"error"`
			Details string `json:"details"`
		}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "can't update samples", response.Err)
		assert.Equal(t, "assert.AnError general error for testing", response.Details)
		assert.Equal(t, 1, len(spamFilterMock.UpdateHamCalls()))
		assert.Equal(t, "error", spamFilterMock.UpdateHamCalls()[0].Msg)
	})

	t.Run("bad request", func(t *testing.T) {
		spamFilterMock.ResetCalls()
		reqBody := []byte("bad request")
		req, err := http.NewRequest("POST", "/update", bytes.NewBuffer(reqBody))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.updateSampleHandler(spamFilterMock.UpdateHam))
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code")
	})
}

func TestServer_deleteSampleHandler(t *testing.T) {
	spamFilterMock := &mocks.SpamFilterMock{
		RemoveDynamicHamSampleFunc: func(sample string) (int, error) { return 1, nil },
		DynamicSamplesFunc: func() ([]string, []string, error) {
			return []string{"spam1", "spam2"}, []string{"ham1", "ham2"}, nil
		},
	}
	server := NewServer(Config{SpamFilter: spamFilterMock})

	t.Run("successful delete ham sample", func(t *testing.T) {
		spamFilterMock.ResetCalls()
		reqBody, err := json.Marshal(map[string]string{
			"msg": "test message",
		})
		require.NoError(t, err)
		req, err := http.NewRequest("POST", "/delete/ham", bytes.NewBuffer(reqBody))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.deleteSampleHandler(spamFilterMock.RemoveDynamicHamSample))
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
		var response struct {
			Deleted bool   `json:"deleted"`
			Msg     string `json:"msg"`
		}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Deleted)
		assert.Equal(t, "test message", response.Msg)
		require.Equal(t, 1, len(spamFilterMock.RemoveDynamicHamSampleCalls()))
		assert.Equal(t, "test message", spamFilterMock.RemoveDynamicHamSampleCalls()[0].Sample)
	})

	t.Run("delete ham sample from htmx", func(t *testing.T) {
		spamFilterMock.ResetCalls()
		req, err := http.NewRequest("POST", "/delete/ham", http.NoBody)
		require.NoError(t, err)
		req.Header.Add("HX-Request", "true") // Simulating HTMX request

		// set form htmx request, msg in r.FormValue("msg")
		req.Form = url.Values{}
		req.Form.Set("msg", "test message")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.deleteSampleHandler(spamFilterMock.RemoveDynamicHamSample))
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
		body := rr.Body.String()
		t.Log(body)
		assert.Contains(t, body, "<h4>Spam Samples</h4>", "response should contain spam samples")
		assert.Contains(t, body, "<h4>Ham Samples</h4>", "response should contain ham samples")
		require.Equal(t, 1, len(spamFilterMock.RemoveDynamicHamSampleCalls()))
		assert.Equal(t, "test message", spamFilterMock.RemoveDynamicHamSampleCalls()[0].Sample)
	})

	t.Run("delete ham sample with error", func(t *testing.T) {
		spamFilterMock.RemoveDynamicHamSampleFunc = func(sample string) (int, error) { return 0, assert.AnError }
		spamFilterMock.ResetCalls()
		reqBody, err := json.Marshal(map[string]string{
			"msg": "test message",
		})
		require.NoError(t, err)
		req, err := http.NewRequest("POST", "/delete/ham", bytes.NewBuffer(reqBody))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.deleteSampleHandler(spamFilterMock.RemoveDynamicHamSample))
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code, "handler returned wrong status code")
	})
}

func TestServer_updateApprovedUsersHandler(t *testing.T) {
	mockDetector := &mocks.DetectorMock{
		AddApprovedUsersFunc: func(ids ...string) {
			if len(ids) == 0 {
				panic("no ids")
			}
		},
		ApprovedUsersFunc: func() []string {
			return []string{"user1", "user2"}
		},
	}
	locatorMock := &mocks.LocatorMock{
		UserIDByNameFunc: func(userName string) int64 {
			if userName == "user1" {
				return 12345
			}
			return 0
		},
	}
	server := NewServer(Config{Detector: mockDetector, Locator: locatorMock})

	t.Run("successful update by name", func(t *testing.T) {
		mockDetector.ResetCalls()
		locatorMock.ResetCalls()
		req, err := http.NewRequest("POST", "/users/add", bytes.NewBuffer([]byte(`{"user_name" : "user1"}`)))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.updateApprovedUsersHandler(mockDetector.AddApprovedUsers))
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
		var response struct {
			Updated  bool   `json:"updated"`
			UserID   string `json:"user_id"`
			UserName string `json:"user_name"`
		}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Updated)
		assert.Equal(t, "12345", response.UserID)
		assert.Equal(t, "user1", response.UserName)
		assert.Equal(t, 1, len(mockDetector.AddApprovedUsersCalls()))
		assert.Equal(t, []string{"12345"}, mockDetector.AddApprovedUsersCalls()[0].Ids)
		assert.Equal(t, 1, len(locatorMock.UserIDByNameCalls()))
		assert.Equal(t, "user1", locatorMock.UserIDByNameCalls()[0].UserName)
	})

	t.Run("successful update from htmx", func(t *testing.T) {
		mockDetector.ResetCalls()
		locatorMock.ResetCalls()

		req, err := http.NewRequest("POST", "/users/add", http.NoBody)
		require.NoError(t, err)
		req.Header.Add("HX-Request", "true") // Simulating HTMX request

		req.Form = url.Values{}
		req.Form.Set("user_id", "id1")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.updateApprovedUsersHandler(mockDetector.AddApprovedUsers))
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
		body := rr.Body.String()
		t.Log(body)
		assert.Contains(t, body, "<h4>Approved Users</h4>", "response should contain approved users header")
		assert.Contains(t, body, "user1")
		assert.Contains(t, body, "user2")

	})

	t.Run("successful update by id", func(t *testing.T) {
		mockDetector.ResetCalls()
		locatorMock.ResetCalls()
		req, err := http.NewRequest("POST", "/users/add", bytes.NewBuffer([]byte(`{"user_id" : "id1"}`)))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.updateApprovedUsersHandler(mockDetector.AddApprovedUsers))
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
		var response struct {
			Updated  bool   `json:"updated"`
			UserID   string `json:"user_id"`
			UserName string `json:"user_name"`
		}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Updated)
		assert.Equal(t, "id1", response.UserID)
		assert.Equal(t, "", response.UserName)
		assert.Equal(t, 1, len(mockDetector.AddApprovedUsersCalls()))
		assert.Equal(t, []string{"id1"}, mockDetector.AddApprovedUsersCalls()[0].Ids)
		assert.Equal(t, 0, len(locatorMock.UserIDByNameCalls()))
	})
	t.Run("bad request", func(t *testing.T) {
		mockDetector.ResetCalls()
		reqBody := []byte("bad request")
		req, err := http.NewRequest("POST", "/users/add", bytes.NewBuffer(reqBody))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.updateApprovedUsersHandler(mockDetector.AddApprovedUsers))
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code")
	})
}

func TestGenerateRandomPassword(t *testing.T) {
	res1, err := GenerateRandomPassword(32)
	require.NoError(t, err)
	t.Log(res1)
	assert.Len(t, res1, 32)

	res2, err := GenerateRandomPassword(32)
	require.NoError(t, err)
	t.Log(res2)
	assert.Len(t, res2, 32)

	assert.NotEqual(t, res1, res2)
}

func TestServer_checkHandler_HTMX(t *testing.T) {
	mockDetector := &mocks.DetectorMock{
		CheckFunc: func(msg string, userID string) (bool, []lib.CheckResult) {
			return msg == "spam example", []lib.CheckResult{{Spam: msg == "spam example", Name: "test", Details: "result details"}}
		},
	}

	server := NewServer(Config{
		Detector: mockDetector,
		Version:  "1.0",
	})

	t.Run("HTMX request", func(t *testing.T) {
		form := url.Values{}
		form.Set("msg", "spam example")
		form.Set("user_id", "user123")
		req, err := http.NewRequest("POST", "/check", strings.NewReader(form.Encode()))
		require.NoError(t, err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("HX-Request", "true") // Simulating HTMX request

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.checkHandler)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

		// Check if the response contains expected HTML snippet
		assert.Contains(t, rr.Body.String(), "strong>Result:</strong> Spam detected", "response should contain spam result")
		assert.Contains(t, rr.Body.String(), "result details")
	})
}

func TestServer_htmlSpamCheckHandler(t *testing.T) {
	server := NewServer(Config{Version: "1.0"})
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", http.NoBody)
	require.NoError(t, err)

	handler := http.HandlerFunc(server.htmlSpamCheckHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler should return status OK")
	body := rr.Body.String()
	assert.Contains(t, body, "<title>Checker - TG-Spam</title>", "template should contain the correct title")
	assert.Contains(t, body, "Version: 1.0", "template should contain the correct version")
	assert.Contains(t, body, "<form", "template should contain a form")
}

func TestServer_htmlManageSamplesHandler(t *testing.T) {
	spamFilterMock := &mocks.SpamFilterMock{
		DynamicSamplesFunc: func() ([]string, []string, error) {
			return []string{"spam1", "spam2"}, []string{"ham1", "ham2"}, nil
		},
	}

	server := NewServer(Config{Version: "1.0", SpamFilter: spamFilterMock})
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/manage_samples", http.NoBody)
	require.NoError(t, err)

	handler := http.HandlerFunc(server.htmlManageSamplesHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler should return status OK")
	body := rr.Body.String()
	assert.Contains(t, body, "<title>Manage Samples - TG-Spam</title>", "template should contain the correct title")
	assert.Contains(t, body, `<div class="row" id="samples-list">`, "template should contain a samples list")
}

func TestServer_htmlManageUsersHandler(t *testing.T) {
	spamFilterMock := &mocks.SpamFilterMock{}
	detectorMock := &mocks.DetectorMock{
		ApprovedUsersFunc: func() []string {
			return []string{"user1", "user2"}
		},
	}

	server := NewServer(Config{Version: "1.0", SpamFilter: spamFilterMock, Detector: detectorMock})
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/manage_users", http.NoBody)
	require.NoError(t, err)

	handler := http.HandlerFunc(server.htmlManageUsersHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler should return status OK")
	body := rr.Body.String()
	assert.Contains(t, body, "<title>Manage Users - TG-Spam</title>", "template should contain the correct title")
	assert.Contains(t, body, "<h4>Approved Users</h4>", "template should contain users list")
}

func TestServer_stylesHandler(t *testing.T) {
	server := NewServer(Config{Version: "1.0"})
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/style.css", http.NoBody)
	require.NoError(t, err)

	handler := http.HandlerFunc(server.stylesHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler should return status OK")
	assert.Equal(t, "text/css; charset=utf-8", rr.Header().Get("Content-Type"), "handler should return CSS content type")
	assert.Contains(t, rr.Body.String(), "body", "handler should return CSS content")
}

func TestServer_getDynamicSamplesHandler(t *testing.T) {
	mockSpamFilter := &mocks.SpamFilterMock{
		DynamicSamplesFunc: func() ([]string, []string, error) {
			return []string{"spam1", "spam2"}, []string{"ham1", "ham2"}, nil
		},
	}

	server := NewServer(Config{
		SpamFilter: mockSpamFilter,
	})

	t.Run("successful response", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/samples", http.NoBody)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.getDynamicSamplesHandler)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var response struct {
			Spam []string `json:"spam"`
			Ham  []string `json:"ham"`
		}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, []string{"spam1", "spam2"}, response.Spam)
		assert.Equal(t, []string{"ham1", "ham2"}, response.Ham)
	})

	t.Run("error handling", func(t *testing.T) {
		mockSpamFilter.DynamicSamplesFunc = func() ([]string, []string, error) {
			return nil, nil, errors.New("test error")
		}

		req, err := http.NewRequest("GET", "/samples", http.NoBody)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.getDynamicSamplesHandler)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		var response struct {
			Error   string `json:"error"`
			Details string `json:"details"`
		}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "can't get dynamic samples", response.Error)
		assert.Equal(t, "test error", response.Details)
	})
}

func TestServer_reloadDynamicSamplesHandler(t *testing.T) {
	mockSpamFilter := &mocks.SpamFilterMock{
		ReloadSamplesFunc: func() error {
			return nil // Simulate successful reload
		},
	}

	server := NewServer(Config{
		SpamFilter: mockSpamFilter,
	})

	t.Run("successful reload", func(t *testing.T) {
		req, err := http.NewRequest("PUT", "/samples", http.NoBody)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.reloadDynamicSamplesHandler)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var response struct {
			Reloaded bool `json:"reloaded"`
		}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response.Reloaded)
	})

	t.Run("error during reload", func(t *testing.T) {
		mockSpamFilter.ReloadSamplesFunc = func() error {
			return errors.New("test error") // Simulate error during reload
		}

		req, err := http.NewRequest("PUT", "/samples", http.NoBody)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.reloadDynamicSamplesHandler)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		var response struct {
			Error   string `json:"error"`
			Details string `json:"details"`
		}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "can't reload samples", response.Error)
		assert.Equal(t, "test error", response.Details)
	})
}

func TestServer_reverseSamples(t *testing.T) {
	tests := []struct {
		name    string
		spam    []string
		ham     []string
		revSpam []string
		revHam  []string
	}{
		{
			name:    "Empty slices",
			spam:    []string{},
			ham:     []string{},
			revSpam: []string{},
			revHam:  []string{},
		},
		{
			name:    "Single element slices",
			spam:    []string{"a"},
			ham:     []string{"1"},
			revSpam: []string{"a"},
			revHam:  []string{"1"},
		},
		{
			name:    "Multiple elements",
			spam:    []string{"a", "b", "c"},
			ham:     []string{"1", "2", "3"},
			revSpam: []string{"c", "b", "a"},
			revHam:  []string{"3", "2", "1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{}
			gotSpam, gotHam := s.reverseSamples(tt.spam, tt.ham)
			assert.Equal(t, tt.revSpam, gotSpam)
			assert.Equal(t, tt.revHam, gotHam)
		})
	}
}

func TestServer_renderSamples(t *testing.T) {
	mockSpamFilter := &mocks.SpamFilterMock{
		DynamicSamplesFunc: func() ([]string, []string, error) {
			return []string{"spam1", "spam2"}, []string{"ham1", "ham2"}, nil
		},
	}

	server := NewServer(Config{
		SpamFilter: mockSpamFilter,
	})
	w := httptest.NewRecorder()
	server.renderSamples(w)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
	t.Log(w.Body.String())
	assert.Contains(t, w.Body.String(), "<h4>Spam Samples</h4>")
	assert.Contains(t, w.Body.String(), "spam1")
	assert.Contains(t, w.Body.String(), "spam2")
	assert.Contains(t, w.Body.String(), "<h4>Ham Samples</h4>")
	assert.Contains(t, w.Body.String(), "ham1")
	assert.Contains(t, w.Body.String(), "ham2")
}
