package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	store *Store
}

func main() {
	store := NewStore()
	server := &Server{store: store}

	r := mux.NewRouter()
	r.HandleFunc("/employees", server.ListEmployees).Methods("GET")
	r.HandleFunc("/employees", server.CreateEmployee).Methods("POST")
	r.HandleFunc("/employees/batch", server.BatchCreateEmployees).Methods("POST")
	r.HandleFunc("/employees/{id:[0-9]+}", server.GetEmployeeByID).Methods("GET")
	r.HandleFunc("/employees/{id:[0-9]+}", server.UpdateEmployee).Methods("PUT")
	r.HandleFunc("/employees/{id:[0-9]+}", server.DeleteEmployee).Methods("DELETE")

	http.ListenAndServe(":5000", r)
}
func (s *Server) ListEmployees(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 50
	}

	s.store.Lock()
	defer s.store.Unlock()

	// Extract and sort employees by ID
	employees := make([]Employee, 0, len(s.store.employees))
	for _, emp := range s.store.employees {
		employees = append(employees, emp)
	}
	sort.Slice(employees, func(i, j int) bool {
		return employees[i].ID < employees[j].ID
	})

	// Pagination logic
	start := (page - 1) * limit
	end := start + limit
	totalEmployees := len(employees)
	if start >= totalEmployees {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	if end > totalEmployees {
		end = totalEmployees
	}

	response := map[string]interface{}{
		"page":      page,
		"limit":     limit,
		"total":     totalEmployees,
		"employees": employees[start:end],
	}
	final := JSONMessageWrappedObjwithStatus(http.StatusOK, response)
	WebResponseJSONObjectNoCache(w, r, http.StatusOK, final)
	return

}
func (s *Server) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var emp Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdEmp := s.store.CreateEmployee(emp.Name, emp.Position, emp.Salary)
	w.WriteHeader(http.StatusCreated)
	final := JSONMessageWrappedObjwithStatus(http.StatusOK, createdEmp)
	WebResponseJSONObjectNoCache(w, r, http.StatusOK, final)
	return
}
func (s *Server) BatchCreateEmployees(w http.ResponseWriter, r *http.Request) {
	var employees []Employee
	if err := json.NewDecoder(r.Body).Decode(&employees); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdEmployees := s.store.BatchCreateEmployees(employees)
	final := JSONMessageWrappedObjwithStatus(http.StatusOK, createdEmployees)
	WebResponseJSONObjectNoCache(w, r, http.StatusOK, final)
	return
}
func (s *Server) GetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	emp, exists := s.store.GetEmployeeByID(id)
	if !exists {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	final := JSONMessageWrappedObjwithStatus(http.StatusOK, emp)
	WebResponseJSONObjectNoCache(w, r, http.StatusOK, final)
	return
}

func (s *Server) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var emp Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedEmp, exists := s.store.UpdateEmployee(id, emp.Name, emp.Position, emp.Salary)
	if !exists {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	final := JSONMessageWrappedObjwithStatus(http.StatusOK, updatedEmp)
	WebResponseJSONObjectNoCache(w, r, http.StatusOK, final)
	return
}

func (s *Server) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	if !s.store.DeleteEmployee(id) {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}
	fmt.Println("Employee info delted")

	w.WriteHeader(http.StatusNoContent)
}

// JSONMessageObj returns an encoded JSON of the object provided.
func JSONMessageWrappedObjwithStatus(code int, obj interface{}) []byte {
	utc, _ := time.LoadLocation("UTC")
	currentUTCTime := time.Now().In(utc).Format("2006-01-02 15:04:05")
	jsonString := JsonWrappedContent{
		StatusCode:  code,
		LastUpdated: currentUTCTime,
		Content:     obj,
	}

	result, err := json.MarshalIndent(jsonString, "", "    ")
	if err != nil {
		fmt.Println(err)
	}

	return result
}

// WebResponseJSONObjectNoCache is a wrapper function that returns an already prepared JSON object as web response with the no-cache header added.
func WebResponseJSONObjectNoCache(w http.ResponseWriter, r *http.Request, code int, obj interface{}) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS,POST,PUT,UPDATE,DELETE")
	w.WriteHeader(code)
	w.Write(obj.([]byte))
}
