# Movie Streaming API Backend

A RESTful API backend for a movie streaming platform built with Go (Gin) and PostgreSQL.

## Features

- User authentication (register, login, JWT tokens)
- Admin and user roles
- Movie upload (admin only)
- Movie streaming with range support
- Movie search and filtering
- View history tracking
- PostgreSQL database with GORM

## Prerequisites

- Go 1.21+
- PostgreSQL
- Git

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd movie-api-backend
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up PostgreSQL database:
```sql
CREATE DATABASE moviedb;
```

4. Configure environment variables:
```bash
cp .env.example .env
# Edit .env with your database credentials
```

5. Run the application:
```bash
go run cmd/main.go
```

## API Endpoints

### Authentication

#### Register User
```
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123",
  "role": "user" // optional, defaults to "user"
}
```

#### Login
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}
```

#### Get Profile
```
GET /api/v1/auth/profile
Authorization: Bearer <token>
```

### Movies

#### Get Movies (with pagination and search)
```
GET /api/v1/movies?page=1&limit=10&genre=action&search=title
```

#### Get Single Movie
```
GET /api/v1/movies/:id
```

#### Stream Movie
```
GET /api/v1/movies/:id/stream
Authorization: Bearer <token>
```

#### Upload Movie (Admin only)
```
POST /api/v1/movies
Authorization: Bearer <admin_token>
Content-Type: multipart/form-data

Form fields:
- movie: video file
- title: string
- description: string
- genre: string
- director: string
- release_year: number
- duration: number (minutes)
- rating: number
```

#### Delete Movie (Admin only)
```
DELETE /api/v1/movies/:id
Authorization: Bearer <admin_token>
```

### User

#### Get View History
```
GET /api/v1/user/history
Authorization: Bearer <token>
```

## Environment Variables

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=moviedb
DB_SSL=disable

JWT_SECRET=your-secret-key-here
PORT=8080

UPLOAD_PATH=./uploads/movies
MAX_FILE_SIZE=1073741824
```

## Video Streaming

The API supports HTTP range requests for efficient video streaming. This allows:
- Seeking to different parts of the video
- Resuming playback from where it left off
- Bandwidth optimization

## Database Schema

### Users Table
- id (primary key)
- username (unique)
- email (unique)
- password (hashed)
- role (user/admin)
- created_at, updated_at, deleted_at

### Movies Table
- id (primary key)
- title
- description
- genre
- director
- release_year
- duration
- rating
- file_path
- thumbnail_path
- file_size
- uploaded_by (foreign key to users)
- created_at, updated_at, deleted_at

### View History Table
- id (primary key)
- user_id (foreign key)
- movie_id (foreign key)
- watched_at

## Security Features

- JWT-based authentication
- Password hashing with bcrypt
- Role-based access control
- Input validation
- CORS support

## Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -o movie-api cmd/main.go
```

## API Usage Examples

### Create Admin User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "admin123",
    "role": "admin"
  }'
```

### Upload Movie
```bash
curl -X POST http://localhost:8080/api/v1/movies \
  -H "Authorization: Bearer <admin_token>" \
  -F "movie=@/path/to/movie.mp4" \
  -F "title=Sample Movie" \
  -F "description=A great movie" \
  -F "genre=Action" \
  -F "director=John Doe" \
  -F "release_year=2023" \
  -F "duration=120" \
  -F "rating=8.5"
```

### Stream Movie
```bash
curl -H "Authorization: Bearer <user_token>" \
  -H "Range: bytes=0-1048576" \
  http://localhost:8080/api/v1/movies/1/stream
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License.
# muve-service
