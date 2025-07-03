# Hacker News Clone API Documentation

*generated with Claude, with my `main.go` used as context*

Base URL: `http://localhost/api/v1`

---

## Authentication

Most endpoints require authentication via JWT token. Include the token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

You can obtain a JWT token by using the /login method.

---

## `GET /`

**Description:**  
Check authentication status. Returns whether the user is logged in and their username.

### Headers
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `Authorization` | string | Yes | Bearer token for authentication |

### Sample Response
```json
{
  "Message": "Logged in as john_doe",
  "Status": 200
}
```

### Possible HTTP Status Codes
- `200 OK` – User is authenticated
- `401 Unauthorized` – Not signed in or invalid token

---

## `POST /api/v1/login`

**Description:**  
Initiate login process by sending a magic link to the provided email address.

### Request Body Parameters
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `email` | string | Yes | User's email address (max 100 chars) |
| `username` | string | Yes | Desired username (max 100 chars, alphanumeric + underscores) |
| `captchaToken` | string | Yes | Google reCAPTCHA token |

### Sample Request
```json
{
  "email": "user@example.com",
  "username": "john_doe",
  "captchaToken": "03AGdBq24..."
}
```

### Sample Response
```json
{
  "message": "Emailed a magic link to user@example.com"
}
```

### Possible HTTP Status Codes
- `200 OK` – Magic link sent successfully
- `400 Bad Request` – Invalid input (missing fields, invalid email/username format, failed captcha)

---

## `GET /api/v1/magic`

**Description:**  
Validate magic link token and return JWT for authentication.

### Query Parameters
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `token` | string | Yes | Magic link token from email |

### Sample Response
```json
{
  "username": "john_doe",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Possible HTTP Status Codes
- `200 OK` – Login successful, JWT returned
- `400 Bad Request` – Missing token parameter
- `502 Bad Gateway` – Invalid or expired magic link

---

## `POST /api/v1/submit`

**Description:**  
Submit a new post/link to Hacker News.

### Headers
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `Authorization` | string | Yes | Bearer token for authentication |

### Request Body Parameters
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `link` | string | Yes | URL to submit (max 255 chars) |
| `title` | string | Yes | Title of the submission |
| `body` | string | No | Optional text body |
| `captchaToken` | string | Yes | Google reCAPTCHA token |

### Sample Request
```json
{
  "link": "https://example.com/article",
  "title": "Interesting Article About Go",
  "body": "This article discusses...",
  "captchaToken": "03AGdBq24..."
}
```

### Sample Response
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000"
}
```

### Possible HTTP Status Codes
- `201 Created` – Submission created successfully
- `400 Bad Request` – Invalid input (missing fields, invalid URL)
- `401 Unauthorized` – Not authenticated

---

## `POST /api/v1/bio`

**Description:**  
Update user profile metadata (bio, full name, birthdate).

### Headers
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `Authorization` | string | Yes | Bearer token for authentication |

### Request Body Parameters
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `fullName` | string | No | User's full name (max 100 chars) |
| `birthdate` | string | No | Birth date in MM-DD-YYYY format |
| `bioText` | string | No | Biography text |

### Sample Request
```json
{
  "fullName": "John Doe",
  "birthdate": "01-15-1990",
  "bioText": "Software developer interested in Go and distributed systems."
}
```

### Sample Response
```json
{
  "message": "Updated metadata for user john_doe"
}
```

### Possible HTTP Status Codes
- `200 OK` – Profile updated successfully
- `400 Bad Request` – Invalid input (wrong date format, name too long)
- `401 Unauthorized` – Not authenticated

---

## `POST /api/v1/vote`

**Description:**  
Vote on a submission (upvote or downvote).

### Headers
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `Authorization` | string | Yes | Bearer token for authentication |

### Request Body Parameters
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | string | Yes | ID of the submission to vote on |
| `upvote` | bool | Yes | true for upvote, false for downvote |

### Sample Request
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "upvote": true
}
```

### Sample Response
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "voteSuccess": true
}
```

### Possible HTTP Status Codes
- `200 OK` – Vote processed (success indicates if new vote or duplicate)
- `400 Bad Request` – Missing ID parameter
- `401 Unauthorized` – Not authenticated

---

## `GET /api/v1/submission`

**Description:**  
Get details and vote count for a specific submission.

### Query Parameters
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | string | Yes | ID of the submission |

### Sample Response
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "metdata": {
    "title": "Interesting Article About Go",
    "link": "https://example.com/article",
    "body": "This article discusses...",
    "author": "john_doe",
    "isFlagged": false
  },
  "votes": {
    "upvotes": 42,
    "downvotes": 3,
    "total": 39
  }
}
```

### Possible HTTP Status Codes
- `200 OK` – Submission data returned
- `400 Bad Request` – Missing ID parameter
- `500 Internal Server Error` – Database error

---

## `GET /api/v1/user`

**Description:**  
Get public profile information for a user.

### Query Parameters
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `username` | string | Yes | Username to look up |

### Sample Response
```json
{
  "username": "john_doe",
  "joined": "2024-01-15T10:30:00Z",
  "metadata": {
    "full_name": "John Doe",
    "birthday": "01-15-1990",
    "bio": "Software developer interested in Go and distributed systems."
  }
}
```

### Possible HTTP Status Codes
- `200 OK` – User data returned
- `400 Bad Request` – Missing username parameter

---

## `GET /api/v1/me`

**Description:**  
Shorthand endpoint that redirects to the authenticated user's profile.

### Headers
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `Authorization` | string | Yes | Bearer token for authentication |

### Sample Response
Redirects to `/api/v1/user?username=<authenticated_username>`

### Possible HTTP Status Codes
- `302 Found` – Redirect to user profile
- `401 Unauthorized` – Not authenticated

---

## `DELETE /api/v1/submission`

**Description:**  
Delete a submission. Only the author can delete their own posts, and flagged posts cannot be deleted.

### Headers
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `Authorization` | string | Yes | Bearer token for authentication |

### Request Body Parameters
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | string | Yes | ID of the submission to delete |

### Sample Request
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000"
}
```

### Sample Response
```json
{
  "message": "OK"
}
```

### Possible HTTP Status Codes
- `200 OK` – Submission deleted successfully
- `400 Bad Request` – Missing ID parameter
- `401 Unauthorized` – Not authenticated
- `403 Forbidden` – User doesn't own this post
- `404 Not Found` – Submission not found
- `423 Locked` – Post is flagged and under review, and cannot be deleted at this time

---

## `GET /api/v1/status`

**Description:**  
Health check endpoint to verify API is running.

### Sample Response
```json
{
  "message": "Healthy",
  "status": 200
}
```

### Possible HTTP Status Codes
- `200 OK` – API is healthy

---

## Error Responses

All error responses follow a consistent format:

```json
{
  "error": "Description of what went wrong"
}
```

Or for some endpoints:

```json
{
  "message": "Error description",
  "status": 400
}
```

---

## Rate Limiting

The API implements rate limiting to prevent abuse. If you exceed the rate limit, you'll receive a `429 Too Many Requests` response.

---

## Notes

- All timestamps are in ISO 8601 format
- Usernames must contain only letters, numbers, and underscores
- Email addresses are validated against standard email regex
- Google reCAPTCHA is required for login and submission endpoints
- Magic links expire after a certain time period
- JWT tokens expire after 60 minutes