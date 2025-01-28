# kyoto-common

## Getting started

### Prerequisites

- Go 1.23.5+
- [Air](https://github.com/air-verse/air) (for hot-reloading in development)

### Setup

1. Clone the repository.
2. Install the dependencies:
    ```sh
    go mod tidy
    ```
3. Create an `.env` file from `.env.example` in the root folder, and fill it with its appropriate values.
4. Run development server:
    ```sh
    make dev
    ```
5. You can check for other available commands in `Makefile`.
