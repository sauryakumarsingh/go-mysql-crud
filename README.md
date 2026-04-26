# Go MySQL CRUD API — Product Catalog

A focused REST API demonstrating **real MySQL CRUD** in Go with zero ORM — raw `database/sql`, connection pooling, parameterized queries, and dynamic filtering.

## Project Structure

```
go-mysql-crud/
├── cmd/api/
│   └── main.go                     # Entry point
├── internal/
│   ├── db/
│   │   └── db.go                   # MySQL connection pool
│   ├── models/
│   │   └── models.go               # Structs & DTOs
│   ├── repository/
│   │   └── product_repo.go         # ALL SQL lives here
│   ├── handlers/
│   │   └── product_handler.go      # HTTP logic
│   └── router/
│       └── router.go               # Routes + middleware
├── migrations/
│   └── 001_create_products.sql     # DB schema + seed data
├── api.http                        # Test all endpoints
├── .env                            # DB credentials
├── Dockerfile
└── docker-compose.yml              # MySQL + API together
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /health | Health + DB ping |
| POST | /api/v1/products | Create product |
| GET | /api/v1/products | List all (with filters) |
| GET | /api/v1/products?category=X | Filter by category |
| GET | /api/v1/products?search=X | Search name/description |
| GET | /api/v1/products/stats | Aggregate stats |
| GET | /api/v1/products/{id} | Get one product |
| PUT | /api/v1/products/{id} | Partial update |
| DELETE | /api/v1/products/{id} | Delete product |

## Quick Start

### Option A — Local MySQL

1. Start MySQL locally (XAMPP / MySQL Workbench / Homebrew)
2. Run the migration:
```bash
mysql -u root -p < migrations/001_create_products.sql
```
3. Update `.env` with your credentials
4. Run:
```bash
go mod tidy
go run ./cmd/api
```

### Option B — Docker (MySQL + API together)

```bash
docker-compose up --build
# MySQL starts first, API waits for it to be healthy
```

## Key Concepts Demonstrated

- **`database/sql`** — Go's standard DB interface, driver-agnostic
- **Connection Pool** — `SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime`
- **Parameterized queries** — `?` placeholders prevent SQL injection
- **`sql.ErrNoRows`** — correctly handled for 404 responses
- **`rows.Err()`** — checked after iteration (common interview gotcha)
- **`result.RowsAffected()`** — used on UPDATE/DELETE to confirm action
- **Dynamic query building** — filter by category and/or search term
- **Partial updates** — PUT applies only non-zero fields
- **Aggregate SQL** — COUNT, SUM, AVG, COUNT DISTINCT in one query
