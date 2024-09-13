CREATE TABLE IF NOT EXISTS section_instructors (
  id INT AUTO_INCREMENT PRIMARY KEY,
  section_id INT NOT NULL,
  instructor_id INT NOT NULL,
  FOREIGN KEY (section_id) REFERENCES sections(id),
  FOREIGN KEY (instructor_id) REFERENCES instructors(id)
);
