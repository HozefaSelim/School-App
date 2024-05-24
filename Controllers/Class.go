package Controllers

import (
	db "SchoolProject/Config"
	"SchoolProject/Models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func CreateClass(c *fiber.Ctx) error {
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
				"message": "Class Name is required",
			})
	}
	if data["school_id"] == "" {
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "school_id is required",
			})
	}
	schoolId, err := strconv.ParseUint(data["school_id"], 10, 64)
	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "School_id must be Uint ",
			})
	}
	var school Models.School
	if err := db.DB.Where("id = ?", schoolId).First(&school).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "School not found",
		})
	}

	// now saving class to db
	class := Models.Class{
		Name:      data["name"],
		SchoolID:  uint(schoolId),
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	db.DB.Create(&class)
	return c.Status(200).JSON(
		fiber.Map{
			"success": true,
			"message": "Success",
			"data":    class,
		})
}

func UpdateClass(c *fiber.Ctx) error {
	classId := c.Params("classId")
	var data map[string]string

	// Parse the request body
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid data",
		})
	}

	// Find the class by ID
	var class Models.Class
	if err := db.DB.Where("id = ?", classId).First(&class).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "Class not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Error finding class",
		})
	}

	// Updating class fields
	if name, ok := data["name"]; ok {
		class.Name = name
	}
	if schoolIdStr, ok := data["school_id"]; ok {
		schoolId, err := strconv.ParseUint(schoolIdStr, 10, 64)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Invalid school_id",
			})
		}

		// Verify if the schoolId exists
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

		class.SchoolID = uint(schoolId)
	}

	class.UpdatedAt = time.Now()

	// Save changes
	if err := db.DB.Save(&class).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Error updating class",
		})
	}

	// Prepare the response map with only necessary class information
	response := map[string]interface{}{
		"ID":         class.ID,
		"name":       class.Name,
		"school_id":  class.SchoolID,
		"created_at": class.CreatedAt,
		"updated_at": class.UpdatedAt,
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Class updated successfully",
		"data":    response,
	})
}

func DeleteClass(c *fiber.Ctx) error {
	classId := c.Params("classId")
	var class Models.Class
	db.DB.Where("id=?", classId).Find(&class)

	if class.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Class not found",
		})
	}
	db.DB.Unscoped().Where("id=?", classId).Delete(&class)
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Class deleted successfully",
	})
}

func GetClassDetails(c *fiber.Ctx) error {
	classId := c.Params("classId")
	var class Models.Class
	db.DB.Select("id,name,school_id,created_at,updated_at").Where("id =?", classId).First(&class)
	classData := make(map[string]interface{})
	classData["classId"] = class.ID
	classData["name"] = class.Name
	classData["school_id"] = class.SchoolID
	classData["createdAt"] = class.CreatedAt
	classData["updatedAt"] = class.UpdatedAt

	//what if there is no class or class id not provided
	if class.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Error fetching class",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "success",
		"data":    classData,
	})
}

func ClassesList(c *fiber.Ctx) error {

	var classes []Models.Class

	db.DB.Find(&classes)

	var classesData []map[string]interface{}
	for _, class := range classes {
		classData := map[string]interface{}{
			"id":   class.ID,
			"name": class.Name,
		}
		classesData = append(classesData, classData)
	}

	var count int64
	db.DB.Model(&Models.Class{}).Count(&count)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Class list API",
		"data":    classesData,
		"total":   count,
	})

}

func GetStudents(c *fiber.Ctx) error {
	classId := c.Params("classId")

	var class Models.Class
	// Fetch the class details
	db.DB.Select("id").Where("id = ?", classId).First(&class)

	// Check if class ID is valid
	if class.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Class not found",
		})
	}

	// Fetch the students for the class
	var students []Models.Student
	db.DB.Select("id", "name", "gender").Where("class_id = ?", classId).Find(&students)

	// Prepare the response data including student names and IDs
	var responseData []map[string]interface{}
	for _, student := range students {
		studentData := map[string]interface{}{
			"id":   student.ID,
			"name": student.Name,
		}
		responseData = append(responseData, studentData)
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    responseData,
	})
}
