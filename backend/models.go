package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoDB client and database
var client *mongo.Client
var db *mongo.Database

// Initialize MongoDB connection

// Collection names
const (
	userCollection              = "users"
	eventCollection             = "events"
	teacherCollection           = "teachers"
	roleCollection              = "roles"
	teacherAssignmentCollection = "teacherAssignments"
	departmentCollection        = "departments"
)

// User struct
type User struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password,omitempty" bson:"password"`
	Role     string             `json:"role" bson:"role"`
	UserID   primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
}

// Event struct
type Event struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	EventID     primitive.ObjectID `json:"event_id,omitempty" bson:"event_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	StartDate   string             `json:"start_date" bson:"start_date"`
	StartTime   string             `json:"start_time" bson:"start_time"`
	EndDate     string             `json:"end_date" bson:"end_date"`
	EndTime     string             `json:"end_time" bson:"end_time"`
	Description string             `json:"description" bson:"description"`
	// Roles       []primitive.ObjectID `json:"roles,omitempty" bson:"roles,omitempty"`
	Roles            []RoleRef `json:"roles,omitempty" bson:"roles,omitempty"`
	Assginedteachers []RoleRef `json:"assginedteachers,omitempty" bson:"assginedteachers,omitempty"`
}

// Teacher struct
type Teacher struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name           string             `json:"name" bson:"name"`
	Email          string             `json:"email" bson:"email"`
	Departmentname string             `json:"departmentname" bson:"departmentname"`
	ProfilePhoto   string             `json:"profile_photo" bson:"profile_photo"`
	Point          int                `json:"point,omitempty" bson:"point,omitempty"`
	// UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	UserID           primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Assginedteachers []RoleRef          `json:"assginedteachers,omitempty" bson:"assginedteachers,omitempty"`
}
type RoleRef struct {
	ID   primitive.ObjectID `json:"id" bson:"id"`
	Name string             `json:"name" bson:"name"`
}

type RoleRef1 struct {
	ID            primitive.ObjectID `json:"id" bson:"id"`
	RoleName      string             `json:"rolename" bson:"rolename"`
	TeacherleName string             `json:"teachername" bson:"teachername"`
	Assignment_ID primitive.ObjectID `json:"AssignmentID,omitempty" bson:"AssignmentID,omitempty"`
}

// Role struct
type Role struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	RoleID    primitive.ObjectID `json:"roleid,omitempty" bson:"roleid,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Point     int                `json:"point" bson:"point"`
	HeadCount int                `json:"head_count" bson:"head_count"`
	EventID   primitive.ObjectID `json:"event_id" bson:"event_id"`
	// EventName string             `json:"eventname,omitempty" bson:"evenetname,omitempty"`
	EventName string `json:"eventname,omitempty" bson:"eventname,omitempty"`
}

// Assignment struct
type Assignment struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	AssignmentID primitive.ObjectID `json:"assignmentid,omitempty" bson:"assignmentid,omitempty"`
	EventID      primitive.ObjectID `json:"event_id" bson:"event_id"`
	EventName    string             `json:"eventname" bson:"eventname"`
	TeacherID    primitive.ObjectID `json:"teacher_id" bson:"teacher_id"`
	RoleID       primitive.ObjectID `json:"role_id" bson:"role_id"`
	RoletName    string             `json:"roletname" bson:"roletname"`
}

// Department struct
