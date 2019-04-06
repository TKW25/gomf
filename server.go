package main

import "fmt"

func main() {
	// people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
	// people = append(people, Person{ID: "2", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})
	// people = append(people, Person{ID: "3", Firstname: "Francis", Lastname: "Sunday"})

	// router := mux.NewRouter()
	// router.HandleFunc("/people", GetPeople).Methods("GET")
	// router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
	// router.HandleFunc("/people/{id}", CreatePerson).Methods("POST")
	// router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")
	// log.Fatal(http.ListenAndServe(":8000", router))
	LoadConfig("config/gomf.config")
	fmt.Println(Config)
	fileName, _ := GetRandomFileName("png")
	fmt.Println(fileName)
}
