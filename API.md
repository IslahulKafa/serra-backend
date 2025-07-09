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
