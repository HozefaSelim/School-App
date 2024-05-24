package Controllers

import (
	db "SchoolProject/Config"
	"SchoolProject/Models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func CreateTeacher(c *fiber.Ctx) error {
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
				"message": "Teacher Name is required",
			})
	}
	if data["phone"] == "" {
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "Teacher Phone is required",
			})
	}
	if data["email"] == "" {
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "Teacher email is required",
			})
	}

	if data["salary"] == "" {
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "salary is required",
			})
	}
	salary, err := strconv.ParseUint(data["salary"], 10, 64)
	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "salary must be int",
			})
	}

	// school id adding
	schoolId, err := strconv.ParseUint(data["school_id"], 10, 64)
	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "School_id must be Uint ",
			})
	}
	var school Models.Subject
	if err := db.DB.Where("id = ?", schoolId).First(&school).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Subject not found",
		})
	}

	// subject id adding
	subjectId, err := strconv.ParseUint(data["subject_id"], 10, 64)
	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "Subject_id must be Uint ",
			})
	}
	var subject Models.Subject
	if err := db.DB.Where("id = ?", subjectId).First(&subject).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Subject not found",
		})
	}

	// now saving teacher to db
	teacher := Models.Teacher{
		Name:      data["name"],
		Phone:     data["phone"],
		Email:     data["email"],
		Salary:    int(salary),
		SchoolID:  uint(schoolId),
		SubjectID: uint(subjectId),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db.DB.Create(&teacher)
	return c.Status(200).JSON(
		fiber.Map{
			"success": true,
			"message": "Success",
			"data":    teacher,
		})
}

func UpdateTeacher(c *fiber.Ctx) error {

	teacherId := c.Params("teacherId")
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid data",
		})
	}

	var teacher Models.Teacher
	if err := db.DB.Where("id = ?", teacherId).First(&teacher).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "Teacher not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Error finding teacher",
		})
	}

	// Updating teacher fields
	if name, ok := data["name"]; ok {
		teacher.Name = name
	}
	if email, ok := data["email"]; ok {
		teacher.Email = email
	}
	salary, err := strconv.ParseUint(data["salary"], 10, 64)
	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "salary must be Uint ",
			})
		teacher.Salary = int(salary)
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

		teacher.SchoolID = uint(schoolId)
	}
	if subjectIdStr, ok := data["subject_id"]; ok {
		subjectId, err := strconv.ParseUint(subjectIdStr, 10, 64)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Invalid subject_id",
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

		teacher.SubjectID = uint(subjectId)
	}
	teacher.UpdatedAt = time.Now()

	// Save changes
	db.DB.Save(&teacher)
	// Prepare the response map with only necessary class information
	response := map[string]interface{}{
		"ID":         teacher.ID,
		"name":       teacher.Name,
		"email":      teacher.Email,
		"salary":     teacher.Salary,
		"school_id":  teacher.SchoolID,
		"subject_id": teacher.SubjectID,
		"created_at": teacher.CreatedAt,
		"updated_at": teacher.UpdatedAt,
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Teacher updated successfully",
		"data":    response,
	})
}

func DeleteTeacher(c *fiber.Ctx) error {
	teacherId := c.Params("teacherId")
	var teacher Models.Teacher
	db.DB.Where("id=?", teacherId).Find(&teacher)

	if teacher.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Teacher not found",
		})
	}
	db.DB.Unscoped().Where("id=?", teacherId).Delete(&teacher)
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Teacher deleted successfully",
	})
}

func GetTeacherDetails(c *fiber.Ctx) error {

	teacherId := c.Params("teacherId")
	var teacher Models.Teacher
	db.DB.Select("id,name,phone,salary,school_id,subject_id,created_at,updated_at").Where("id =?", teacherId).First(&teacher)
	teacherData := make(map[string]interface{})
	teacherData["teacherId"] = teacher.ID
	teacherData["name"] = teacher.Name
	teacherData["salary"] = teacher.Salary
	teacherData["phone"] = teacher.Phone
	teacherData["school_id,"] = teacher.SchoolID
	teacherData["subject_id"] = teacher.SubjectID
	teacherData["createdAt"] = teacher.CreatedAt
	teacherData["updatedAt"] = teacher.UpdatedAt

	//what if there is no teacher or teacher id not provided
	if teacher.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Error fetching teacher",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "success",
		"data":    teacherData,
	})
}

func TeachersList(c *fiber.Ctx) error {
	var teachers []Models.Teacher

	db.DB.Find(&teachers)

	var teachersData []map[string]interface{}
	for _, teacher := range teachers {
		teacherData := map[string]interface{}{
			"id":   teacher.ID,
			"name": teacher.Name,
		}
		teachersData = append(teachersData, teacherData)
	}

	var count int64
	db.DB.Model(&Models.Teacher{}).Count(&count)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Teacher list API",
		"data":    teachersData,
		"total":   count,
	})
}

func GetTeacherSubject(c *fiber.Ctx) error {
	teacherId := c.Params("teacherId")

	// Find the teacher by ID to check if they exist
	var teacher Models.Teacher
	if err := db.DB.Preload("Subject").First(&teacher, teacherId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "Teacher not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Error retrieving teacher",
		})
	}

	// Check if the teacher has an associated subject
	if teacher.SubjectID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "No subject found for the teacher",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Subject retrieved successfully",
		"data":    teacher.Subject,
	})
}
