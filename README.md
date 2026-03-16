## Table of Contents

- [GET /api/chirps](#get-apichirps)
- [GET /api/chirps/{chirpID}](#get-apichirpschirpid)
- [POST /api/chirps](#post-apichirps)
- [DELETE /api/chirps/{chirpID}](#delete-apichirpschirpid)
- [GET /api/healthz](#get-apihealthz)
- [POST /api/users](#post-apiusers)
- [PUT /api/users](#put-apiusers)
- [POST /api/login](#post-apilogin)
- [POST /api/refresh](#post-apirefresh)
- [POST /api/revoke](#post-apirevoke)
- [POST /api/polka/webhooks](#post-apipolkawebhooks)
- [GET /admin/metrics](#get-adminmetrics)
- [POST /admin/reset](#post-adminreset)

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

## `GET /api/chirps/{chirpID}`

Retrieves a single chirp by ID.

### Path Parameters

|Parameter|Type|Required|Description|
|---|---|---|---|
|`chirpID`|`string (uuid)`|Yes|The chirp ID to retrieve|

### Response

**`200 OK`**

```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "body": "I'm the one who knocks!",
  "user_id": "uuid"
}
```

### Error Responses

- `400 Bad Request` for invalid UUID
- `404 Not Found` when chirp does not exist

## `POST /api/chirps`

Creates a chirp for the authenticated user.

### Headers

|Header|Required|Description|
|---|---|---|
|`Authorization`|Yes|Bearer access token (`Bearer <jwt>`)|

### Request Body

```json
{
  "body": "This is my chirp"
}
```

### Notes

- Max length is `140` characters.
- Profanity filtering replaces these words (case-insensitive) with `****`: `kerfuffle`, `sharbert`, `fornax`.

### Response

**`201 Created`**

```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "body": "This is my chirp",
  "user_id": "uuid"
}
```

### Error Responses

- `400 Bad Request` when chirp is too long
- `401 Unauthorized` when token is missing/invalid

## `DELETE /api/chirps/{chirpID}`

Deletes a chirp owned by the authenticated user.

### Path Parameters

|Parameter|Type|Required|Description|
|---|---|---|---|
|`chirpID`|`string (uuid)`|Yes|The chirp ID to delete|

### Headers

|Header|Required|Description|
|---|---|---|
|`Authorization`|Yes|Bearer access token (`Bearer <jwt>`)|

### Response

**`204 No Content`**

### Error Responses

- `400 Bad Request` for invalid UUID
- `401 Unauthorized` when token is missing/invalid
- `403 Forbidden` when chirp belongs to a different user
- `404 Not Found` when chirp does not exist

## `GET /api/healthz`

Returns server health status.

### Response

**`200 OK`** (`text/plain`)

```text
OK
```

## `POST /api/users`

Creates a new user.

### Request Body

```json
{
  "email": "user@example.com",
  "password": "super-secret"
}
```

### Response

**`201 Created`**

```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "email": "user@example.com",
  "is_chirpy_red": false
}
```

## `PUT /api/users`

Updates the authenticated user's email and password.

### Headers

|Header|Required|Description|
|---|---|---|
|`Authorization`|Yes|Bearer access token (`Bearer <jwt>`)|

### Request Body

```json
{
  "email": "new-email@example.com",
  "password": "new-password"
}
```

### Response

**`200 OK`**

```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-02T00:00:00Z",
  "email": "new-email@example.com",
  "is_chirpy_red": false
}
```

### Error Responses

- `400 Bad Request` for malformed JSON/unknown fields
- `401 Unauthorized` when token is missing/invalid

## `POST /api/login`

Authenticates a user and returns access and refresh tokens.

### Request Body

```json
{
  "email": "user@example.com",
  "password": "super-secret"
}
```

### Response

**`200 OK`**

```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "email": "user@example.com",
  "is_chirpy_red": false,
  "token": "<jwt>",
  "refresh_token": "<refresh_token>"
}
```

### Notes

- Access token expires after `1 hour`.
- Refresh token expires after `60 days`.

### Error Responses

- `401 Unauthorized` for incorrect email/password

## `POST /api/refresh`

Exchanges a valid refresh token for a new access token.

### Headers

|Header|Required|Description|
|---|---|---|
|`Authorization`|Yes|Bearer refresh token (`Bearer <refresh_token>`)|

### Response

**`200 OK`**

```json
{
  "token": "<jwt>"
}
```

### Error Responses

- `400 Bad Request` when token header is missing/malformed
- `401 Unauthorized` for invalid, revoked, or expired refresh token

## `POST /api/revoke`

Revokes a refresh token.

### Headers

|Header|Required|Description|
|---|---|---|
|`Authorization`|Yes|Bearer refresh token (`Bearer <refresh_token>`)|

### Response

**`204 No Content`**

### Error Responses

- `400 Bad Request` when token header is missing/malformed
- `500 Internal Server Error` when revoke fails

## `POST /api/polka/webhooks`

Receives Polka webhook events.

### Headers

|Header|Required|Description|
|---|---|---|
|`Authorization`|Yes|API key header format expected by server (`ApiKey <key>`)|

### Request Body

```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "uuid"
  }
}
```

### Behavior

- If `event` is not `user.upgraded`, returns `204 No Content`.
- If `event` is `user.upgraded`, the matching user is upgraded to Chirpy Red.

### Responses

- `204 No Content` on success or ignored event
- `401 Unauthorized` for invalid API key
- `400 Bad Request` for malformed JSON or invalid `user_id`
- `404 Not Found` when `user_id` does not exist

## `GET /admin/metrics`

Returns a simple admin HTML page showing app hit count.

### Response

**`200 OK`** (`text/html`)

```html
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited X times!</p>
  </body>
</html>
```

## `POST /admin/reset`

Resets users and file-server hit counter. Only available when `PLATFORM=dev`.

### Response

- `200 OK` with body `Hits reset to 0` (when `PLATFORM=dev`)
- `403 Forbidden` otherwise