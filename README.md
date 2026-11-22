# Authentication and Authorization API

A Go REST API built with Gin and MongoDB. It supports user signup and login with JWTs, plus authenticated endpoints for retrieving user records.

## Requirements

- Go 1.26 or newer
- Docker
- MongoDB (the examples below use MongoDB 7.0)

## Run MongoDB

Start a local MongoDB container:

```sh
docker run -d --name mongodb -p 27017:27017 -v mongodb_data:/data/db mongo:7.0
```

The application connects to `mongodb://localhost:27017` and uses the `userdb` database with the `users` collection. Check that MongoDB is running with `docker ps`.

## Run the API

From the repository root, download dependencies and start the server:

```sh
go mod tidy
go run .
```

The API listens on `http://localhost:8080`. A JWT signing key is generated at startup, so tokens become invalid when the server restarts.

## Endpoints

| Method | Path | Access | Purpose |
| --- | --- | --- | --- |
| `POST` | `/signup` | Public | Create a user and return a success message. |
| `POST` | `/login` | Public | Authenticate a user and return access and refresh tokens. |
| `GET` | `/users` | Bearer token, `ADMIN` role | List users. |
| `GET` | `/user/:id` | Bearer token; own record or `ADMIN` | Retrieve a user by user ID. |

The `/users` and `/user/:id` endpoints require an auth bearer token. Get the token from the `token` field in the `/login` response, then send it as `Authorization: Bearer <token>`. `/users` requires the token to belong to an `ADMIN`; `/user/:id` allows the user's own record or an `ADMIN`. User roles accepted by the model are `ADMIN` and `USER`.

## Request examples

### Sign up

```sh
curl --location 'http://localhost:8080/signup' \
--data-raw '{
    "first_name": "joshi",
    "last_name": "shubham",
    "email": "shubham.jsh@gmail.com",
    "phone": "7753073467",
    "password": "password123",
    "role": "ADMIN"
}'
```

### Log in

```sh
curl --location 'http://localhost:8080/login' \
--data-raw '{
    "email": "shubham.jsh@gmail.com",
    "password": "password123"
}'
```

### List users (ADMIN)

```sh
curl --location 'http://localhost:8080/users' \
--header 'Authorization: Bearer <token from /login response>' \
--data ''
```

### Get a user by ID

```sh
curl --location 'http://localhost:8080/user/6b02d6f3e44e4d1c891560b4' \
--header 'Authorization: Bearer <token from /login response>' \
--data ''
```

## Project layout

- `main.go` starts the API and initializes the JWT key.
- `routes/` registers HTTP routes.
- `controllers/` contains signup, login, and user handlers.
- `middleware/` validates bearer tokens for protected routes.
- `models/` defines the user document shape.
- `config/` configures MongoDB and JWT key generation.
- `helpers/` contains password hashing and JWT helpers.
