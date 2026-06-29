# ecommerce

## Running Unit Tests

```bash
go test ./auth/...
```

To run with verbose output:

```bash
go test -v ./auth/...
```

## Running Integration Tests

The integration tests spin up a real HTTP server and exercise the signup, login, and authenticate flows end-to-end.

```bash
go test ./integration/...
```

To run with verbose output:

```bash
go test -v ./integration/...
```
