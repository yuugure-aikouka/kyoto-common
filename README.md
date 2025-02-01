# kyoto-common

## Getting started

### Prerequisites

- Go 1.23.5+
- Postgres
- [Air](https://github.com/air-verse/air) (for hot-reloading in development)
- [Goose](https://github.com/pressly/goose) (for database migration purposes)

### Setup

1. Clone the repository.
2. Install the dependencies:
    ```sh
    go mod tidy
    ```
3. Create an `.env` file from `.env.example` in the root folder, and fill it with its appropriate values.
4. Make sure to have a postgres database running and run the migration script:
    ```sh
    make migrate_up
    ```
5. Run development server:
    ```sh
    make dev
    ```
6. You can check for other available commands in `Makefile`.

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
