# Employees

## PIC

...

## Background:

Manager can manage their employees

## Contract:

**POST /v1/employee**

Request Header:

|      key      |   value    |
| :-----------: | :--------: |
| Authorization | bearer ... |

Request Body:

```json
{
  "identityNumber": "", // string | minlength 5 | maxlength 33
  "name": "", // string | minlength 4 | maxlength 33
  "employeeImageUri": "", // string | should be an uri
  "gender": "male|female", // string | enum
  "departmentId": "" // string | should be a valid departmentId
}
```

Response:

- `201` Created

```json
{
  "identityNumber": "",
  "name": "",
  "employeeImageUri": "",
  "gender": "",
  "departmentId": ""
}
```

- `400` Bad Request case:
  - Validation error
- `401` Unauthorized for:
  - expired / invalid / missing request token
- `409` Conflict for:
  - identity number
- `500` Server Error

**GET /v1/employee**

Request parameters (all optional)

- `limit` & `offset` limit the output of the data
  - default `limit=5&offset=0`
  - value should be a number
  - invalid `limit` / `offset` value will use the default value
- `identityNumber` filter the result based on the identity number
  - search should be a wildcard (`123` should return information like `11123333`)
  - value should be a string
- `gender` filter the result based on gender
  - value should be an enum of `male` | `female`
  - invalid `gender` will cause the search come up empty (`[]`)
- `departmentId` filter the result based on the department
  - value should be a valid `departmentId`
  - invalid `departmentId` will cause the search come up empty (`[]`)

Response:

- `200` Ok

```json
[
  {
    "identityNumber": "",
    "name": "",
    "employeeImageUri": "",
    "gender": "",
    "departmentId": ""
  }
]
```

- `401` Unauthorized for:
  - expired / invalid / missing request token
- `500` Server Error

**PATCH /v1/employee/:identityNumber**

Request Header:

|      key      |   value    |
| :-----------: | :--------: |
| Authorization | bearer ... |

Request Body:

```json
{
  "identityNumber": "", // string | minlength 5 | maxlength 33
  "name": "", // string | minlength 4 | maxlength 33
  "employeeImageUri": "", // string | should be an uri
  "gender": "male|female", // string | enum
  "departmentId": "" // string | should be a valid departmentId
}
```

Response:

- `200` Ok

```json
{
  "identityNumber": "",
  "name": "",
  "employeeImageUri": "",
  "gender": "",
  "departmentId": ""
}
```

- `400` Bad Request for:
  - Validation error
- `401` Unauthorized for:
  - expired / invalid / missing request token
- `404` Not Found for:
  - `identityNumber` is not found
- `409` Conflict for:
  - identity number
- `500` Server Error

**DELETE /v1/employee/:identityNumber**

Request Header:

|      key      |   value    |
| :-----------: | :--------: |
| Authorization | bearer ... |

Response:

- `200` Ok deleted
- `401` Unauthorized for:
  - expired / invalid / missing request token
- `404` Not Found for:
  - `identityNumber` is not found
- `500` Server Error
