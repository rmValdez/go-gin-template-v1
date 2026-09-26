# Services (Business Logic)

This directory contains the core business logic of your application.

## 🏛️ Architecture Rules

1. **Pure Go, No HTTP:** Services should have absolutely no knowledge of HTTP, Gin contexts (`*gin.Context`), or JSON serialization. They only accept and return standard Go types and structs.
2. **Interface Driven:** Services should define an interface (e.g., `AuthService`) and a concrete implementation. This allows the Handler layer to mock the service during tests.
3. **Dependency Injection:** Services do not access the database directly. They accept a `Repository` interface via their constructor (e.g., `NewAuthService(repo)`).
4. **Domain Errors:** Services should return typed domain errors from `internal/apperrors/` (e.g., `apperrors.ErrUserNotFound`) rather than HTTP errors (like 404). The handler is responsible for mapping domain errors to HTTP codes.
