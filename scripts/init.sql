-- Create authors table
CREATE TABLE IF NOT EXISTS authors (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL
);

-- Create articles table
CREATE TABLE IF NOT EXISTS articles (
    id TEXT PRIMARY KEY,
    author_id TEXT NOT NULL,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (author_id) REFERENCES authors(id)
);

-- Insert sample data
INSERT INTO authors (id, name) VALUES 
    ('author-1', 'John Doe'),
    ('author-2', 'Jane Smith')
ON CONFLICT (id) DO NOTHING;

INSERT INTO articles (id, author_id, title, body) VALUES 
    ('article-1', 'author-1', 'Getting Started with Go', 'This is a comprehensive guide to getting started with Go programming language.'),
    ('article-2', 'author-2', 'Database Design Best Practices', 'Learn the essential principles of database design and normalization.')
ON CONFLICT (id) DO NOTHING;
