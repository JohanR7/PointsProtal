package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Routes

func main() {
	// Initialize MongoDB
	if err := initMongoDB(); err != nil {
		panic(err)
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// User routes
	r.POST("/signup", Signup)
	r.POST("/login", Login)
	r.POST("/teachers", CreateTeacher)
	r.POST("/events", CreateEvent)
	r.POST("/roles/:eventid", CreateRole)
	// Event routes
	r.GET("/events", ListEvents)
	// r.GET("/events/:id", GetEventByID)
	r.GET("/teachers", ListTeachers)

	r.PUT("/events/:id", UpdateEvent)

	// Role routes

	r.GET("/events/:id/roles", GetRolesByEventID)

	// Teacher routes

	r.GET("/teachers/top", GetTopTeachers)

	// Assignment routes
	r.POST("/assignments", AssignTeacherToRole)

	r.DELETE("/delete-role-assignment", DeleteRoleAssignment)

	// GET: Get all assignments for a specific teacher
	r.GET("/teacher-assignments/:id", GetTeacherAssignments)

	// GET: Get all assignments for a specific role
	r.GET("/role-assignments/:id", GetRoleAssignments)

	r.DELETE("/event", DeleteEvent)

	r.GET("/event/:id/roles", GetRolesByEventID)

	r.GET("/teacher/:teacherid/event/:eventid/roles", GetTeacherRolesInEvent)

	r.GET("/events/assigned-teachers/:eventid", GetAssignedTeachersForEvent)

	r.Run(":8080")
}
