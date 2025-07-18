-- Create database if it doesn't exist
CREATE DATABASE moviedb;

-- Connect to the database
\c moviedb;

-- Create extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- The tables will be created automatically by GORM Auto-Migration
-- This file is just for any additional database setup if needed

-- Create indexes for better performance
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_movies_title ON movies(title);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_movies_genre ON movies(genre);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_movies_release_year ON movies(release_year);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_view_histories_user_id ON view_histories(user_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_view_histories_movie_id ON view_histories(movie_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_view_histories_watched_at ON view_histories(watched_at);

-- Insert default admin user (password: admin123)
INSERT INTO users (username, email, password, role, created_at, updated_at) 
VALUES (
  'admin', 
  'admin@movieapi.com', 
  '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 
  'admin', 
  NOW(), 
  NOW()
) ON CONFLICT (username) DO NOTHING;
