package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/chingsley/go-movies-api/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	collection *mongo.Collection
}

func (db *DB) GetMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var movie models.Movie
	objectID, _ := primitive.ObjectIDFromHex(vars["id"])
	filter := bson.M{"_id": objectID}
	err := db.collection.FindOne(context.TODO(), filter).Decode(&movie)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(movie)
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func (db *DB) PostMovie(w http.ResponseWriter, r *http.Request) {
	var movie models.Movie
	postBody, _ := ioutil.ReadAll(r.Body)

	json.Unmarshal(postBody, &movie)
	result, err := db.collection.InsertOne(context.TODO(), movie)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(result)
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://getir_assignment_tester2:6NyfP6gvkykZ9Pcd@cluster0.orvcd.mongodb.net/getir_assignment_test?retryWrites=true&w=majority")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database("getir_assignment_test").Collection("movies")
	db := &DB{collection: collection}

	r := mux.NewRouter()
	r.HandleFunc("/v1/movies/{id:[a-zA-Z0-9]*}",
		db.GetMovie).Methods("GET")
	r.HandleFunc("/v1/movies", db.PostMovie).Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println(fmt.Sprintf("Staring application on port %s", "8000"))
	log.Fatal(srv.ListenAndServe())
}
