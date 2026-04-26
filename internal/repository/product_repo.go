package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yourname/go-mysql-crud/internal/models"
)

// Sentinel errors — handlers check these to decide HTTP status
var (
	ErrNotFound = errors.New("product not found")
)

// ProductRepository holds the DB connection pool.
// All SQL lives here — handlers never touch SQL directly.
type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// ─── CREATE ────────────────────────────────────────────────────────────────

func (r *ProductRepository) Create(req models.CreateProductRequest) (*models.Product, error) {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO products (id, name, description, price, stock, category, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		id, req.Name, req.Description, req.Price, req.Stock, req.Category, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("insert product: %w", err)
	}

	return &models.Product{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// ─── READ ALL (with optional filters) ─────────────────────────────────────

type ListFilter struct {
	Category string // filter by category if non-empty
	Search   string // search in name/description if non-empty
}

func (r *ProductRepository) GetAll(f ListFilter) ([]*models.Product, error) {
	// Build query dynamically based on filters
	base := `SELECT id, name, description, price, stock, category, created_at, updated_at FROM products`
	var conditions []string
	var args []interface{}

	if f.Category != "" {
		conditions = append(conditions, "category = ?")
		args = append(args, f.Category)
	}
	if f.Search != "" {
		conditions = append(conditions, "(name LIKE ? OR description LIKE ?)")
		like := "%" + f.Search + "%"
		args = append(args, like, like)
	}

	if len(conditions) > 0 {
		base += " WHERE " + strings.Join(conditions, " AND ")
	}
	base += " ORDER BY created_at DESC"

	rows, err := r.db.Query(base, args...)
	if err != nil {
		return nil, fmt.Errorf("query products: %w", err)
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		p := &models.Product{}
		err := rows.Scan(
			&p.ID, &p.Name, &p.Description,
			&p.Price, &p.Stock, &p.Category,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}
		products = append(products, p)
	}

	// Always check rows.Err() after iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return products, nil
}

// ─── READ ONE ──────────────────────────────────────────────────────────────

func (r *ProductRepository) GetByID(id string) (*models.Product, error) {
	query := `
		SELECT id, name, description, price, stock, category, created_at, updated_at
		FROM products WHERE id = ?`

	p := &models.Product{}
	err := r.db.QueryRow(query, id).Scan(
		&p.ID, &p.Name, &p.Description,
		&p.Price, &p.Stock, &p.Category,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get product by id: %w", err)
	}
	return p, nil
}

// ─── UPDATE (partial — only non-zero fields) ───────────────────────────────

func (r *ProductRepository) Update(id string, req models.UpdateProductRequest) (*models.Product, error) {
	// First confirm it exists
	existing, err := r.GetByID(id)
	if err != nil {
		return nil, err // preserves ErrNotFound
	}

	// Apply partial updates
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.Price > 0 {
		existing.Price = req.Price
	}
	if req.Stock >= 0 {
		existing.Stock = req.Stock
	}
	if req.Category != "" {
		existing.Category = req.Category
	}
	existing.UpdatedAt = time.Now()

	query := `
		UPDATE products
		SET name=?, description=?, price=?, stock=?, category=?, updated_at=?
		WHERE id=?`

	result, err := r.db.Exec(query,
		existing.Name, existing.Description,
		existing.Price, existing.Stock,
		existing.Category, existing.UpdatedAt,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("update product: %w", err)
	}

	// Sanity check: rows affected
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, ErrNotFound
	}

	return existing, nil
}

// ─── DELETE ────────────────────────────────────────────────────────────────

func (r *ProductRepository) Delete(id string) error {
	result, err := r.db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// ─── STATS (bonus — shows aggregate SQL) ──────────────────────────────────

type Stats struct {
	TotalProducts int     `json:"total_products"`
	TotalStock    int     `json:"total_stock"`
	AvgPrice      float64 `json:"avg_price"`
	Categories    int     `json:"categories"`
}

func (r *ProductRepository) GetStats() (*Stats, error) {
	query := `
		SELECT
			COUNT(*)                    AS total_products,
			COALESCE(SUM(stock), 0)    AS total_stock,
			COALESCE(AVG(price), 0)    AS avg_price,
			COUNT(DISTINCT category)   AS categories
		FROM products`

	s := &Stats{}
	err := r.db.QueryRow(query).Scan(
		&s.TotalProducts, &s.TotalStock, &s.AvgPrice, &s.Categories,
	)
	if err != nil {
		return nil, fmt.Errorf("get stats: %w", err)
	}
	return s, nil
}
