package todo

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

// GetLists returns all todo lists
func (s Service) GetLists(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	w.Header().Add("Content-Type", "application/json")

	dbclient := s.Service.DBClient.Client
	collection := dbclient.Database("todo").Collection("lists")

	// Here's an array in which you can store the decoded documents
	var results []*List

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Find(ctx, bson.D{{}})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem List
		err := cur.Decode(&elem)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	rsp := MessageResponse{
		Message: "Ok",
		Data:    results,
	}

	bytes, err := json.Marshal(rsp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

// GetListByID returns a list given its ID
func (s Service) GetListByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	w.Header().Add("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]

	dbclient := s.Service.DBClient.Client
	collection := dbclient.Database("todo").Collection("lists")

	// Here's an array in which you can store the decoded documents
	var results []*List

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Find(ctx, bson.M{"id": id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem List
		err := cur.Decode(&elem)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	rsp := MessageResponse{
		Message: "Ok",
		Data:    results,
	}

	bytes, err := json.Marshal(rsp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

// AddNewList adds a new list
func (s Service) AddNewList(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	w.Header().Add("Content-Type", "application/json")

	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	var newListReq NewListRequest
	err = json.Unmarshal(req, &newListReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	newList := List{
		ID:          uuid.New().String(),
		Name:        newListReq.Name,
		List:        make([]Todo, 0),
		Description: newListReq.Description,
	}

	dbclient := s.Service.DBClient.Client
	collection := dbclient.Database("todo").Collection("lists")

	insertResult, err := collection.InsertOne(ctx, newList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	rsp := MessageResponse{
		Message: fmt.Sprintf("Inserted a single document: %s", insertResult.InsertedID),
		Data:    newList.ID,
	}

	bytes, err := json.Marshal(rsp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

// DeleteList deletes a list
func (s Service) DeleteList(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	w.Header().Add("Content-Type", "application/json")

	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	var deleteListReq DeleteListRequest
	err = json.Unmarshal(req, &deleteListReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	dbclient := s.Service.DBClient.Client
	collection := dbclient.Database("todo").Collection("lists")
	deleteResult, err := collection.DeleteOne(ctx, bson.M{"id": deleteListReq.ListID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	rsp := MessageResponse{
		Message: fmt.Sprintf("Deleted %v documents in the lists collection", deleteResult.DeletedCount),
		Data: struct {
			DeleteCount int64  `json:"deleteCount"`
			ID          string `json:"id"`
		}{
			DeleteCount: deleteResult.DeletedCount,
			ID:          deleteListReq.ListID,
		},
	}

	bytes, err := json.Marshal(rsp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

// AddTodo adds a new todo element to a list
func (s Service) AddTodo(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	w.Header().Add("Content-Type", "application/json")

	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	var addReq AddTodoRequest
	err = json.Unmarshal(req, &addReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	log.Printf("addReq %+v\n", addReq)

	newTodo := addReq.Todo
	newTodo.ID = uuid.New().String()

	log.Printf("newTodo %+v\n", newTodo)

	dbclient := s.Service.DBClient.Client
	collection := dbclient.Database("todo").Collection("lists")

	filter := bson.M{"id": addReq.ListID}
	update := bson.M{"$addToSet": bson.M{"list": newTodo}}
	updateResult, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	log.Printf("updateResult %+v\n", updateResult)

	rsp := MessageResponse{
		Message: fmt.Sprintf("updated %d document(s)", updateResult.MatchedCount),
		Data:    addReq.ListID,
	}

	bytes, err := json.Marshal(rsp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

// DeleteTodo deletes a todo element from a list
func (s Service) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	w.Header().Add("Content-Type", "application/json")

	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	var delReq DeleteTodoRequest
	err = json.Unmarshal(req, &delReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	log.Printf("delReq %+v\n", delReq)

	toDelete := delReq.TodoID

	log.Printf("toDelete %+v\n", toDelete)

	dbclient := s.Service.DBClient.Client
	collection := dbclient.Database("todo").Collection("lists")

	filter := bson.M{"id": delReq.ListID}
	update := bson.M{"$pull": bson.M{"list": bson.M{"id": toDelete}}}
	updateResult, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	log.Printf("updateResult %+v\n", updateResult)

	rsp := MessageResponse{
		Message: fmt.Sprintf("updated %d document(s)", updateResult.MatchedCount),
		Data:    delReq.ListID,
	}

	bytes, err := json.Marshal(rsp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
