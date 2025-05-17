package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// if user.Role != "admin" && user.Role != "teacher" {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
	// 	return
	// }
	user.Role = "teacher"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userCollectionRef := db.Collection(userCollection)
	teacherCollectionRef := db.Collection(teacherCollection)

	// Check if email already exists in users
	count, err := userCollectionRef.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already in use"})
		return
	}

	// Check if a teacher exists with the same email
	var existingTeacher Teacher
	err = teacherCollectionRef.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingTeacher)
	if err == nil {
		// Teacher found, use Teacher's UserID as the new User's _id
		user.ID = existingTeacher.UserID
		user.UserID = existingTeacher.UserID
	} else {
		// No matching teacher, generate new ObjectID
		user.ID = primitive.NewObjectID()
		user.UserID = user.ID
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	// Insert user with assigned ID
	_, err = userCollectionRef.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user_id": user.ID,
	})
}

// Login handler
func Login(c *gin.Context) {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type LoginResponse struct {
		Message string             `json:"message"`
		Role    string             `json:"role,omitempty"`
		Name    string             `json:"name,omitempty"`
		UserID  primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	collection := db.Collection(userCollection)
	err := collection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Message: "Login successful",
		Name:    user.Name,
		Role:    user.Role,
		UserID:  user.UserID,
	})
}

// CreateEvent handler
func CreateEvent(c *gin.Context) {
	var event Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.Collection(eventCollection)
	event.ID = primitive.NewObjectID()
	event.EventID = event.ID
	_, err := collection.InsertOne(ctx, event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, event)
}

