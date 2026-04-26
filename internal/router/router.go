package router

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/yourname/go-mysql-crud/internal/handlers"
	"github.com/yourname/go-mysql-crud/internal/repository"
)

func New(db *sql.DB) http.Handler {
	repo := repository.NewProductRepository(db)
	ph := handlers.NewProductHandler(repo)

	mux := http.NewServeMux()

	// ── Health ──────────────────────────────────────────────────────────────
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		// Also ping the DB on health check — very useful in production
		dbStatus := "ok"
		if err := db.Ping(); err != nil {
			dbStatus = "unreachable"
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":   "ok",
			"database": dbStatus,
			"version":  "1.0.0",
		})
	})

	// ── Product routes ───────────────────────────────────────────────────────
	mux.HandleFunc("GET /api/v1/products/stats", ph.Stats)   // must be before {id} route
	mux.HandleFunc("GET /api/v1/products", ph.List)
	mux.HandleFunc("POST /api/v1/products", ph.Create)
	mux.HandleFunc("GET /api/v1/products/{id}", ph.Get)
	mux.HandleFunc("PUT /api/v1/products/{id}", ph.Update)
	mux.HandleFunc("DELETE /api/v1/products/{id}", ph.Delete)

	return logger(cors(mux))
}

// ── Inline middleware (small project — no separate package needed) ──────────

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		println(r.Method, r.URL.Path, time.Since(start).String())
	})
}
