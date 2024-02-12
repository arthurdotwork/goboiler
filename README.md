<center>
<img src="https://github.com/MariaLetta/free-gophers-pack/blob/master/characters/png/47.png?raw=true" width="300px">

# Goboiler
</center>

Goboiler is my personal golang boilerplate for services & APIs.
It uses the following packages:
- [sqlx](https://github.com/jackc/sqlx) and [sqalx](https://github.com/heetch/sqalx) for database access.
- [gin](https://github.com/gin-gonic/gin) for the http router.
- [zerolog](https://github.com/rs/zerolog) for logging.
- [tern](https://github.com/jackc/tern) for database migrations.
- [task](https://taskfile.dev) as a makefile replacement.

## Getting started

### Prerequisites

- [Go](https://golang.org/dl/)
- [Docker](https://www.docker.com/products/docker-desktop)
- [Task](https://taskfile.dev/#/installation)

### Installation

1. Clone the repository
```sh
git clone github.com/arthureichelberger/goboiler
```

2. Install dependencies
```sh
task install
```

3. Start the associated services
```sh
docker-compose up -d
```

4. Run the migrations
```sh
task migrate:fresh
```

5. Run the server
```sh
task run
```

## Usage

### Run the server
```sh
task run
```

### Run the tests
```sh
task test # also resets the database
```

### Run the linter
```sh
task lint
```

### Run the migrations
```sh
task migrate:up # to run the next migrations
task migrate:down # to revert the database
task migrate:fresh # to revert the database and run all the migrations
```

### Create a new migration
```sh
task migrate:new -- {{migration_name}} # ex: create_users_table
```
