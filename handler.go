package main

// CreateEmployee adds a new employee to the store
func (s *Store) CreateEmployee(name, position string, salary float64) Employee {
	s.Lock()
	defer s.Unlock()
	employee := Employee{
		ID:       s.nextID,
		Name:     name,
		Position: position,
		Salary:   salary,
	}
	s.employees[s.nextID] = employee
	s.nextID++
	return employee
}

// CreateEmployee adds a batch of employee to the store
func (s *Store) BatchCreateEmployees(employees []Employee) []Employee {
	s.Lock()
	defer s.Unlock()

	createdEmployees := make([]Employee, len(employees))
	for i, emp := range employees {
		emp.ID = s.nextID
		s.employees[s.nextID] = emp
		s.nextID++
		createdEmployees[i] = emp
	}
	return createdEmployees
}

// GetEmployeeByID retrieves an employee by ID
func (s *Store) GetEmployeeByID(id int) (Employee, bool) {
	s.Lock()
	defer s.Unlock()
	employee, exists := s.employees[id]
	return employee, exists
}

// UpdateEmployee updates an existing employee's details
func (s *Store) UpdateEmployee(id int, name, position string, salary float64) (Employee, bool) {
	s.Lock()
	defer s.Unlock()
	employee, exists := s.employees[id]
	if !exists {
		return Employee{}, false
	}
	employee.Name = name
	employee.Position = position
	employee.Salary = salary
	s.employees[id] = employee
	return employee, true
}

// DeleteEmployee removes an employee from the store
func (s *Store) DeleteEmployee(id int) bool {
	s.Lock()
	defer s.Unlock()
	_, exists := s.employees[id]
	if !exists {
		return false
	}
	delete(s.employees, id)
	return true
}
