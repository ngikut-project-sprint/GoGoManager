# Authentication & Authorization

## PIC

...

## Background:

Manager can register to the app and start manage their employee

## Contract:

**POST /v1/auth**

Request Body:

```json
{
  "email": "name@name.com", // should in email format
  "password": "asdfasdf", // string | minLength: 8 | maxLength: 32
  "action": "create|login" // string | enum
}
```

Response:

- `200` Ok for existing user

```json
{
  "email": "name@name.com",
  "token": "asdfasdf" // use any token you want
}
```

- `201` Created for new user

```json
{
  "email": "name@name.com",
  "token": "asdfasdf" // use any token you want
}
```

- `400` Bad Request case:
  - Validation error
- `401` Not Found case:
  - Email is not found if `action == 'loging'`
- `409` Conflict case:
  - Email is existed if `action == 'create'`
- `500` Server Error
