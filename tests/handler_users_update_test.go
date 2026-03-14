package tests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/tbone317/chirpy/internal/auth"
)

// makeHandler returns an http.HandlerFunc that calls the provided handler
// with a real-looking authenticated request built from the given token string.
func makeAuthRequest(method, path, authorization, body string) *http.Request {
	var bodyReader *strings.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	} else {
		bodyReader = strings.NewReader("")
	}
	req := httptest.NewRequest(method, path, bodyReader)
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	return req
}

// makeValidToken signs a JWT with the given secret for the given user.
func makeValidToken(secret string, userID uuid.UUID) string {
	token, _ := auth.MakeJWT(userID, secret, time.Hour)
	return "Bearer " + token
}

func runHandler(handler http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler(rr, req)
	return rr
}

func assertStatus(t *testing.T, rr *httptest.ResponseRecorder, want int) {
	t.Helper()
	if rr.Code != want {
		t.Fatalf("status = %d, want %d (body: %s)", rr.Code, want, rr.Body.String())
	}
}

func assertBodyContains(t *testing.T, rr *httptest.ResponseRecorder, sub string) {
	t.Helper()
	body, _ := io.ReadAll(rr.Body)
	if !strings.Contains(string(body), sub) {
		t.Fatalf("body = %q, want substring %q", string(body), sub)
	}
}

// ---------- PUT /api/users auth tests ----------

func TestUpdateUserAuth_MissingHeader(t *testing.T) {
	handler := makeUpdateUserAuthOnlyHandler("test-secret")
	req := makeAuthRequest(http.MethodPut, "/api/users", "", `{"email":"a@b.com","password":"pw"}`)
	rr := runHandler(handler, req)
	assertStatus(t, rr, http.StatusUnauthorized)
}

func TestUpdateUserAuth_InvalidFormat(t *testing.T) {
	handler := makeUpdateUserAuthOnlyHandler("test-secret")
	req := makeAuthRequest(http.MethodPut, "/api/users", "Token abc123", `{"email":"a@b.com","password":"pw"}`)
	rr := runHandler(handler, req)
	assertStatus(t, rr, http.StatusUnauthorized)
}

func TestUpdateUserAuth_InvalidJWT(t *testing.T) {
	handler := makeUpdateUserAuthOnlyHandler("test-secret")
	req := makeAuthRequest(http.MethodPut, "/api/users", "Bearer not.a.jwt", `{"email":"a@b.com","password":"pw"}`)
	rr := runHandler(handler, req)
	assertStatus(t, rr, http.StatusUnauthorized)
}

// makeUpdateUserAuthOnlyHandler returns a handler that only checks auth
// (no DB), so we can test the auth layer in isolation.
func makeUpdateUserAuthOnlyHandler(secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := auth.GetBearerToken(r.Header)
		if err != nil {
			http.Error(w, `{"error":"Couldn't find token"}`, http.StatusUnauthorized)
			return
		}
		_, err = auth.ValidateJWT(tokenStr, secret)
		if err != nil {
			http.Error(w, `{"error":"Couldn't validate token"}`, http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// ---------- DELETE /api/chirps/{chirpID} auth tests ----------

func TestDeleteChirpAuth_MissingHeader(t *testing.T) {
	handler := makeDeleteChirpAuthOnlyHandler("test-secret", uuid.New(), uuid.New())
	req := makeAuthRequest(http.MethodDelete, "/api/chirps/"+uuid.New().String(), "", "")
	rr := runHandler(handler, req)
	assertStatus(t, rr, http.StatusUnauthorized)
}

func TestDeleteChirpAuth_InvalidJWT(t *testing.T) {
	handler := makeDeleteChirpAuthOnlyHandler("test-secret", uuid.New(), uuid.New())
	req := makeAuthRequest(http.MethodDelete, "/api/chirps/"+uuid.New().String(), "Bearer not.a.jwt", "")
	rr := runHandler(handler, req)
	assertStatus(t, rr, http.StatusUnauthorized)
}

func TestDeleteChirpAuth_WrongUser(t *testing.T) {
	const secret = "test-secret"
	chirpOwner := uuid.New()
	otherUser := uuid.New()
	// Handler simulates: chirp belongs to chirpOwner, request is from otherUser.
	handler := makeDeleteChirpAuthOnlyHandler(secret, chirpOwner, otherUser)
	req := makeAuthRequest(http.MethodDelete, "/api/chirps/"+uuid.New().String(), makeValidToken(secret, otherUser), "")
	rr := runHandler(handler, req)
	assertStatus(t, rr, http.StatusForbidden)
}

func TestDeleteChirpAuth_Success(t *testing.T) {
	const secret = "test-secret"
	owner := uuid.New()
	handler := makeDeleteChirpAuthOnlyHandler(secret, owner, owner)
	req := makeAuthRequest(http.MethodDelete, "/api/chirps/"+uuid.New().String(), makeValidToken(secret, owner), "")
	rr := runHandler(handler, req)
	assertStatus(t, rr, http.StatusNoContent)
}

// makeDeleteChirpAuthOnlyHandler returns a handler that checks auth and
// ownership (no DB), simulating that the chirp belongs to chirpOwnerID.
func makeDeleteChirpAuthOnlyHandler(secret string, chirpOwnerID, requestingUserID uuid.UUID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := auth.GetBearerToken(r.Header)
		if err != nil {
			http.Error(w, `{"error":"Couldn't find token"}`, http.StatusUnauthorized)
			return
		}
		userID, err := auth.ValidateJWT(tokenStr, secret)
		if err != nil {
			http.Error(w, `{"error":"Couldn't validate token"}`, http.StatusUnauthorized)
			return
		}
		if userID != chirpOwnerID {
			http.Error(w, `{"error":"Forbidden"}`, http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
