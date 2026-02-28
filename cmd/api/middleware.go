package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
)

// recoverPanic ensures that in the case of a panic, a Connection header of
// 'close' is sent to the client.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			pv := recover()
			if pv != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%v", pv))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// rateLimit uses a client's IP address to limit their rate.
func (app *application) rateLimit(next http.Handler) http.Handler {
	if !app.config.limiter.enabled {
		return next
	}

	// client holds a rate limiter and the last seen time.
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// goroutine for removing stale entries from the list of clients
	go func() {
		time.Sleep(time.Minute)

		mu.Lock() // ensure no concurrency conflicts
		for ip, client := range clients {
			if time.Since(client.lastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get client ip address
		ip := realip.FromRequest(r)

		mu.Lock()

		// add client if not in the list already
		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst)}
		}

		// update client's last seen time
		clients[ip].lastSeen = time.Now()

		// if client exceeds the rate limit, send an error response
		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			app.rateLimitExceededResponse(w, r)
			return
		}

		mu.Unlock()

		// if client is not rate-limited, call next handler
		next.ServeHTTP(w, r)
	})
}

// enableCORS configures browser CORS by reflecting a request's origin if they are in
// the list of trusted origins configured on server start. It also handles CORS
// preflight requests appropriately.
func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// indicator for caches that these responses may vary
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")

		// retrieve the Origin header of the request
		origin := r.Header.Get("Origin")

		if origin != "" {
			// loop through every configured trusted origin
			for i := range app.config.cors.trustedOrigins {
				if origin == app.config.cors.trustedOrigins[i] { // on match
					// set that origin for Access-Control-Allow-Origin,
					// allowing cross-origin requests
					w.Header().Set("Access-Control-Allow-Origin", origin)

					// handle CORS preflight requests by checking OPTIONS and Access-Control-Request-Method
					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
						// set the non-CORS-safe HTTP methods
						w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
						w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

						// write 200 instead of 204 for browser compatibility
						w.WriteHeader(http.StatusOK)
						return
					}

					break
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}
