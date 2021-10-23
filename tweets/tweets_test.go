package tweets

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShouldReturn404SinceTweetDoesNotExist(t *testing.T) {
	w := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "http://localhost:4000/users/1/")
}
