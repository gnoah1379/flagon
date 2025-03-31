# Project Structure

This document provides an overview of the project structure for the Flagon repository.

## Directories

### `api/`
- Contains API-related code.
  - `v1/`: Version 1 of the API.

### `cmd/`
- Command-line interface.

### `docs/`
- Documentation files.

### `migrations/`
- Database migration files and migrate logic.
  - `postgres/`: PostgreSQL migration files.
  - `sqlite/`: SQLite migration files.

### `model/`
- Contains data models for the application.

### `pkg/`
- Shared utility packages.
  - `config/`: Configuration management.
  - `database/`: Database-related utilities.
  - `log/`: Logging utilities.
  - `cache/`: Cache utilities.

### `repository/`
- Data access layer for interacting with the database.

### `server/`
- Transport layer.

### `service/`
- Business logic and service layer.

### `ui/`
- User interface code.