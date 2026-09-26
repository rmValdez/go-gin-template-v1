# Repositories (Data Access Layer)

This directory is strictly for interacting with the database.

## 🏛️ Architecture Rules

1. **Abstract the Database:** The rest of the application should not know we are using GORM or PostgreSQL. 
2. **Interface First:** Every repository must export an interface. For example, `UserRepository` should define methods like `FindByEmail()` and `Create()`.
3. **Keep it Dumb:** Repositories do not perform business logic (like checking passwords). They only execute CRUD operations and return `models`.
4. **Mockability:** Because we use interfaces here, the `services` layer can easily mock the database during unit testing without needing a real PostgreSQL connection.
