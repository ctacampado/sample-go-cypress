package todo

import (
	"react-go-cypress/internal/service"
)

// Todo is an element of our todo list
type Todo struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// List contains a list of todos
type List struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	List        []Todo `json:"list"`
}

// Service is our Todo Service
type Service struct {
	Service *service.Service
}

// NewListRequest contains data for a newly
// created list
type NewListRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// DeleteListRequest contains data for list that
// is to be deleted
type DeleteListRequest struct {
	ListID string `json:"listID"`
}

// AddTodoRequest contains data for a todo
// element to be added to a list
type AddTodoRequest struct {
	ListID string `json:"listID"`
	Todo   Todo   `json:"todo"`
}

// DeleteTodoRequest contains data for the
// deletion of a list element
type DeleteTodoRequest struct {
	ListID string `json:"listID"`
	TodoID string `json:"todoID"`
}

// MessageResponse is the response for any
// list related request
type MessageResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
