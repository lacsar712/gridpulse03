package relay

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// stallHandler mimics an unresponsive station endpoint. It blocks until the
// test releases it (quit) or its request context is done, so it never wedges
// the test server's Close regardless of disconnect-propagation quirks.
func stallHandler(quit <-chan struct{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-quit:
		case <-r.Context().Done():
		}
	}
}

// PostOutbound must honor a canceled context and return promptly, not drag the
// call out to the http.Client timeout. Asserts only the observable wall-clock
// behavior required by the business symptom.
func TestPostOutboundCancelEndsEarly(t *testing.T) {
	quit := make(chan struct{})
	srv := httptest.NewServer(stallHandler(quit))
	defer func() {
		close(quit) // release the in-flight handler first...
		srv.Close() // ...so Close returns promptly.
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	start := time.Now()
	done := make(chan error, 1)
	go func() { done <- PostOutbound(ctx, srv.URL, []byte("Trip=TG01-12")) }()

	// Let the request reach the server, then cancel the context bound to it.
	time.Sleep(150 * time.Millisecond)
	cancel()

	// Well under the 3s http.Client timeout PostOutbound otherwise uses.
	const cancelBudget = 1500 * time.Millisecond
	select {
	case <-done:
		if elapsed := time.Since(start); elapsed > cancelBudget {
			t.Fatalf("cancel did not end PostOutbound promptly: %v (want < %v)", elapsed, cancelBudget)
		}
	case <-time.After(cancelBudget):
		t.Fatalf("PostOutbound did not return after cancel within %v", cancelBudget)
	}
}
