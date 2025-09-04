-- Migration: Create authors table
-- Created: 2025-09-04

CREATE TABLE IF NOT EXISTS authors (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL
);
