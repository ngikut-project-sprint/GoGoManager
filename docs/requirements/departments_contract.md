# Departments

## PIC

...

## Background:

Manager can manage department for their employee

## Contract:

**POST /v1/department**

Request Header:

|      key      |   value    |
| :-----------: | :--------: |
| Authorization | bearer ... |

Request Body:

```js
{
  "name": "" // string | minlength 4 | maxlength 33
}
```

Response:

- `201` Created

```js
{
  "departmentId": "", // use any id you want
  "name": ""
}
```

- `400` Bad Request for:
  - Validation error
- `401` Unauthorized for:
  - expired / invalid / missing request token
- `500` Server Error

**GET /v1/department**

Request parameters (all optional)

- `limit` & `offset` limit the output of the data
  - default `limit=5&offset=0`
  - value should be a number
  - invalid `limit` / `offset` value will use the default value
- `name` filter the result based on name
  - search should be a wildcard (`123` should return information like `11123333`)
  - value should be a string

Response:

- `200` Ok

```js
[
  {
    departmentId: "",
    name: "",
  },
];
```

- `401` Unauthorized for:
  - expired / invalid / missing request token
- `500` Server Error

**PATCH /v1/department/:departmentId**

Request Header:

|      key      |   value    |
| :-----------: | :--------: |
| Authorization | bearer ... |

Request Body:

```js
{
  "name": "" // string | minlength 4 | maxlength 33
}
```

Response:

- `200` Ok

```js
{
  "departmentId": "", // use any id you want
  "name": ""
}
```

- `400` Bad Request for:
  - Validation error
- `401` Unauthorized for:
  - expired / invalid / missing request token
- `404` Not Found for:
  - `departmentId` is not found
- `500` Server Error

**DELETE /v1/department/:departmentId**

Request Header:

|      key      |   value    |
| :-----------: | :--------: |
| Authorization | bearer ... |

Response:

- `200` Ok deleted
- `401` Unauthorized for:
  - expired / invalid / missing request token
- `404` Not Found for:
  - `departmentId` is not found
- `500` Server Error
