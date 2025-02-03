# kyoto-common

## Getting started

### Prerequisites

- Go 1.23.5+
- Postgres
- [Air](https://github.com/air-verse/air) (for hot-reloading in development)
- [Goose](https://github.com/pressly/goose) (for database migration purposes)
- [sqlc](https://sqlc.dev/) (for generating type-safe DAO code from SQL)

### Setup

1. Clone the repository.
2. Install the dependencies:
    ```sh
    go mod tidy
    ```
3. Create an `.env` file from `.env.example` in the root folder, and fill it with its appropriate values.
4. Run `docker compose up` to run all app dependencies (postgres, etc.)
5. Run the migration script:
    ```sh
    make migrate_up
    ```
6. Run development server:
    ```sh
    make dev
    ```
7. You can check for other available commands in `Makefile`.

## Guide for developers

### Creating a new migration

Suppose we want to create a new table called `users`. Then do the following steps:

1. Run this to create the new migration:
    ```sh
    make new_migration name=create_users_table
    ```
2. A new migration file will be created in `db/migration` folder.
3. Fill in your migrate up/down logic there.
4. Run `make migrate_up` to run the migrations.

### Working with sqlc

sqlc is a tool that generates type-safe Go code from raw SQL queries. Instead of manually writing the cumbersome query logic, you define SQL statements, and sqlc creates Go functions for you.

#### How to use sqlc

To work with sqlc, follow these steps:

1. Define your queries in the `db/queries` folder. Each query must have a `-- name:` comment to specify the function name in Go.  Example (`db/queries/users.sql`):
    ```sql
    -- name: GetUserById :one
    SELECT * FROM users WHERE id = $1;
    ```

2. Run the following command to generate or sync the Go code:
    ```bash
    make sqlc
    ```
    This will process all queries in `db/queries` and generate Go functions under `db/store`, based on the `sqlc.yaml` configuration.

    📝 Note: Everything inside `db/store/` is auto-generated, except for `store.go`.

3. Now, you can use the generated functions in your API handlers under the `api/` folder:
    ```go
    func (s *Server) getFirstUserHandler(c echo.Context) error {
        id := 1
        firstUser, err := s.store.GetUserById(c.Request().Context(), id)
        if err != nil {
            return jsonResponse(c, http.StatusInternalServerError, nil)
        }

        return jsonResponse(c, http.StatusOK, firstUser)
    }
    ```
