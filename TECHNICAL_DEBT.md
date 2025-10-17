# Technical Debt Checklist

This document tracks the identified technical debt in the codebase.

## Auction Service

- [ ] **Persistence Layer**: Replace mock data and placeholder logic in `services/auction.go` with a real persistence layer (e.g., database or Firebase integration) for all service methods (`CreateAuction`, `GetAuction`, etc.).
- [ ] **Secure ID Generation**: Refactor the ID generation helper functions in `services/auction.go` to use the `crypto/rand` package for generating cryptographically secure, random identifiers instead of relying on `time.Now().UnixNano()`.
- [ ] **Robust Error Handling**: In `internal/api/handlers/auction.go`, add proper error handling for the `strconv.Atoi` conversion when parsing `page` and `limit` query parameters. Return a `400 Bad Request` for invalid input.
- [ ] **Safe Type Assertions**: In `internal/api/handlers/auction.go`, refactor the `userID` retrieval to use the `value, ok := c.Get("userID").(string)` idiom to safely handle cases where the user ID is missing or not a string, preventing potential panics.
- [ ] **Configuration Security**: Remove default secrets (e.g., `JWT_SECRET`, `DATABASE_URL` password) from `internal/config/config.go`. The application should fail fast on startup if critical secrets are not provided via the environment, especially in production.

## Auth Service

- [ ] **Unified Persistence for Auth**: Refactor the `betterauth` client (`pkg/betterauth/client.go`) to use the persistent database via the `Database` interface, removing its internal in-memory user map. This will unify data storage and ensure user data persists across restarts.
- [ ] **Dynamic Role-Based Access Control (RBAC)**: In `internal/middleware/auth.go`, replace the hardcoded user role in `RoleMiddleware` with a dynamic lookup from the database to enable proper role-based access control.
- [ ] **Production Secret Management**: Remove the hardcoded `JWT_SECRET` from `docker-compose.yml`. In a production environment, this should be injected securely using a secrets management system (e.g., Docker secrets, environment variables from a CI/CD pipeline) instead of being committed to the repository.