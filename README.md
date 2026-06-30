# ecommerce

A small online store backend in Go: product catalogue, per-user shopping cart with session based auth, and coupon-based pricing with correct decimal money math.

## Prerequisites

- Go 1.22+

## Run the server

```bash
go run ./cmd/server
```

The server starts on `http://localhost:8080`. Sample products and coupons are loaded automatically on startup.

## Demo

With the server running, execute the full user journey (signup → login → browse → add to cart → apply coupon):

```bash
go run ./cmd/demo
```

Optionally override the server address:

```bash
BASE_URL=http://localhost:9090 go run ./cmd/demo
```

## Run tests

```bash
go test ./...
```

Integration tests only:

```bash
go test ./integration/...
```

With verbose output:

```bash
go test -v ./integration/...
```

## API

All routes except `/signup` and `/login` require `Authorization: Bearer <token>`.

| Method | Path | Description |
|--------|------|-------------|
| POST | `/signup` | Register a new user |
| POST | `/login` | Get a session token |
| GET | `/products?page=1&page_size=10` | List products (paginated) |
| GET | `/products/{id}` | Get product details |
| GET | `/products/{id}/inventory` | Get stock level |
| POST | `/cart/items` | Add item to cart |
| PATCH | `/cart/items/{productID}` | Update item quantity (0 removes it) |
| DELETE | `/cart/items/{productID}` | Remove item |
| GET | `/cart` | Get cart total |
| GET | `/cart?coupon=CODE` | Get cart total with coupon applied |

### Sample data loaded on startup

**Products**

| Name | Price |
|------|-------|
| Laptop | 799.99 |
| Phone | 499.99 |
| Headphones | 199.99 |

**Coupons**

| Code | Discount | Max discount |
|------|----------|--------------|
| `SAVE10` | 10% | 500.00 |
| `WELCOME20` | 20% | 200.00 |
