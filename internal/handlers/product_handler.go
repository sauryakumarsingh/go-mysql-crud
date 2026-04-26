package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/yourname/go-mysql-crud/internal/models"
	"github.com/yourname/go-mysql-crud/internal/repository"
)

type ProductHandler struct {
	repo *repository.ProductRepository
}

func NewProductHandler(repo *repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

// ─── POST /api/v1/products ─────────────────────────────────────────────────

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, models.Fail("invalid request body"))
		return
	}

	// Validate
	if req.Name == "" {
		writeJSON(w, http.StatusBadRequest, models.Fail("name is required"))
		return
	}
	if req.Price < 0 {
		writeJSON(w, http.StatusBadRequest, models.Fail("price cannot be negative"))
		return
	}
	if req.Stock < 0 {
		writeJSON(w, http.StatusBadRequest, models.Fail("stock cannot be negative"))
		return
	}

	product, err := h.repo.Create(req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.Fail("failed to create product: "+err.Error()))
		return
	}

	writeJSON(w, http.StatusCreated, models.OK(product, "product created successfully"))
}

// ─── GET /api/v1/products ──────────────────────────────────────────────────
// Query params: ?category=electronics  ?search=phone

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := repository.ListFilter{
		Category: r.URL.Query().Get("category"),
		Search:   r.URL.Query().Get("search"),
	}

	products, err := h.repo.GetAll(filter)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.Fail("failed to fetch products"))
		return
	}

	// Return empty array, never null
	if products == nil {
		products = []*models.Product{}
	}

	writeJSON(w, http.StatusOK, models.OK(products, ""))
}

// ─── GET /api/v1/products/{id} ─────────────────────────────────────────────

func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	if id == "" {
		writeJSON(w, http.StatusBadRequest, models.Fail("missing product id"))
		return
	}

	product, err := h.repo.GetByID(id)
	if errors.Is(err, repository.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, models.Fail("product not found"))
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.Fail("failed to fetch product"))
		return
	}

	writeJSON(w, http.StatusOK, models.OK(product, ""))
}

// ─── PUT /api/v1/products/{id} ─────────────────────────────────────────────

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	if id == "" {
		writeJSON(w, http.StatusBadRequest, models.Fail("missing product id"))
		return
	}

	var req models.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, models.Fail("invalid request body"))
		return
	}

	product, err := h.repo.Update(id, req)
	if errors.Is(err, repository.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, models.Fail("product not found"))
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.Fail("failed to update product"))
		return
	}

	writeJSON(w, http.StatusOK, models.OK(product, "product updated successfully"))
}

// ─── DELETE /api/v1/products/{id} ──────────────────────────────────────────

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	if id == "" {
		writeJSON(w, http.StatusBadRequest, models.Fail("missing product id"))
		return
	}

	err := h.repo.Delete(id)
	if errors.Is(err, repository.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, models.Fail("product not found"))
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.Fail("failed to delete product"))
		return
	}

	writeJSON(w, http.StatusOK, models.OK(nil, "product deleted successfully"))
}

// ─── GET /api/v1/products/stats ────────────────────────────────────────────

func (h *ProductHandler) Stats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.repo.GetStats()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.Fail("failed to fetch stats"))
		return
	}
	writeJSON(w, http.StatusOK, models.OK(stats, ""))
}

// ─── Helpers ───────────────────────────────────────────────────────────────

func extractID(path string) string {
	parts := strings.Split(strings.TrimSuffix(path, "/"), "/")
	last := parts[len(parts)-1]
	// avoid returning "products" or "stats" as an id
	if last == "products" || last == "stats" {
		return ""
	}
	return last
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
