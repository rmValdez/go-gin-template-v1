# Middleware

This directory contains reusable Gin middleware functions.

## 🏛️ Architecture Rules

1. **Request Interception:** Middleware sits between the router and the handlers. It is used for cross-cutting concerns like:
   - Logging requests
   - Global Error Recovery (Catching panics)
   - Authentication (Extracting and verifying JWTs or Cookies)
   - CORS headers
   - Request ID generation
2. **Context Passing:** Use `c.Set("key", value)` to pass data (like the authenticated User ID) down to the handler.
3. **Early Return:** Use `c.Abort()` or `c.AbortWithStatusJSON()` if a request should be rejected early (e.g., unauthorized) so it never reaches the handler.
