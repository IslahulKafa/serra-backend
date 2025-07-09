# Serra API Documentation

## Overview

This document describes the RESTful API endpoints for the Serra project.

---

## Authentication

All endpoints require authentication via a Bearer token in the `Authorization` header.

---

## Endpoints

### 1. Users

#### Register a new user

- **POST** `http:localhost:8080/api/v1/register`
- **Body:**
  ```json
  {
    "email": "string",
    "password": "string"
  }
  ```
- **Response:** `201 Created`

#### Login

- **POST** `http:localhost:8080/api/v1/login`
- **Body:**
  ```json
  {
    "email": "string",
    "password": "string"
  }
  ```
- **Response:** `200 OK`
  ```json
  {
    "token": "jwt-token-string"
  }
  ```

#### Get user profile

- **GET** `http:localhost:8080/api/v1/me`
- **Headers:**
  - `Authorization: Bearer <token>`
- **Response:** `200 OK`
  ```json
  {
    "id": "user-id",
    "email": "user@example.com",
    "username": "string"
  }
  ```

## Error Handling

All errors return a JSON object:

```json
{
  "error": "Error message"
}
```

---

## Version

Current API version: `v1`
