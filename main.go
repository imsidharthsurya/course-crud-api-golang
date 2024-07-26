package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// creating models for course and it's author
type Course struct {
	CourseId    string  `json:"courseid"`
	CourseName  string  `json:"coursename"`
	CoursePrice int     `json:"price"`
	Author      *Author `json:"author"`
}

type Author struct {
	FullName string `json:"fullname"`
	Website  string `json:"website"`
}

// creating fake db ie. slice
var courses []Course

// middleware/helper function to check if valid body
// like courseId & coursename mandatory
func (c *Course) IsEmpty() bool {
	//b/c courseId we'll generate by ourselves
	return c.CourseName == ""
}

func main() {
	fmt.Println("Building Course Backend API")
	r := mux.NewRouter()

	//seeding the data into slice
	courses = append(courses, Course{CourseId: "1",
		CourseName:  "C++ & DSA",
		CoursePrice: 199,
		Author:      &Author{FullName: "Stiver", Website: "takeuforward.com"}})

	courses = append(courses, Course{CourseId: "2",
		CourseName:  "Frontend with react",
		CoursePrice: 1999,
		Author:      &Author{FullName: "Akshay", Website: "namastedev.com"}})

	//routing
	r.HandleFunc("/", serveHome).Methods("GET")
	//get all courses
	r.HandleFunc("/courses", getAllCourses).Methods("GET")
	//get a course with id
	r.HandleFunc("/course/{id}", getOneCourse).Methods("GET")
	//create a course
	r.HandleFunc("/course", createOneCourse).Methods("POST")
	//update one course grab by id
	r.HandleFunc("/course/{id}", updateOneCourse).Methods("PUT")
	//delete a course by id
	r.HandleFunc("/course/{id}", deleteOneCourse).Methods("DELETE")

	//host the server on port
	log.Fatal(http.ListenAndServe(":4000", r))
}

//controllers: later these will be on there own files & a seperate folder

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>This is home page of course backend api</h1>"))
}

func getAllCourses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get all courses")
	//to set header
	w.Header().Set("Content-Type", "application/json")
	//now throw the slice courses data as a json
	json.NewEncoder(w).Encode(courses)
}

func getOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get one course data")
	//set header
	w.Header().Set("Content-Type", "application/json")
	//get id from request using mux
	params := mux.Vars(r)

	//loop through the slice & return matching data
	for _, course := range courses {
		if course.CourseId == params["id"] {
			json.NewEncoder(w).Encode(course)
			return
		}

	}
	//out of for loop ie. no course found with given id
	json.NewEncoder(w).Encode("No Course found with given id")
	return
}

func createOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Creating a course")
	w.Header().Set("Content-Type", "application/json")
	//if body is empty
	if r.Body == nil {
		json.NewEncoder(w).Encode("Please send some data")
	}

	var course Course
	//decode the body json & store it in course var.
	_ = json.NewDecoder(r.Body).Decode(&course)

	//also check that isEmpty middleware
	if course.IsEmpty() {
		json.NewEncoder(w).Encode("Please send some data")
		return
	}

	//now generate unique id convert it into string
	//and store it into courseId
	//then append the course data
	rand.Seed(time.Now().UnixNano())
	course.CourseId = strconv.Itoa(rand.Intn(100))
	courses = append(courses, course)
	json.NewEncoder(w).Encode(course)
	return
}

func updateOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update one course")
	w.Header().Set("Content-Type", "application/json")

	//1st grab id from params
	params := mux.Vars(r)

	//loop get the value, remove, and add again new value with same id
	for index, course := range courses {
		if course.CourseId == params["id"] {
			//remove the item from slice
			courses = append(courses[:index], courses[index+1:]...)

			//now grab the json from req body decode it and create data
			//and insert it into the slice
			var course Course
			_ = json.NewDecoder(r.Body).Decode(&course)
			//now id has to be the same as params id
			course.CourseId = params["id"]

			//now add this new course
			courses = append(courses, course)
			json.NewEncoder(w).Encode(course)
			return
		}
	}
}

func deleteOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete one course")
	w.Header().Set("Content-Type", "application/json")

	//1st grab the id
	params := mux.Vars(r)

	//loop through get the data & remove it
	for index, course := range courses {
		if course.CourseId == params["id"] {
			courses = append(courses[:index], courses[index+1:]...)
			json.NewEncoder(w).Encode("Deleted Successfully")
			break
		}
	}
}
