#!/bin/bash

# Movie API Test Script

BASE_URL="http://localhost:8080/api/v1"

echo "🎬 Testing Movie API Backend"
echo "=========================="

# Test health check
echo "1. Testing health check..."
curl -s "$BASE_URL/../health" | jq '.'

# Test user registration
echo -e "\n2. Testing user registration..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }')
echo $REGISTER_RESPONSE | jq '.'

# Test admin registration
echo -e "\n3. Testing admin registration..."
ADMIN_REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "admin123",
    "role": "admin"
  }')
echo $ADMIN_REGISTER_RESPONSE | jq '.'

# Test user login
echo -e "\n4. Testing user login..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }')
echo $LOGIN_RESPONSE | jq '.'

# Extract token
USER_TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')

# Test admin login
echo -e "\n5. Testing admin login..."
ADMIN_LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }')
echo $ADMIN_LOGIN_RESPONSE | jq '.'

# Extract admin token
ADMIN_TOKEN=$(echo $ADMIN_LOGIN_RESPONSE | jq -r '.token')

# Test profile endpoint
echo -e "\n6. Testing profile endpoint..."
curl -s "$BASE_URL/auth/profile" \
  -H "Authorization: Bearer $USER_TOKEN" | jq '.'

# Test get movies (should be empty initially)
echo -e "\n7. Testing get movies..."
curl -s "$BASE_URL/movies" | jq '.'

# Test get view history (should be empty)
echo -e "\n8. Testing view history..."
curl -s "$BASE_URL/user/history" \
  -H "Authorization: Bearer $USER_TOKEN" | jq '.'

echo -e "\n✅ API tests completed!"
echo "📝 Notes:"
echo "- To upload movies, use: curl -X POST $BASE_URL/movies -H \"Authorization: Bearer \$ADMIN_TOKEN\" -F \"movie=@/path/to/video.mp4\" -F \"title=Movie Title\""
echo "- To stream movies, use: curl -H \"Authorization: Bearer \$USER_TOKEN\" $BASE_URL/movies/{id}/stream"
echo "- Make sure PostgreSQL is running before starting the API"
