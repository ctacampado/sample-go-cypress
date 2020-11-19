package todo

// InitRoutes initializes all todo routes
func (s Service) InitRoutes() {
	r := s.Service.Router
	r.HandleFunc("/", s.GetLists).Methods("GET")
	r.HandleFunc("/list", s.GetLists).Methods("GET")
	r.HandleFunc("/list/{id}", s.GetListByID).Methods("GET")
	r.HandleFunc("/list", s.AddNewList).Methods("POST")
	r.HandleFunc("/list", s.DeleteList).Methods("DELETE")
	r.HandleFunc("/list/add", s.AddTodo).Methods("POST")
	r.HandleFunc("/list/delete", s.DeleteTodo).Methods("DELETE")
}
