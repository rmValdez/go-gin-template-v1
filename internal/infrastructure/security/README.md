# Security

This directory handles cryptography, hashing, and token generation.

## 🏛️ Architecture Rules

1. **Focused Utility:** Do not use this as a generic `utils/` dumping ground. This package is strictly for security functions.
2. **Stateless:** Functions here (like `HashPassword`, `CheckPassword`, `GenerateToken`) should generally be stateless utility functions.
3. **Best Practices:** 
   - Passwords should use bcrypt or Argon2.
   - JWT generation should rely on the `config` environment variables.
