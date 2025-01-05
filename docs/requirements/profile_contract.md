# Profile Management

## PIC

...

## Background:

Manager can manage their information

## Contract:

**GET /v1/user**

Request Header:

|      key      |   value    |
| :-----------: | :--------: |
| Authorization | bearer ... |

Response:

- `200` Ok

```json
{
  "email": "name@name.com",
  "name": "",
  "userImageUri": "",
  "companyName": "",
  "companyImageUri": ""
}
```

- `401` Unauthorized for:
  - expired / invalid / missing request token
- `500` Server Error

**PATCH /v1/user**

Request Header:

|      key      |   value    |
| :-----------: | :--------: |
| Authorization | bearer ... |

Request Body:

```json
{
  "email": "name@name.com", // should in email format
  "name": "", // string | minLength 4 | maxLength 52
  "userImageUri": "", // string | should be an URI
  "companyName": "", // string | minLength 4 | maxLength 52
  "companyImageUri": "" // string | should be an URI
}
```

Response:

- `200` Ok

```json
{
  "email": "name@name.com",
  "name": "",
  "userImageUri": "",
  "companyName": "",
  "companyImageUri": ""
}
```

- `400` Bad Request for:
  - Validation error
- `401` Unauthorized for:
  - expired / invalid / missing request token
- `409` Conflict for:
  - Email is used by another person
- `500` Server Error
