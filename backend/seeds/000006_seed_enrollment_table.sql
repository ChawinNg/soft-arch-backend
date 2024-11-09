-- Insert seed data for specific users with specified conditions
INSERT INTO enrollments (user_id, course_id, course_name, course_credit, section_id, section, points, round, summarized) VALUES
-- User 672a433d591618f3bbc02221
('672a433d591618f3bbc02221', '2110413', 'COMP SECURITY', 3, 1, 21, FLOOR(RAND() * 9), '2024/1', FALSE),
('672a433d591618f3bbc02221', '2110201', 'COMP ENG MATH', 3, 2, 1, FLOOR(RAND() * 9), '2024/1', FALSE),
('672a433d591618f3bbc02221', '2110233', 'COMP ENG MATH LAB', 3, 3, 1, FLOOR(RAND() * 9), '2024/1', FALSE),
('672a433d591618f3bbc02221', 'SCI101', 'Intro to Science', 3, 8, 1, FLOOR(RAND() * 9), '2024/1', FALSE),
('672a433d591618f3bbc02221', 'ART102', 'Western Art History', 3, 10, 1, FLOOR(RAND() * 9), '2024/1', FALSE),

-- User 67287a0f15417e87427db6e1
('67287a0f15417e87427db6e1', '2110233', 'COMP ENG MATH LAB', 3, 4, 2, FLOOR(RAND() * 9), '2024/1', FALSE),
('67287a0f15417e87427db6e1', 'ART102', 'Western Art History', 3, 11, 2, FLOOR(RAND() * 9), '2024/1', FALSE),
('67287a0f15417e87427db6e1', 'MED103', 'Human Anatomy', 4, 12, 1, FLOOR(RAND() * 9), '2024/1', FALSE),
('67287a0f15417e87427db6e1', 'LAW104', 'Criminal Law', 3, 13, 1, FLOOR(RAND() * 9), '2024/1', FALSE),
('67287a0f15417e87427db6e1', 'AHS105', 'Medical Technologies', 3, 14, 1, FLOOR(RAND() * 9), '2024/1', FALSE),

-- User 6728795b32b1868cf860e0e8
('6728795b32b1868cf860e0e8', 'COM107', 'Financial Accounting', 3, 16, 1, FLOOR(RAND() * 9), '2024/1', FALSE),
('6728795b32b1868cf860e0e8', 'COMA108', 'Media and Journalism', 3, 17, 1, FLOOR(RAND() * 9), '2024/1', FALSE),
('6728795b32b1868cf860e0e8', 'DEN109', 'Intro to Dentistry', 4, 18, 1, FLOOR(RAND() * 9), '2024/1', FALSE),
('6728795b32b1868cf860e0e8', 'ECO110', 'Economics 101', 3, 19, 1, FLOOR(RAND() * 9), '2024/1', FALSE),
('6728795b32b1868cf860e0e8', 'EDU111', 'Philosophy of Education', 3, 20, 1, FLOOR(RAND() * 9), '2024/1', FALSE),

-- User 6727aeefa5c87b4e553fa352
('6727aeefa5c87b4e553fa352', 'ENG112', 'Engineering Mechanics', 4, 21, 1, FLOOR(RAND() * 9), '2024/1', FALSE),
('6727aeefa5c87b4e553fa352', 'FAA113', 'Sculpture Art', 3, 22, 1, FLOOR(RAND() * 9), '2024/1', FALSE),
('6727aeefa5c87b4e553fa352', 'NUR114', 'Nursing Practices', 3, 23, 1, FLOOR(RAND() * 9), '2024/1', FALSE),
('6727aeefa5c87b4e553fa352', 'PHAR115', 'Pharmaceutical Sciences', 4, 24, 1, FLOOR(RAND() * 9), '2024/1', FALSE),
('6727aeefa5c87b4e553fa352', 'POL116', 'Political Theories', 3, 25, 1, FLOOR(RAND() * 9), '2024/1', FALSE);