// ListEvents handler
func ListEvents(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.Collection(eventCollection)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var events []Event
	if err := cursor.All(ctx, &events); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

// GetEventByID handler
func GetEventByID(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.Collection(eventCollection)
	var event Event
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&event)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// UpdateEvent handler
func UpdateEvent(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var event Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.Collection(eventCollection)
	update := bson.M{
		"$set": bson.M{
			"name":        event.Name,
			"start_date":  event.StartDate,
			"start_time":  event.StartTime,
			"end_date":    event.EndDate,
			"end_time":    event.EndTime,
			"description": event.Description,
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event updated successfully"})
}

func GetTopTeachers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// MongoDB aggregation pipeline to get top teachers
	pipeline := mongo.Pipeline{
		{
			{"$lookup", bson.M{
				"from":         teacherAssignmentCollection,
				"localField":   "_id",
				"foreignField": "teacher_id",
				"as":           "assignments",
			}},
		},
		{
			{"$unwind", bson.M{
				"path":                       "$assignments",
				"preserveNullAndEmptyArrays": true,
			}},
		},
		{
			{"$lookup", bson.M{
				"from":         roleCollection,
				"localField":   "assignments.role_id",
				"foreignField": "_id",
				"as":           "role",
			}},
		},
		{
			{"$unwind", bson.M{
				"path":                       "$role",
				"preserveNullAndEmptyArrays": true,
			}},
		},
		{
			{"$group", bson.M{
				"_id":          "$_id",
				"teacher_name": bson.M{"$first": "$name"},
				"total_points": bson.M{"$sum": bson.M{
					"$ifNull": []interface{}{"$role.point", 0},
				}},
			}},
		},
		{
			{"$sort", bson.M{"total_points": -1}},
		},
		{
			{"$limit", 10},
		},
	}

	collection := db.Collection(teacherCollection)
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	type TopTeacher struct {
		ID     primitive.ObjectID `json:"teacher_id" bson:"_id"`
		Name   string             `json:"teacher_name" bson:"teacher_name"`
		Points int                `json:"points" bson:"total_points"`
	}

	var teachers []TopTeacher
	if err := cursor.All(ctx, &teachers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teachers)
}

func CreateRole(c *gin.Context) {
	eventID := c.Param("eventid")

	var role Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert eventID string param to ObjectID
	oid, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	roleCollection := db.Collection("roles")
	eventCollection := db.Collection("events")

	// Retrieve the Event document to get event name BEFORE inserting role
	var event Event
	err = eventCollection.FindOne(ctx, bson.M{"_id": oid}).Decode(&event)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	// Assign all fields before inserting
	role.ID = primitive.NewObjectID()
	role.EventID = oid
	role.EventName = event.Name // <- now it will get saved
	role.RoleID = role.ID

	// Insert the new role document into roles collection
	_, err = roleCollection.InsertOne(ctx, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert role: " + err.Error()})
		return
	}

	// Update the event document: push role object with id and name to the roles array
	update := bson.M{
		"$push": bson.M{
			"roles": bson.M{
				"id":   role.ID,
				"name": role.Name,
			},
		},
	}
	_, err = eventCollection.UpdateOne(ctx, bson.M{"_id": oid}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event with role: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}

func CreateTeacher(c *gin.Context) {
	var teacher Teacher
	if err := c.ShouldBindJSON(&teacher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.Collection(teacherCollection)
	teacher.ID = primitive.NewObjectID()
	teacher.UserID = teacher.ID
	_, err := collection.InsertOne(ctx, teacher)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update user if email matches
	userCollectionRef := db.Collection(userCollection)
	filter := bson.M{"email": teacher.Email}
	update := bson.M{"$set": bson.M{"user_id": teacher.ID}} // Add new field
	_, _ = userCollectionRef.UpdateOne(ctx, filter, update)

	c.JSON(http.StatusOK, teacher)
}

// ListTeachers handler
func ListTeachers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// MongoDB aggregation pipeline to get teachers with department info
	pipeline := mongo.Pipeline{
		{
			{"$lookup", bson.M{
				"from":         departmentCollection,
				"localField":   "department_id",
				"foreignField": "_id",
				"as":           "department",
			}},
		},
		{
			{"$unwind", bson.M{
				"path":                       "$department",
				"preserveNullAndEmptyArrays": true,
			}},
		},
		{
			{"$project", bson.M{
				"_id":             1,
				"name":            1,
				"email":           1,
				"profile_photo":   1,
				"point":           1,
				"department_name": "$departmentname",
			}},
		},
	}

	collection := db.Collection(teacherCollection)
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	type TeacherWithDepartment struct {
		ID             primitive.ObjectID `json:"id" bson:"_id"`
		Name           string             `json:"name" bson:"name"`
		Email          string             `json:"email" bson:"email"`
		ProfilePhoto   string             `json:"profile_photo" bson:"profile_photo"`
		DepartmentName string             `json:"department_name" bson:"department_name"`
		Point          int                `json:"point" bson:"point"`
	}

	var teachers []TeacherWithDepartment
	if err := cursor.All(ctx, &teachers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teachers)
}

// AssignTeacherToRole assigns a teacher to a role and updates their points
func AssignTeacherToRole(c *gin.Context) {
	type AssignmentRequest struct {
		TeacherID string `json:"teacher_id" binding:"required"`
		RoleID    string `json:"role_id" binding:"required"`
		EventID   string `json:"event_id" binding:"required"`
	}

	var req AssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert string IDs to ObjectIDs
	teacherID, err := primitive.ObjectIDFromHex(req.TeacherID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid teacher ID format"})
		return
	}

	roleID, err := primitive.ObjectIDFromHex(req.RoleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID format"})
		return
	}

	eventID, err := primitive.ObjectIDFromHex(req.EventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID format"})
		return
	}

	// Get the role details to obtain points and event name
	roleCollection := db.Collection(roleCollection)
	var role Role
	err = roleCollection.FindOne(ctx, bson.M{"_id": roleID}).Decode(&role)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	// Get the event name
	eventCollection := db.Collection(eventCollection)
	var event Event
	err = eventCollection.FindOne(ctx, bson.M{"_id": eventID}).Decode(&event)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	// Check if the teacher exists
	teacherCollection := db.Collection(teacherCollection)
	var teacher Teacher
	err = teacherCollection.FindOne(ctx, bson.M{"_id": teacherID}).Decode(&teacher)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher not found"})
		return
	}

	// Check if this assignment already exists
	assignmentCollection := db.Collection(teacherAssignmentCollection)
	count, err := assignmentCollection.CountDocuments(ctx, bson.M{
		"teacher_id": teacherID,
		"role_id":    roleID,
		"event_id":   eventID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Teacher is already assigned to this role in this event"})
		return
	}

	// Check if the role has reached its head count limit
	assignedCount, err := assignmentCollection.CountDocuments(ctx, bson.M{"role_id": roleID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if int(assignedCount) >= role.HeadCount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role has reached its maximum head count"})
		return
	}

	// Create the assignment using the exact Assignment struct
	assignment := Assignment{
		ID:        primitive.NewObjectID(),
		EventID:   eventID,
		EventName: event.Name,
		TeacherID: teacherID,
		RoleID:    roleID,
		RoletName: role.Name,
	}

	_, err = assignmentCollection.InsertOne(ctx, assignment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create assignment"})
		return
	}

	// Update teacher's points
	_, err = teacherCollection.UpdateOne(
		ctx,
		bson.M{"_id": teacherID},
		bson.M{"$inc": bson.M{"point": role.Point}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update teacher points"})
		return
	}

	// Add role reference to teacher's Assginedteachers array using the RoleRef struct
	roleRef := RoleRef1{
		ID:            roleID,
		RoleName:      role.Name,
		TeacherleName: teacher.Name,
		Assignment_ID: assignment.ID,
	}

	_, err = eventCollection.UpdateOne(
		ctx,
		bson.M{"_id": eventID},
		bson.M{"$push": bson.M{"assginedteachers": roleRef}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update teacher's assigned roles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Teacher assigned to role successfully",
		"assignment": assignment,
	})
}

// DeleteRoleAssignment removes a teacher's role assignment with optional point handling
func DeleteRoleAssignment(c *gin.Context) {
	type DeleteAssignmentRequest struct {
		AssignmentID string `json:"assignment_id" binding:"required"`
		DeductPoints bool   `json:"deduct_points"` // Whether to deduct points from teacher
	}

	var req DeleteAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	assignmentID, err := primitive.ObjectIDFromHex(req.AssignmentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID format"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find the assignment first to get the role ID and teacher ID
	assignmentCollection := db.Collection(teacherAssignmentCollection)
	var assignment Assignment
	err = assignmentCollection.FindOne(ctx, bson.M{"_id": assignmentID}).Decode(&assignment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// If we need to deduct points, we need to get the role's point value
	if req.DeductPoints {
		// Get the role to determine how many points to deduct
		roleCollection := db.Collection(roleCollection)
		var role Role
		err = roleCollection.FindOne(ctx, bson.M{"_id": assignment.RoleID}).Decode(&role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Role not found"})
			return
		}

		// Deduct points from the teacher
		teacherCollection := db.Collection(teacherCollection)
		_, err = teacherCollection.UpdateOne(
			ctx,
			bson.M{"_id": assignment.TeacherID},
			bson.M{"$inc": bson.M{"point": -role.Point}},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update teacher points"})
			return
		}
	}

	// Remove the role reference from teacher's assginedteachers array
	// Match using the exact field name from the Teacher struct
	eventCollection := db.Collection(eventCollection)
	_, err = eventCollection.UpdateOne(
		ctx,
		bson.M{"_id": assignment.TeacherID},
		bson.M{"$pull": bson.M{"assginedteachers": bson.M{"id": assignment.RoleID}}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update teacher's assigned roles"})
		return
	}

	// Delete the assignment
	_, err = assignmentCollection.DeleteOne(ctx, bson.M{"_id": assignmentID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete assignment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Role assignment deleted successfully",
		"deducted_points": req.DeductPoints,
	})
}

// GetTeacherAssignments retrieves all role assignments for a specific teacher
func GetTeacherAssignments(c *gin.Context) {
	teacherID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(teacherID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid teacher ID format"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.Collection(teacherAssignmentCollection)
	cursor, err := collection.Find(ctx, bson.M{"teacher_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var assignments []Assignment
	if err := cursor.All(ctx, &assignments); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assignments)
}

// GetRoleAssignments retrieves all teacher assignments for a specific role
func GetRoleAssignments(c *gin.Context) {
	roleID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(roleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID format"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.Collection(teacherAssignmentCollection)
	cursor, err := collection.Find(ctx, bson.M{"role_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var assignments []Assignment
	if err := cursor.All(ctx, &assignments); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assignments)
}

//

func DeleteEvent(c *gin.Context) {
	type DeleteEventRequest struct {
		EventID      string `json:"event_id" binding:"required"`
		DeductPoints bool   `json:"deduct_points"` // Whether to deduct points from teachers
	}

	var req DeleteEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	eventID, err := primitive.ObjectIDFromHex(req.EventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID format"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the event details to confirm it exists
	eventCollection := db.Collection(eventCollection)
	var event Event
	err = eventCollection.FindOne(ctx, bson.M{"_id": eventID}).Decode(&event)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Find all roles associated with the event
	roleCollection := db.Collection(roleCollection)
	roleCursor, err := roleCollection.Find(ctx, bson.M{"event_id": eventID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query roles"})
		return
	}
	defer roleCursor.Close(ctx)

	var roles []Role
	if err := roleCursor.All(ctx, &roles); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse roles"})
		return
	}

	// Process teacher assignments for each role
	assignmentCollection := db.Collection(teacherAssignmentCollection)
	teacherCollection := db.Collection(teacherCollection)

	for _, role := range roles {
		// Find all assignments for this role
		assignmentCursor, err := assignmentCollection.Find(ctx, bson.M{"role_id": role.ID})
		if err != nil {
			continue // Skip if error
		}

		var assignments []Assignment
		if err := assignmentCursor.All(ctx, &assignments); err != nil {
			assignmentCursor.Close(ctx)
			continue // Skip if error
		}
		assignmentCursor.Close(ctx)

		// If deducting points is requested, update each teacher's points
		if req.DeductPoints {
			for _, assignment := range assignments {
				_, _ = teacherCollection.UpdateOne(
					ctx,
					bson.M{"_id": assignment.TeacherID},
					bson.M{"$inc": bson.M{"point": -role.Point}},
				)

				// Remove role reference from event's assginedteachers array
				_, _ = eventCollection.UpdateOne(
					ctx,
					bson.M{"_id": eventID},
					bson.M{"$pull": bson.M{"assginedteachers": bson.M{"assignment_id": assignment.ID}}},
				)
			}
		}

		// Delete all assignments for this role
		_, _ = assignmentCollection.DeleteMany(ctx, bson.M{"role_id": role.ID})
	}

	// Delete all roles associated with the event
	_, err = roleCollection.DeleteMany(ctx, bson.M{"event_id": eventID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete roles"})
		return
	}

	// Finally delete the event itself
	_, err = eventCollection.DeleteOne(ctx, bson.M{"_id": eventID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Event and all associated data deleted successfully",
		"deducted_points": req.DeductPoints,
		"event_name":      event.Name,
	})
}

// GetRolesByEventID retrieves all roles for a specific event
func GetRolesByEventID(c *gin.Context) {
	eventID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID format"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	roleCollection := db.Collection(roleCollection)
	cursor, err := roleCollection.Find(ctx, bson.M{"event_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var roles []Role
	if err := cursor.All(ctx, &roles); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, roles)
}

// GetTeacherRolesInEvent retrieves all roles assigned to a teacher in a specific event
func GetTeacherRolesInEvent(c *gin.Context) {
	teacherID := c.Param("teacherid")
	eventID := c.Param("eventid")

	teacherObjID, err := primitive.ObjectIDFromHex(teacherID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid teacher ID format"})
		return
	}

	eventObjID, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID format"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find all assignments for this teacher in this event
	assignmentCollection := db.Collection(teacherAssignmentCollection)

	// MongoDB aggregation pipeline to get detailed role information
	pipeline := mongo.Pipeline{
		{
			{"$match", bson.M{
				"teacher_id": teacherObjID,
				"event_id":   eventObjID,
			}},
		},
		{
			{"$lookup", bson.M{
				"from":         roleCollection,
				"localField":   "role_id",
				"foreignField": "_id",
				"as":           "role_details",
			}},
		},
		{
			{"$unwind", bson.M{
				"path":                       "$role_details",
				"preserveNullAndEmptyArrays": false,
			}},
		},
		{
			{"$project", bson.M{
				"_id":              1,
				"event_id":         1,
				"event_name":       "$eventname",
				"teacher_id":       1,
				"role_id":          1,
				"role_name":        "$role_details.name",
				"role_description": "$role_details.description",
				"role_point":       "$role_details.point",
			}},
		},
	}

	cursor, err := assignmentCollection.Aggregate(ctx, pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	type TeacherRoleAssignment struct {
		ID              primitive.ObjectID `json:"id" bson:"_id"`
		EventID         primitive.ObjectID `json:"event_id" bson:"event_id"`
		EventName       string             `json:"event_name" bson:"event_name"`
		TeacherID       primitive.ObjectID `json:"teacher_id" bson:"teacher_id"`
		RoleID          primitive.ObjectID `json:"role_id" bson:"role_id"`
		RoleName        string             `json:"role_name" bson:"role_name"`
		RoleDescription string             `json:"role_description" bson:"role_description"`
		RolePoint       int                `json:"role_point" bson:"role_point"`
	}

	var assignments []TeacherRoleAssignment
	if err := cursor.All(ctx, &assignments); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assignments)
}

func GetAssignedTeachersForEvent(c *gin.Context) {
	eventIDParam := c.Param("eventid")
	eventID, err := primitive.ObjectIDFromHex(eventIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	assignmentCollection := db.Collection("assignments")

	// Fetch all assignments for this event
	cursor, err := assignmentCollection.Find(ctx, bson.M{"event_id": eventID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch assignments"})
		return
	}
	defer cursor.Close(ctx)

	var assignments []Assignment
	if err = cursor.All(ctx, &assignments); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode assignments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"assignments": assignments,
	})
}
