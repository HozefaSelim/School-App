package Routes

import (
	"SchoolProject/Controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	// School routes
	app.Post("/Schools", Controllers.CreateSchool)
	app.Get("/Schools", Controllers.SchoolsList)
	app.Get("/Schools/:SchoolId", Controllers.GetSchoolDetails)
	app.Delete("/Schools/:SchoolId", Controllers.DeleteSchool)
	app.Put("/Schools/:SchoolId", Controllers.UpdateSchool)
	app.Get("/SchoolClasses/:SchoolId", Controllers.GetClasses)
	app.Get("/SchoolTeachers/:SchoolId", Controllers.GetTeachersInSchool)
	app.Get("/SchoolStudents/:SchoolId", Controllers.GetStudentsInSchool)

	// Class routes
	app.Post("/Classes", Controllers.CreateClass)
	app.Get("/Classes", Controllers.ClassesList)
	app.Get("/Classes/:classId", Controllers.GetClassDetails)
	app.Delete("/Classes/:classId", Controllers.DeleteClass)
	app.Put("/Classes/:classId", Controllers.UpdateClass)
	app.Get("/ClassStudents/:classId", Controllers.GetStudents)

	// Student routes
	app.Post("/Students", Controllers.CreateStudent)
	app.Get("/Students", Controllers.StudentsList)
	app.Get("/Students/:StudentId", Controllers.GetStudentDetails)
	app.Delete("/Students/:StudentId", Controllers.DeleteStudent)
	app.Put("/Students/:StudentId", Controllers.UpdateStudent)
	app.Put("/StudentSubject/:StudentId", Controllers.AddStudentToSubject)
	app.Get("/StudentSubjects/:StudentId", Controllers.GetStudentSubjects)

	// Teacher routes
	app.Post("/Teachers", Controllers.CreateTeacher)
	app.Get("/Teachers", Controllers.TeachersList)
	app.Get("/Teachers/:TeacherId", Controllers.GetTeacherDetails)
	app.Delete("/Teachers/:TeacherId", Controllers.DeleteTeacher)
	app.Put("/Teachers/:TeacherId", Controllers.UpdateTeacher)
	app.Get("/TeacherSubject/:TeacherId", Controllers.GetTeacherSubject)

	// Subject routes
	app.Post("/Subjects", Controllers.CreateSubject)
	app.Get("/Subjects", Controllers.SubjectsList)
	app.Get("/Subjects/:SubjectId", Controllers.GetSubjectDetails)
	app.Delete("/Subjects/:SubjectId", Controllers.DeleteSubject)
	app.Put("/Subjects/:SubjectId", Controllers.UpdateSubject)
	app.Get("/SubjectTeachers/:SubjectId", Controllers.GetTeachers)
	app.Get("/StudentsAndTeachersSubject/:SubjectId", Controllers.GetTeachersAndStudentsBySubject)

}
