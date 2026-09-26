# App Errors

This directory acts as the central registry for all Domain-Level application errors.

## 🏛️ Architecture Rules

1. **No Magic Strings:** Do not return `errors.New("user not found")` directly in your services. 
2. **Typed Variables:** Define exported error variables here, such as `ErrUserNotFound` or `ErrInvalidCredentials`.
3. **Decoupling:** 
   - Services return these typed `apperrors`.
   - Handlers check them using `errors.Is(err, apperrors.ErrUserNotFound)` and map them to HTTP Status Codes (e.g., 404).
   - This keeps HTTP logic entirely out of the Service layer!
