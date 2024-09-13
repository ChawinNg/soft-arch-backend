CREATE TABLE IF NOT EXISTS sections (
  id INT AUTO_INCREMENT PRIMARY KEY,
  course_id VARCHAR(10) NOT NULL,
  section INT NOT NULL DEFAULT 1,
  room VARCHAR(255),
  capacity INT,
  FOREIGN KEY (course_id) REFERENCES courses(id)
);
