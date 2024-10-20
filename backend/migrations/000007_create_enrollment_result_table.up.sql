CREATE TABLE IF NOT EXISTS enrollment_results (
  id INT AUTO_INCREMENT PRIMARY KEY,
  user_id VARCHAR(255) NOT NULL,
  course_id VARCHAR(10) NOT NULL,
  course_name VARCHAR(255) NOT NULL,
  course_credit INT NOT NULL DEFAULT 3,
  section_id INT NOT NULL,
  section INT NOT NULL DEFAULT 1,
  round VARCHAR(10) NOT NULL,
  FOREIGN KEY (course_id) REFERENCES courses(id),
  FOREIGN KEY (section_id) REFERENCES sections(id)
);
