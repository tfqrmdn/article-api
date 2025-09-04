-- Seeder: Insert sample authors
-- Created: 2025-09-04

INSERT INTO authors (id, name) VALUES 
    ('author-1', 'John Doe'),
    ('author-2', 'Jane Smith'),
    ('author-3', 'Bob Johnson'),
    ('author-4', 'Alice Brown'),
    ('author-5', 'Michael Chen'),
    ('author-6', 'Sarah Wilson'),
    ('author-7', 'David Rodriguez'),
    ('author-8', 'Emily Davis'),
    ('author-9', 'James Thompson'),
    ('author-10', 'Lisa Anderson'),
    ('author-11', 'Robert Taylor'),
    ('author-12', 'Jennifer Martinez'),
    ('author-13', 'William Garcia'),
    ('author-14', 'Amanda White'),
    ('author-15', 'Christopher Lee'),
    ('author-16', 'Michelle Harris'),
    ('author-17', 'Daniel Clark'),
    ('author-18', 'Ashley Lewis'),
    ('author-19', 'Matthew Walker'),
    ('author-20', 'Jessica Hall')
ON CONFLICT (id) DO NOTHING;
