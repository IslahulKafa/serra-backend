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

#### Request OTP

- **POST** `http:localhost:8080/api/v1/request-otp`
- **Body:**
  ```json
  {
    "email": "string"
  }
  ```
- **Response:** `200 OK`
  ```json
  {
    "message": "OTP sent"
  }
  ```

#### Verify OTP

- **POST** `http:localhost:8080/api/v1/verify-otp`
- **Body:**
  ```json
  {
    "email": "string",
    "otp": "string",
    "otp_token": "jwt_token"
  }
  ```
- **Response:** `200 OK`
  ```json
  {
    "message": "OTP verified"
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

### 2. Keys

#### Upload keys

- **POST** `http:localhost:8080/api/v1/keys/upload`
- **Headers:**
  - `Authorization: Bearer <token>`
- **Body:**
  ```json
  {
    "identity_key": "base64-identity-key",
    "signed_prekey": "base64-signed-prekey",
    "signed_prekey_signature": "base64-signature",
    "one_time_prekeys": [
      "base64-prekey-1",
      "base64-prekey-2",
      "base64-prekey-3"
    ]
  }
  ```
- **Response:** `201 Created`

#### Get keys

- **GET** `http:localhost:8080/api/v1/keys/{user_id}`
- **Headers:**
  - `Authorization: Bearer <token>`
- **Response:** `200 OK`
  ```json
  {
    "data": {
      "identity_key": "base64-identity-key",
      "ne_time_prekey": "base64-prekey-1",
      "signed_prekey": "base64-signed-prekey",
      "signed_prekey_signature": "base64-signature"
    },
    "status": "success"
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
