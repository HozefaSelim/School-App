package Controllers

import (
	db "SchoolProject/Config"
	"SchoolProject/Models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

func CreateSubject(c *fiber.Ctx) error {
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
				"message": "Subject Name is required",
			})
	}

	// now saving class to db
	subject := Models.Subject{
		Name:      data["name"],
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	db.DB.Create(&subject)
	return c.Status(200).JSON(
		fiber.Map{
			"success": true,
			"message": "Success",
			"data":    subject,
		})
}

func UpdateSubject(c *fiber.Ctx) error {

	subjectId := c.Params("subjectId")
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid data",
		})
	}

	var subject Models.Subject
	if err := db.DB.Where("id = ?", subjectId).First(&subject).Error; err != nil {
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

	// Updating subject fields
	if name, ok := data["name"]; ok {
		subject.Name = name
	}

	subject.UpdatedAt = time.Now()

	// Save changes
	if err := db.DB.Save(&subject).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Error updating subject",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Subject updated successfully",
		"data":    subject,
	})
}

func DeleteSubject(c *fiber.Ctx) error {
	subjectId := c.Params("subjectId")
	var subject Models.Subject
	db.DB.Where("id=?", subjectId).Find(&subject)

	if subject.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Subject not found",
		})
	}
	db.DB.Unscoped().Where("id=?", subjectId).Delete(&subject)
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Subject deleted successfully",
	})
}

func GetSubjectDetails(c *fiber.Ctx) error {
	subjectId := c.Params("subjectId")
	var subject Models.Subject
	db.DB.Select("id,name,created_at,updated_at").Where("id =?", subjectId).First(&subject)
	subjectData := make(map[string]interface{})
	subjectData["subjectId"] = subject.ID
	subjectData["name"] = subject.Name
	subjectData["createdAt"] = subject.CreatedAt
	subjectData["updatedAt"] = subject.UpdatedAt

	//what if there is no subject or subject id not provided
	if subject.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Error fetching subject",
			"error":   map[string]interface{}{},
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "success",
		"data":    subjectData,
	})
}

func SubjectsList(c *fiber.Ctx) error {

	var subjects []Models.Subject

	db.DB.Find(&subjects)

	var subjectsData []map[string]interface{}
	for _, subject := range subjects {
		subjectData := map[string]interface{}{
			"id":   subject.ID,
			"name": subject.Name,
		}
		subjectsData = append(subjectsData, subjectData)
	}

	var count int64
	db.DB.Model(&Models.Subject{}).Count(&count)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Subject list API",
		"data":    subjectsData,
		"total":   count,
	})

}

func GetTeachers(c *fiber.Ctx) error {
	subjectId := c.Params("subjectId")

	// Find the subject by ID to check if it exists
	var subject Models.Subject
	if err := db.DB.Select("id").First(&subject, subjectId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "Subject not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Error retrieving subject",
		})
	}

	// Fetch the teachers for the subject
	var teachers []Models.Teacher
	if err := db.DB.Select("name").Where("subject_id = ?", subjectId).Find(&teachers).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Error fetching teachers",
		})
	}

	// Extract teacher names
	teacherNames := make([]string, len(teachers))
	for i, teacher := range teachers {
		teacherNames[i] = teacher.Name
	}

	// Count teachers for the specific subject
	var totalTeachers int64
	db.DB.Model(&Models.Teacher{}).Where("subject_id = ?", subjectId).Count(&totalTeachers)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Teacher list API",
		"data":    teacherNames,
		"total":   totalTeachers,
	})
}

func GetTeachersAndStudentsBySubject(c *fiber.Ctx) error {
	subjectID := c.Params("subjectId")

	var teachers []Models.Teacher
	var students []Models.Student

	if err := db.DB.Where("subject_id = ?", subjectID).Find(&teachers).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch teachers",
		})
	}

	if err := db.DB.Joins("JOIN student_subjects ON student_subjects.student_id = students.id").
		Where("student_subjects.subject_id = ?", subjectID).
		Find(&students).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch students",
		})
	}

	teacherNames := make([]string, len(teachers))
	for i, teacher := range teachers {
		teacherNames[i] = teacher.Name
	}

	studentNames := make([]string, len(students))
	for i, student := range students {
		studentNames[i] = student.Name
	}

	return c.JSON(fiber.Map{
		"teacherNames": teacherNames,
		"studentNames": studentNames,
	})
}
