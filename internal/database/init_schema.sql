-- Enum type for gender
CREATE TYPE GENDER AS ENUM ('male', 'female');

-- Managers table (1 manager -> N deparments, N employees indirectly)
CREATE TABLE managers (
  id SERIAL NOT NULL,
  email VARCHAR(255) NOT NULL,
  name VARCHAR(52),
  user_image_uri TEXT,
  company_name VARCHAR(52),
  company_image_uri TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP DEFAULT NULL,
  PRIMARY KEY(id),
  CONSTRAINT unique_email UNIQUE (email)
);

-- Departments table (1 department -> N employees)
CREATE TABLE departments (
  department_id SERIAL NOT NULL,
  name VARCHAR(33) NOT NULL,
  manager_id INT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP DEFAULT NULL,
  PRIMARY KEY(department_id),
  FOREIGN KEY(manager_id) REFERENCES managers(id)
);

-- Employees table (1 department -> N employees)
CREATE TABLE employees (
  id SERIAL NOT NULL,
  identity_number VARCHAR(33) NOT NULL,
  name VARCHAR(33) NOT NULL,
  employee_image_uri VARCHAR(33) NOT NULL,
  gender GENDER NOT NULL,
  department_id INT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP DEFAULT NULL,
  PRIMARY KEY(id),
  FOREIGN KEY (department_id) REFERENCES departments(department_id)
);
