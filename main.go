package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Patient model
type Patient struct {
	ID             uint   `json:"id" gorm:"primaryKey"`
	Name           string `json:"name"`
	Age            int    `json:"age"`
	MedicalHistory string `json:"medical_history"`
}

var db *gorm.DB
var err error

func main() {
	// Initialize the database
	db, err = gorm.Open(sqlite.Open("patients.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database!")
	}

	// Migrate the schema
	db.AutoMigrate(&Patient{})

	// Initialize Gin router
	r := gin.Default()

	// Routes
	r.POST("/patients", createPatient)
	r.GET("/patients", getPatients)
	r.GET("/patients/:id", getPatientByID)
	r.PUT("/patients/:id", updatePatient)
	r.DELETE("/patients/:id", deletePatient)

	// Run the server
	r.Run(":8080")
}

// Create a new patient
func createPatient(c *gin.Context) {
	var patient Patient
	if err := c.ShouldBindJSON(&patient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Create(&patient)
	c.JSON(http.StatusCreated, patient)
}

// Get all patients
func getPatients(c *gin.Context) {
	var patients []Patient
	db.Find(&patients)
	c.JSON(http.StatusOK, patients)
}

// Get a patient by ID
func getPatientByID(c *gin.Context) {
	id := c.Param("id")
	var patient Patient
	if err := db.First(&patient, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}
	c.JSON(http.StatusOK, patient)
}

// Update a patient
func updatePatient(c *gin.Context) {
	id := c.Param("id")
	var updatedData Patient
	if err := db.First(&updatedData, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Model(&updatedData).Updates(updatedData)
	c.JSON(http.StatusOK, updatedData)
}

// Delete a patient
func deletePatient(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&Patient{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Patient deleted successfully"})
}
