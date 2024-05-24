package Controllers

import (
	db "SchoolProject/Config"
	"SchoolProject/Models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

func CreateSchool(c *fiber.Ctx) error {
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
				"message": "School Name is required",
			})
	}
	if data["address"] == "" {
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "School Address is required",
			})
	}
	if data["city"] == "" {
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "School City is required",
			})
	}
	//now saving school to db
	school := Models.School{
		Name:      data["name"],
		Address:   data["address"],
		City:      data["city"],
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	db.DB.Create(&school)
	return c.Status(200).JSON(
		fiber.Map{
			"success": true,
			"message": "Success",
			"data":    school,
		})
}

func UpdateSchool(c *fiber.Ctx) error {
	schoolId := c.Params("schoolId")
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid data",
		})
	}

	var school Models.School
	if err := db.DB.Where("id = ?", schoolId).First(&school).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "School not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Error finding school",
		})
	}

	// Updating school fields
	if name, ok := data["name"]; ok {
		school.Name = name
	}
	if address, ok := data["address"]; ok {
		school.Address = address
	}
	if city, ok := data["city"]; ok {
		school.City = city
	}

	school.UpdatedAt = time.Now()

	// Save changes
	if err := db.DB.Save(&school).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Error updating school",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "School updated successfully",
		"data":    school,
	})
}

func DeleteSchool(c *fiber.Ctx) error {
	schoolId := c.Params("schoolId")
	var school Models.School
	db.DB.Where("id=?", schoolId).Find(&school)

	if school.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "School not found",
		})
	}
	db.DB.Unscoped().Where("id=?", schoolId).Delete(&school)
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "School deleted successfully",
	})
}

func GetSchoolDetails(c *fiber.Ctx) error {
	schoolId := c.Params("schoolId")
	var school Models.School
	db.DB.Select("id,name,address,city,created_at,updated_at").Where("id =?", schoolId).First(&school)
	schoolData := make(map[string]interface{})
	schoolData["schoolId"] = school.ID
	schoolData["name"] = school.Name
	schoolData["address"] = school.Address
	schoolData["city"] = school.City
	schoolData["createdAt"] = school.CreatedAt
	schoolData["updatedAt"] = school.UpdatedAt

	//what if there is no school or school id not provided
	if school.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Error fetching school",
			"error":   map[string]interface{}{},
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "success",
		"data":    schoolData,
	})
}

func SchoolsList(c *fiber.Ctx) error {

	var schools []Models.School

	db.DB.Find(&schools)

	var schoolsData []map[string]interface{}
	for _, school := range schools {
		schoolData := map[string]interface{}{
			"id":   school.ID,
			"name": school.Name,
		}
		schoolsData = append(schoolsData, schoolData)
	}

	var count int64
	db.DB.Model(&Models.School{}).Count(&count)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "School list API",
		"data":    schoolsData,
		"total":   count,
	})

}

func GetClasses(c *fiber.Ctx) error {

	schoolId := c.Params("schoolId")
	var school Models.School

	db.DB.Select("id").Where("id = ?", schoolId).First(&school)

	// Check if school ID is valid
	if school.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "School not found",
		})
	}

	// Fetch the classes for the school
	var classes []Models.Class
	if err := db.DB.Select("name").Where("school_id = ?", schoolId).Find(&classes).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Error fetching classes",
		})
	}

	classesNames := make([]string, len(classes))
	for i, class := range classes {
		classesNames[i] = class.Name
	}
	var count int64
	db.DB.Model(&Models.Class{}).Count(&count) // Count total classes

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"success": true,
			"message": "Class list API",
			"data":    classesNames,
			"total":   count,
		})

}

func GetTeachersInSchool(c *fiber.Ctx) error {

	schoolId := c.Params("schoolId")
	var school Models.School

	db.DB.Select("id").Where("id = ?", schoolId).First(&school)

	// Check if school ID is valid
	if school.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "School not found",
		})
	}

	// Fetch the classes for the school
	var teachers []Models.Teacher
	if err := db.DB.Select("name").Where("school_id = ?", schoolId).Find(&teachers).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Error fetching teachers",
		})
	}

	teachersNames := make([]string, len(teachers))
	for i, teacher := range teachers {
		teachersNames[i] = teacher.Name
	}
	var count int64
	db.DB.Model(&Models.Teacher{}).Count(&count) // Count total teachers

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"success": true,
			"message": "Class list API",
			"data":    teachersNames,
			"total":   count,
		})

}

func GetStudentsInSchool(c *fiber.Ctx) error {
	schoolId := c.Params("schoolId")
	var school Models.School

	// Check if school ID is valid
	if err := db.DB.Select("id").Where("id = ?", schoolId).First(&school).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "School not found",
		})
	}

	var students []Models.Student
	if err := db.DB.Joins("JOIN classes ON classes.id = students.class_id").
		Where("classes.school_id = ?", schoolId).
		Find(&students).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch students",
		})
	}

	studentNames := make([]string, len(students))
	for i, student := range students {
		studentNames[i] = student.Name
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Class list API",
		"data":    studentNames,
		"total":   len(students),
	})
}
