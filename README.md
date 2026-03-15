---

## `GET /api/chirps`

Retrieves a list of all chirps.

### Query Parameters

|Parameter|Type|Required|Default|Description|
|---|---|---|---|---|
|`sort`|`string`|No|`asc`|Sort order by `created_at`. Values: `asc`, `desc`|
|`author_id`|`string`|No|—|Filter chirps by author UUID|

### Response

**`200 OK`**

```json
[
  {
    "id": "uuid",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "body": "I'm the one who knocks!",
    "user_id": "uuid"
  }
]
```

### Examples

```
GET /api/chirps
GET /api/chirps?sort=asc
GET /api/chirps?sort=desc
GET /api/chirps?author_id=3fa85f64-5717-4562-b3fc-2c963f66afa6
GET /api/chirps?author_id=3fa85f64-5717-4562
```