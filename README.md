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

## Request/Response examples

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

**Response**

```json
{
    "message": "User created successfully"
}
```

### Log in

```sh
curl --location 'http://localhost:8080/login' \
--data-raw '{
    "email": "shubham.jsh@gmail.com",
    "password": "password123"
}'
```

**Response**

```json
{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiIiwiZW1haWwiOiIiLCJyb2xlIjoiIiwiZXhwIjoxNzk1MzUyNTY3fQ.MRaVdrWrUZcpm0Ri57q4DCqAVn-Isvd4o4lyhVALndA",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNmIwMmU0M2E4OGUzYmIwMDNiMzhkNTNlIiwiZW1haWwiOiJzaHViaGFAZ21haWwuY29tIiwicm9sZSI6IkFETUlOIiwiZXhwIjoxNzk1MzUxOTY3fQ.D_ValgdDFULgSajZ2BsrAj0nniwsHa34BhQU_6r_1QE",
    "user": {
        "ID": "6b02e43a88e3bb003b38d53e",
        "first_name": "joshi",
        "last_name": "shubham",
        "password": "$2a$10$8nxXqcuIs18IwDQm9p/JQObN0GyIwprZY0T2UjtVk3F5Jwqbk1Q32",
        "email": "shubha@gmail.com",
        "phone": "775307346",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNmIwMmU0M2E4OGUzYmIwMDNiMzhkNTNlIiwiZW1haWwiOiJzaHViaGFAZ21haWwuY29tIiwicm9sZSI6IkFETUlOIiwiZXhwIjoxNzk1MzUxOTEwfQ.6y14ZlbILnCTyzQpgpDFaLM-C7OL1FU5dNBz_qwh8Xg",
        "role": "ADMIN",
        "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiIiwiZW1haWwiOiIiLCJyb2xlIjoiIiwiZXhwIjoxNzk1MzUyNTEwfQ.sJz-l45b47qXUiROiykEH1oXQkIYuTWlglJ7iV6JbLw",
        "created_at": "2026-11-22T12:46:50.74Z",
        "updated_at": "2026-11-22T12:46:50.74Z",
        "user_id": "6b02e43a88e3bb003b38d53e"
    }
}
```

### List users (ADMIN)

```sh
curl --location 'http://localhost:8080/users' \
--header 'Authorization: Bearer <token from /login response>' \
--data ''
```

**Response**

