# Go Blog App

A production-style, server-rendered blog platform written in Go.

This project is intentionally backend-first: it demonstrates clean HTTP architecture, secure authentication flows, PostgreSQL data modeling, middleware composition, and meaningful automated testing. It is designed to be easy to read, easy to extend, and strong enough to showcase in interviews.

## Why This Project Stands Out

- Clean layered architecture: handlers, middleware, models, templates, and shared validation utilities are separated by responsibility.
- Secure-by-default HTTP setup: CSRF protection, security headers, secure cookies, session token rotation, and panic recovery are implemented.
- Real authentication lifecycle: signup, login, protected routes, logout, session-backed identity, and duplicate-email handling.
- Production-minded server config: TLS-enabled server, strict timeouts, and explicit DB connectivity checks.
- Test coverage where it matters: routing/handler behavior, middleware contracts, template helpers, and model integration tests.
- No framework lock-in: built on `net/http` plus small focused libraries.

## Core Features

- View latest blog posts on the home page.
- View an individual post by ID.
- Create new blog posts (authenticated users only).
- User signup and login with password hashing.
- User logout with session token renewal.
- Flash messages for UX feedback after key actions.
- Health endpoint for liveness checks.

## Architecture Overview

```text
HTTP Request
  -> Standard Middleware Chain
     (panic recovery -> request logging -> secure headers)
  -> Dynamic Middleware Chain
     (session load/save -> CSRF -> auth context)
  -> Route Handler
  -> Model Layer (PostgreSQL)
  -> HTML Template Rendering
  -> HTTP Response
```

### Project Structure

```text
cmd/web/                 # Application entrypoint, routes, handlers, middleware
internal/models/         # PostgreSQL data access (blogs/users) + interfaces
internal/models/mocks/   # Mock model implementations for handler tests
internal/models/testdata/# SQL setup/teardown for model integration tests
internal/validator/      # Reusable form validation and error collection
internal/assert/         # Lightweight test assertion helpers
ui/html/                 # Base template, pages, and nav partial
ui/static/               # CSS, JS, images
ui/efs.go                # Embedded FS for templates and static files
Makefile                 # Build/run/clean commands
```

## Technology Stack

- Language: Go `1.25.1`
- HTTP: `net/http`
- Router: `github.com/julienschmidt/httprouter`
- Middleware chaining: `github.com/justinas/alice`
- Templates: `html/template`
- Database: PostgreSQL + `github.com/lib/pq`
- Sessions: `github.com/alexedwards/scs/v2` + PostgreSQL store
- CSRF protection: `github.com/justinas/nosurf`
- Form decoding: `github.com/go-playground/form`
- Password hashing: `bcrypt` (`golang.org/x/crypto/bcrypt`)
- Environment loading: `github.com/joho/godotenv`

## Security Highlights

- CSRF protection on state-changing forms via `nosurf`.
- Secure and HTTP-only CSRF/session cookies.
- Security headers middleware:
  - `Content-Security-Policy`
  - `Referrer-Policy`
  - `X-Content-Type-Options`
  - `X-Frame-Options`
  - `X-XSS-Protection`
- Passwords are never stored in plain text (bcrypt hash cost 12).
- Session token renewal on login/logout to reduce fixation risk.
- Protected routes guarded by auth middleware and `Cache-Control: no-store`.

## API / Routes

### Public

- `GET /` - Home page with latest blog posts
- `GET /blog/view/:id` - View a single blog post
- `GET /user/signup` - Signup form
- `POST /user/signup` - Create user account
- `GET /user/login` - Login form
- `POST /user/login` - Authenticate user
- `GET /health` - Liveness endpoint
- `GET /static/*filepath` - Static assets

### Protected (authenticated)

- `GET /blog/create` - Blog creation form
- `POST /blog/create` - Create blog post
- `POST /user/logout` - Logout

## Data Model Snapshot

`blogs`
- `id` (identity PK)
- `title`
- `content`
- `author`
- `created_at`
- `expires`

`users`
- `id` (identity PK)
- `name`
- `email` (unique)
- `hashed_password`
- `created`

## Local Setup

### 1. Prerequisites

- Go `1.25+`
- PostgreSQL
- `air` (optional, for live reload through `make run`)

Install Air:

```bash
go install github.com/air-verse/air@latest
```

### 2. Environment Variables

Create `.env` in project root:

```env
MAIN_DSN=postgres://<user>:<password>@<host>:<port>/<db>?sslmode=disable
TEST_DSN=postgres://<user>:<password>@<host>:<port>/<test_db>?sslmode=disable
```

- `MAIN_DSN` is used by the running app.
- `TEST_DSN` is used by integration tests under `internal/models`.

### 3. Initialize Database Schema

```bash
psql "$MAIN_DSN" -f internal/models/testdata/setup.sql
```

(Optional cleanup)

```bash
psql "$MAIN_DSN" -f internal/models/testdata/teardown.sql
```

### 4. Build and Run

```bash
go mod download
make build
make run
```

App default address: `https://localhost:8000`

The server uses local TLS certs from:
- `./tls/cert.pem`
- `./tls/key.pem`

## Testing

Run all tests:

```bash
go test ./...
```

What is covered:
- Handler behavior (`cmd/web/handlers_test.go`)
- Middleware header/security behavior (`cmd/web/middleware_test.go`)
- Template date formatting helper (`cmd/web/templates_test.go`)
- Model integration checks against a real test DB (`internal/models/users_test.go`)

## Interview Talking Points

If you are presenting this project, these are the strongest engineering points to discuss:

1. Middleware-first request lifecycle and why each middleware exists.
2. Session architecture with database-backed store and token rotation.
3. Validation strategy (field + non-field errors) and server-side rendering feedback loop.
4. Use of interfaces (`BlogModelInterface`, `UserModelInterface`) to enable fast handler tests with mocks.
5. Security posture: CSRF, secure headers, password hashing, and protected-route cache controls.
6. Why server-rendered Go templates can be a pragmatic choice for internal tools or MVPs.

## Potential Next Enhancements

- Pagination for home feed.
- Edit/delete post workflows.
- Structured logging and request IDs.
- DB migrations (e.g. `golang-migrate`) for schema versioning.
- CI pipeline for test + lint automation.

---

This repository is a strong backend portfolio piece: concise codebase, clear design decisions, real security concerns addressed, and test coverage that demonstrates engineering discipline.
