package main

import (
	"database/sql"
	"log"
	"encoding/json"
	"net/http"
	"io/ioutil"
	
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var (
	mssqldb *sql.DB
)
//Person model
type Person struct {
	ID       int  `json: "id, omitempty"`
	Name     string  `json: "name, omitempty"`
	LastName string  `json: "last_name, omitempty"`
	Age      int     `json: "age, omitempty"`
	Address  Address `json:"address, omitempty"`
	AddressId int `json:"addressId"`
	Status int `json:"status"`
}

//Address model
type Address struct {
	ID     int  `json: "id, omitempty"`
	City   string `json:"city, omitempty"`
	Sector string `json:"sector, omitempty"`
}

//GetPeoples get all people
func GetPeoples(res http.ResponseWriter, req *http.Request) {
	
	var people []Person
	results, err := mssqldb.Query("SELECT * FROM people where status = 0")

	if err != nil {
		panic(err.Error())
	}

	for results.Next(){
		var person Person 
		
		err = results.Scan(&person.ID, &person.Name, &person.LastName, &person.Age, &person.AddressId, &person.Status)	
		if err!= nil{
			panic(err.Error())
		}else{
			var address Address
			address = getAdreess(person.AddressId)
			person.Address = address
			people = append(people, person)
		}
	}
	json.NewEncoder(res).Encode(people)
}

//GetPerson get all people
func GetPerson(res http.ResponseWriter, req *http.Request) {
	var person Person
	params :=mux.Vars(req)
	id := params ["id"]
//the 0 status meand the people is´snt deteted, the 1 status means it is
	result, err := mssqldb.Query("SELECT * FROM PEOPLE WHERE ID = ? ",id)
	if err != nil{
		panic(err.Error)
	}
	for result.Next(){
		err = result.Scan(&person.ID, &person.Name, &person.LastName, &person.Age, &person.AddressId, &person.Status)
		person.Address = getAdreess(person.AddressId)
		if err != nil{
			panic(err.Error())
		}
	}
	if person.ID != 0{
		json.NewEncoder(res).Encode(person)
	}else{
	json.NewEncoder(res).Encode("EL usuario que busca no existe")
	}
}

//newPerson get all people
func newPerson(res http.ResponseWriter, req *http.Request) {

	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil{
		http.Error(res, err.Error(), 500)
		return
	}
	//unmarchal
	var person Person
	err = json.Unmarshal(body, &person)
	if err != nil{
		http.Error(res, err.Error(), 500)
		return
	}
	insert, err := mssqldb.Query("insert into people (name, lastname, age, addressId) values (?, ?, ?, ?)", person.Name, person.LastName, person.Age , person.AddressId)
	if err != nil{
		panic(err.Error())
	}else{
		
		http.Redirect(res, req,  "http://localhost:3000/peoples", 301)
		 insert.Close()
	}
}

//editPerson get all people
func editPerson(res http.ResponseWriter, req *http.Request) {
	params:= mux.Vars(req)
	id :=params["id"]
	if id != ""{
	var person Person
	//new dates
	body, err := ioutil.ReadAll(req.Body)
	if err !=  nil{
		http.Error(res, err.Error(), 500)
	}
	err = json.Unmarshal(body, &person)
	if err !=  nil{
		http.Error(res, err.Error(), 500)
	}
	update, err:= mssqldb.Query("UPDATE people SET name =?, lastname=?, age=?, addressId=? WHERE ID = ?", person.Name, person.LastName, person.Age, person.AddressId, id)
	if err !=  nil{
		http.Error(res, err.Error(), 500)
	}else{
		http.Redirect(res, req,  "http://localhost:3000/peoples", 301)
		defer update.Close()
	}
	}else{
		json.NewEncoder(res).Encode("Nesecita enviar el Id del usuario a editar")
	}
}

//removePerson get all people
func removePerson(res http.ResponseWriter, req *http.Request) {
	
	params:= mux.Vars(req)
	id :=params["id"]
	if id!=""{
	update, err := mssqldb.Query("UPDATE people SET status = 1 WHERE ID = ?", id)
	
	if err !=  nil{
		http.Error(res, err.Error(), 500)
	}else{
		http.Redirect(res, req,  "http://localhost:3000/peoples", 301)
		defer update.Close()
	}
	
	}else{
		json.NewEncoder(res).Encode("Nesecita enviar el Id del usuario a editar")
	}
}
//getAddress method for attaching to persons
func getAdreess(id int) Address{
	addressResult, err:= mssqldb.Query("SELECT * FROM address WHERE ID = ?", id)
	var address Address	
	for addressResult.Next(){
		
		err = addressResult.Scan(&address.ID, &address.City, &address.Sector)	
		
		if err != nil{
			panic(err.Error())
			return Address{ID: 0, City: "", Sector: ""}
		}
			
	}
		return address
}

func openDBServer() (*sql.DB, error) {
	// Open SQL Server DB
	db, err := sql.Open("mysql", "username:password@/dbname")
	if err != nil {
		return nil, err
	}
	// Ping SQLServer
	err = db.Ping()
	if err != nil {
		return db, err
	}
	return db, err
}

func main() {

	// Connect to DB
	db, err := openDBServer()
	
	if err != nil {
		panic(err.Error())
	}

	mssqldb = db
	defer mssqldb.Close()
	
	router := mux.NewRouter()

	router.HandleFunc("/peoples", GetPeoples).Methods("GET")
	router.HandleFunc("/peoples/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/peoples", newPerson).Methods("POST")
	router.HandleFunc("/peoples/{id}/edit", editPerson).Methods("PUT")
	router.HandleFunc("/peoples/{id}", removePerson).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", router))
}