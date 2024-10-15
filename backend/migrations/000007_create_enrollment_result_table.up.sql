CREATE TABLE IF NOT EXISTS enrollments (
  id INT AUTO_INCREMENT PRIMARY KEY,
  user_id INT NOT NULL,
  course_id VARCHAR(10) NOT NULL,
  section_id INT NOT NULL,
  section INT NOT NULL DEFAULT 1,
  round string NOT NULL DEFAULT "2024/1",
  FOREIGN KEY (course_id) REFERENCES courses(id),
  FOREIGN KEY (section_id) REFERENCES sections(id),
);
