-- Seeder: Insert sample articles
-- Created: 2025-09-04

INSERT INTO articles (id, author_id, title, body) VALUES 
    ('article-1', 'author-1', 'Getting Started with Go', 'This is a comprehensive guide to getting started with Go programming language.'),
    ('article-2', 'author-2', 'Database Design Best Practices', 'Learn the essential principles of database design and normalization.'),
    ('article-3', 'author-1', 'Advanced Go Patterns', 'Explore advanced patterns and techniques in Go programming.'),
    ('article-4', 'author-3', 'Microservices Architecture', 'Understanding microservices architecture and implementation strategies.'),
    ('article-5', 'author-4', 'API Design Principles', 'Best practices for designing RESTful APIs.')
ON CONFLICT (id) DO NOTHING;
