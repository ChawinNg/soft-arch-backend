CREATE TABLE IF NOT EXISTS timeslots (
  id INT AUTO_INCREMENT PRIMARY KEY,
  time VARCHAR(255),
  section_id INT NOT NULL,
  FOREIGN KEY (section_id) REFERENCES sections(id)
);
