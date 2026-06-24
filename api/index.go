package handler

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// ponytail: in-memory storage, swap for Supabase/Postgres when ready

type Form struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	CreatedAt time.Time  `json:"created_at"`
	Questions []Question `json:"questions,omitempty"`
}

type Question struct {
	ID      string   `json:"id"`
	Label   string   `json:"label"`
	Type    string   `json:"type"` // text, multiple_choice, rating, file
	Options []string `json:"options,omitempty"`
}

type Response struct {
	FormID      string            `json:"form_id"`
	Answers     map[string]string `json:"answers"`
	SubmittedAt time.Time         `json:"submitted_at"`
}

var (
	router *chi.Mux
	mu     sync.RWMutex
	forms  = map[string]*Form{}
	resps  = map[string][]*Response{}
)

func init() {
	router = chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		MaxAge:           300,
	}))

	router.Get("/", rootHandler)

	router.Route("/api", func(r chi.Router) {
		r.Get("/", apiInfo)
		r.Get("/health", healthHandler)

		r.Route("/forms", func(r chi.Router) {
			r.Get("/", listForms)
			r.Post("/", createForm)
			r.Get("/{id}", getForm)
			r.Put("/{id}", updateForm)
			r.Delete("/{id}", deleteForm)
		})

		r.Post("/responses/{formId}", submitResponse)
		r.Get("/analysis/{formId}", analyzeForm)
	})
}

func Handler(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}

// --- Handlers ---

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"app": "mimir-backend", "docs": "/api/health"})
}

func apiInfo(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]any{
		"endpoints": []string{
			"GET  /api/health",
			"GET  /api/forms",
			"POST /api/forms",
			"GET  /api/forms/{id}",
			"PUT  /api/forms/{id}",
			"DELETE /api/forms/{id}",
			"POST /api/responses/{formId}",
			"GET  /api/analysis/{formId}",
		},
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "app": "mimir-backend"})
}

func listForms(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()
	list := make([]*Form, 0, len(forms))
	for _, f := range forms {
		list = append(list, f)
	}
	json.NewEncoder(w).Encode(list)
}

func createForm(w http.ResponseWriter, r *http.Request) {
	var f Form
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	f.ID = randID()
	f.CreatedAt = time.Now()

	mu.Lock()
	forms[f.ID] = &f
	mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(f)
}

func getForm(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	mu.RLock()
	f, ok := forms[id]
	mu.RUnlock()
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(f)
}

func updateForm(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var f Form
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mu.Lock()
	existing, ok := forms[id]
	if !ok {
		mu.Unlock()
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	existing.Title = f.Title
	if f.Questions != nil {
		existing.Questions = f.Questions
	}
	mu.Unlock()
	json.NewEncoder(w).Encode(existing)
}

func deleteForm(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	mu.Lock()
	delete(forms, id)
	delete(resps, id)
	mu.Unlock()
	w.WriteHeader(http.StatusNoContent)
}

func submitResponse(w http.ResponseWriter, r *http.Request) {
	formID := chi.URLParam(r, "formId")
	var resp Response
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp.FormID = formID
	resp.SubmittedAt = time.Now()

	mu.Lock()
	resps[formID] = append(resps[formID], &resp)
	count := len(resps[formID])
	mu.Unlock()

	json.NewEncoder(w).Encode(map[string]any{"submitted": true, "total_responses": count})
}

func analyzeForm(w http.ResponseWriter, r *http.Request) {
	formID := chi.URLParam(r, "formId")
	mu.RLock()
	form, hasForm := forms[formID]
	responses := resps[formID]
	mu.RUnlock()

	if !hasForm {
		http.Error(w, "form not found", http.StatusNotFound)
		return
	}

	// ponytail: dummy AI analysis, replace with actual LLM call
	json.NewEncoder(w).Encode(map[string]any{
		"form_id":        formID,
		"form_title":     form.Title,
		"total_responses": len(responses),
		"summary":        "AI analysis placeholder — connect Gemini/Claude API here",
	})
}

func randID() string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 12)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
