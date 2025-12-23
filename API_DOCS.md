# API Documentation

## Base URL
`http://localhost:3000`

## Health Check
- **GET** `/health`
  - Returns status of the server.
  - Response: `{"status": "success", "message": "Backend is reachable"}`

## Users

### Create User (Register)
- **POST** `/api/users`
- **Body**:
  ```json
  {
    "username": "cool_gamer",
    "email": "gamer@example.com",
    "password": "securepassword123"
  }
  ```
- **Response**: User object (without password).

### Get Users
- **GET** `/api/users`
- **Response**: List of users.

### Get User by ID
- **GET** `/api/users/:id`
- **Response**: User details.

### Update User
- **PUT** `/api/users/:id`
- **Body** (fields are optional):
  ```json
  {
    "username": "new_username",
    "email": "new@example.com"
  }
  ```

### Delete User
- **DELETE** `/api/users/:id`

---

## Games

### Create Game
- **POST** `/api/games`
- **Body**:
  ```json
  {
    "user_id": "uuid-of-user",
    "title": "Elden Ring",
    "cover_url": "https://example.com/eldenring.jpg",
    "genre": "Action RPG",
    "status": "playing",
    "score": 9.5,
    "hours_played": 45,
    "hltb_estimate": 100,
    "release_year": 2022,
    "review_text": "Amazing game!",
    "tag_ids": ["uuid-of-tag-1", "uuid-of-tag-2"]
  }
  ```
  - *Note*: `status` can be "backlog", "playing", "completed", "dropped".

### Get Games
- **GET** `/api/games`
- **Response**: List of games with their tags.

### Get Game by ID
- **GET** `/api/games/:id`

### Update Game
- **PUT** `/api/games/:id`
- **Body** (any field can be updated):
  ```json
  {
    "status": "completed",
    "hours_played": 120,
    "score": 10.0,
    "tag_ids": [] 
  }
  ```
  - *Note*: Sending an empty `tag_ids` array will remove all tags. Sending `null` or omitting it will keep existing tags.

### Delete Game
- **DELETE** `/api/games/:id`

---

## Tags

### Create Tag
- **POST** `/api/tags`
- **Body**:
  ```json
  {
    "name": "Cozy",
    "user_id": "optional-uuid"
  }
  ```

### Get Tags
- **GET** `/api/tags`

### Update Tag
- **PUT** `/api/tags/:id`
- **Body**:
  ```json
  {
    "name": "Hardcore"
  }
  ```

### Delete Tag
- **DELETE** `/api/tags/:id`
