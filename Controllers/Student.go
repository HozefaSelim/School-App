package Controllers

import (
	db "SchoolProject/Config"
	"SchoolProject/Models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func CreateStudent(c *fiber.Ctx) error {
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "Invalid data",
			})
	}
	if data["name"] == "" {
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "Student Name is required",
			})
	}
	if data["phone"] == "" {
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "Student Phone is required",
			})
	}
	if data["gender"] == "" {
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "Student Gender is required",
			})
	}

	if data["class_id"] == "" {
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "class_id is required",
			})
	}
	classId, err := strconv.ParseUint(data["class_id"], 10, 64)
	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "class_id must be Uint ",
			})
	}
	var class Models.Class
	if err := db.DB.Where("id = ?", classId).First(&class).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Class not found",
		})
	}

	// now saving class to db
	student := Models.Student{
		Name:      data["name"],
		Phone:     data["phone"],
		Gender:    data["gender"],
		ClassID:   uint(classId),
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	db.DB.Create(&student)
	return c.Status(200).JSON(
		fiber.Map{
			"success": true,
			"message": "Success",
			"data":    student,
		})
}

func UpdateStudent(c *fiber.Ctx) error {
	studentId := c.Params("studentId")
	var data map[string]string

	// Parse the request body
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid data",
		})
	}

	// Find the student by ID
	var student Models.Student
	if err := db.DB.Where("id = ?", studentId).First(&student).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "Student not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Error finding student",
		})
	}

	// Updating student fields
	if name, ok := data["name"]; ok {
		student.Name = name
	}
	if phone, ok := data["phone"]; ok {
		student.Phone = phone
	}
	if gender, ok := data["gender"]; ok {
		student.Gender = gender
	}
	if classIdStr, ok := data["class_id"]; ok {
		classId, err := strconv.ParseUint(classIdStr, 10, 64)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Invalid class_id",
			})
		}
		student.ClassID = uint(classId)
	}

	student.UpdatedAt = time.Now()

	// Save changes
	if err := db.DB.Save(&student).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Error updating student",
		})
	}

	// Prepare the response map with only necessary student information
	response := map[string]interface{}{
		"id":         student.ID,
		"name":       student.Name,
		"phone":      student.Phone,
		"gender":     student.Gender,
		"class_id":   student.ClassID,
		"updated_at": student.UpdatedAt,
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Student updated successfully",
		"data":    response,
	})
}

func DeleteStudent(c *fiber.Ctx) error {
	studentId := c.Params("studentId")
	var student Models.Student
	db.DB.Where("id=?", studentId).Find(&student)

	if student.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Student not found",
		})
	}
	db.DB.Unscoped().Where("id=?", studentId).Delete(&student)
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Student deleted successfully",
	})
}

func GetStudentDetails(c *fiber.Ctx) error {
	studentId := c.Params("studentId")
	var student Models.Student
	db.DB.Select("id,name,phone,gender,class_id,created_at,updated_at").Where("id =?", studentId).First(&student)
	studentData := make(map[string]interface{})
	studentData["Id"] = student.ID
	studentData["name"] = student.Name
	studentData["gender"] = student.Gender
	studentData["phone"] = student.Phone
	studentData["class_id"] = student.ClassID
	studentData["createdAt"] = student.CreatedAt
	studentData["updatedAt"] = student.UpdatedAt

	//what if there is no student or student id not provided
	if student.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Error fetching student",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "success",
		"data":    studentData,
	})
}

func StudentsList(c *fiber.Ctx) error {

	var students []Models.Student

	db.DB.Find(&students)

	var studentsData []map[string]interface{}
	for _, student := range students {
		studentData := map[string]interface{}{
			"id":   student.ID,
			"name": student.Name,
		}
		studentsData = append(studentsData, studentData)
	}

	var count int64
	db.DB.Model(&Models.Student{}).Count(&count)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Student list API",
		"data":    studentsData,
		"total":   count,
	})

}

func AddStudentToSubject(c *fiber.Ctx) error {
	studentId := c.Params("studentId")

	var data map[string]string
	// Parse the body to get the subject ID
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid data",
		})
	}

	// Find the student by ID
	var student Models.Student
	if err := db.DB.Preload("Subjects").First(&student, studentId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "Student not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Error finding student",
		})
	}

	// Parse and find the subject by ID
	subjectId, err := strconv.ParseUint(data["subject_id"], 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Subject ID must be a valid uint",
		})
	}

	var subject Models.Subject
	if err := db.DB.First(&subject, subjectId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "Subject not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Error finding subject",
		})
	}

	// Add the subject to the student's subjects
	if err := db.DB.Model(&student).Association("Subjects").Append(&subject); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Error adding subject to student",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Subject added to student successfully",
	})
}

func GetStudentSubjects(c *fiber.Ctx) error {

	studentId := c.Params("studentId")

	// Find the student by ID and preload subjects
	var student Models.Student
	if err := db.DB.Preload("Subjects").First(&student, studentId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "Student not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Error finding student",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    student.Subjects,
	})
}