```json
[
    {
        "ID": "6b02d6f3e44e4d1c891560b4",
        "first_name": "joshi",
        "last_name": "shubham",
        "password": "$2a$10$YebC2zUqdoEVrNqiVnyMwe0HrXngDaqIYPcUoAV2bta1AEM7PYD4i",
        "email": "shubham.jsh@gmail.com",
        "phone": "7753073467",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNmIwMmQ2ZjNlNDRlNGQxYzg5MTU2MGI0IiwiZW1haWwiOiJzaHViaGFtLmpzaEBnbWFpbC5jb20iLCJyb2xlIjoiQURNSU4iLCJleHAiOjE3OTUzNTA5ODF9.AewVZ9BgBvWGTpQqzPRWySdLKQmo5lXCOnQOJJ6SQlI",
        "role": "ADMIN",
        "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiIiwiZW1haWwiOiIiLCJyb2xlIjoiIiwiZXhwIjoxNzk1MzUxNTgxfQ.w2atGGQuk_izeLP0MmzvhZP2tWD3L_mIv-5PLmqs4jA",
        "created_at": "2026-11-22T11:50:11.502Z",
        "updated_at": "2026-11-22T12:31:21.268Z",
        "user_id": "6b02d6f3e44e4d1c891560b4"
    },
    {
        "ID": "6b02e41f88e3bb003b38d53d",
        "first_name": "joshi",
        "last_name": "shubham",
        "password": "$2a$10$jVYZFvZ6P8wuByAC7f39Re3Nvperarr7.cLu7UvWwAM28KQRvlcVS",
        "email": "shubham@gmail.com",
        "phone": "7753073461",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNmIwMmU0MWY4OGUzYmIwMDNiMzhkNTNkIiwiZW1haWwiOiJzaHViaGFtQGdtYWlsLmNvbSIsInJvbGUiOiJBRE1JTiIsImV4cCI6MTc5NTM1MTg4M30.tAJ9Gq6iS7tHgHwuyuQ6UIp_v_B1TvQ7zq9BNRF5hhs",
        "role": "ADMIN",
        "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiIiwiZW1haWwiOiIiLCJyb2xlIjoiIiwiZXhwIjoxNzk1MzUyNDgzfQ.vH6htRl5InWzWwMnmekHip7imJLf9f7KphEqkxw_XEA",
        "created_at": "2026-11-22T12:46:23.021Z",
        "updated_at": "2026-11-22T12:46:23.021Z",
        "user_id": "6b02e41f88e3bb003b38d53d"
    },
    {
        "ID": "6b02e43a88e3bb003b38d53e",
        "first_name": "joshi",
        "last_name": "shubham",
        "password": "$2a$10$8nxXqcuIs18IwDQm9p/JQObN0GyIwprZY0T2UjtVk3F5Jwqbk1Q32",
        "email": "shubha@gmail.com",
        "phone": "775307346",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNmIwMmU0M2E4OGUzYmIwMDNiMzhkNTNlIiwiZW1haWwiOiJzaHViaGFAZ21haWwuY29tIiwicm9sZSI6IkFETUlOIiwiZXhwIjoxNzk1MzUxOTY3fQ.D_ValgdDFULgSajZ2BsrAj0nniwsHa34BhQU_6r_1QE",
        "role": "ADMIN",
        "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiIiwiZW1haWwiOiIiLCJyb2xlIjoiIiwiZXhwIjoxNzk1MzUyNTY3fQ.MRaVdrWrUZcpm0Ri57q4DCqAVn-Isvd4o4lyhVALndA",
        "created_at": "2026-11-22T12:46:50.74Z",
        "updated_at": "2026-11-22T12:47:47.988Z",
        "user_id": "6b02e43a88e3bb003b38d53e"
    }
]
```

### Get a user by ID

```sh
curl --location 'http://localhost:8080/user/6b02d6f3e44e4d1c891560b4' \
--header 'Authorization: Bearer <token from /login response>' \
--data ''
```

**Response**

```json
{
    "ID": "6b02d6f3e44e4d1c891560b4",
    "first_name": "joshi",
    "last_name": "shubham",
    "password": "$2a$10$YebC2zUqdoEVrNqiVnyMwe0HrXngDaqIYPcUoAV2bta1AEM7PYD4i",
    "email": "shubham.jsh@gmail.com",
    "phone": "7753073467",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNmIwMmQ2ZjNlNDRlNGQxYzg5MTU2MGI0IiwiZW1haWwiOiJzaHViaGFtLmpzaEBnbWFpbC5jb20iLCJyb2xlIjoiQURNSU4iLCJleHAiOjE3OTUzNTA5ODF9.AewVZ9BgBvWGTpQqzPRWySdLKQmo5lXCOnQOJJ6SQlI",
    "role": "ADMIN",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiIiwiZW1haWwiOiIiLCJyb2xlIjoiIiwiZXhwIjoxNzk1MzUxNTgxfQ.w2atGGQuk_izeLP0MmzvhZP2tWD3L_mIv-5PLmqs4jA",
    "created_at": "2026-11-22T11:50:11.502Z",
    "updated_at": "2026-11-22T12:31:21.268Z",
    "user_id": "6b02d6f3e44e4d1c891560b4"
}
```

## Project layout

- `main.go` starts the API and initializes the JWT key.
- `routes/` registers HTTP routes.
- `controllers/` contains signup, login, and user handlers.
- `middleware/` validates bearer tokens for protected routes.
- `models/` defines the user document shape.
- `config/` configures MongoDB and JWT key generation.
- `helpers/` contains password hashing and JWT helpers.
