CREATE TABLE IF NOT EXISTS courses (
    id VARCHAR(10) PRIMARY KEY,
    description TEXT,
    course_name VARCHAR(255) NOT NULL,
    course_full_name VARCHAR(255) NOT NULL,
    course_type VARCHAR(255),
    grading_type ENUM('Letter Grade', 'S/U') NOT NULL DEFAULT 'S/U',
    faculty ENUM(
    'Science', 
    'Arts', 
    'Medicine', 
    'Law', 
    'Allied Health Science', 
    'Architecture',
    'Commerce and Accountancy', 
    'Communication Arts', 
    'Dentistry', 
    'Economics', 
    'Education', 
    'Engineering', 
    'Fine and Applied Arts', 
    'Nursing',
    'Pharmaceutical Sciences',
    'Political Science',
    'Psychology',
    'Sports Science',
    'Veterinary Science',
    'Integrated Innovation',
    'Agricultural Resources',
    'Graduate School'
    ),
    midterm_exam_date DATE,
    final_exam_date DATE,
    credit INT NOT NULL DEFAULT 3,
    course_group_id INT
);