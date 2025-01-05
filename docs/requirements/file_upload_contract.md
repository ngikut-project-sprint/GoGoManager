# File Upload

## PIC

...

## Background:

Manager can manage uploaded files

## Contract:

**POST /v1/file**

Request Header:

|      key      |   value    |
| :-----------: | :--------: |
| Authorization | bearer ... |

Request Multipart Form-Data:

| key  |       value       |
| :--: | :---------------: |
| file | file (max 100KiB) |

Response:

- `200` Ok

```js
{
  "fileId": "", // use whatever id you want
  "uri": "name@name.com/file.jpg" // should be the URI
}
```

- `400` Bad Request case:
  - Validation error
- `401` Unauthorized for
  - expired / invalid / missing request token
- `500` Server Error
