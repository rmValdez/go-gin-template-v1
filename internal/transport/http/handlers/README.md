# Handlers (Controllers)

This directory contains the HTTP handlers (often called controllers in MVC frameworks).

## 🏛️ Architecture Rules

1. **No Business Logic:** Handlers should NEVER contain business logic or direct database queries. 
2. **Delegation:** The handler's only job is to:
   - Receive the HTTP request
   - Parse and validate the incoming JSON body into a **DTO** (Data Transfer Object)
   - Call the appropriate method in the `internal/services/` layer
   - Format the response using the `internal/response/` package and return the correct HTTP status code.
3. **Dependency Injection:** Handlers should accept their required Services via constructor functions (e.g., `NewAuthHandler(service)`), making them easy to test.

## 📄 Example Flow
```go
// 1. Parse DTO
var req LoginRequest
if err := c.ShouldBindJSON(&req); err != nil { ... }

// 2. Call Service
token, err := h.authService.Login(req.Email, req.Password)

// 3. Format Response
response.Success(c, http.StatusOK, token)
```
