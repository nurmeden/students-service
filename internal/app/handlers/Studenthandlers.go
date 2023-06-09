package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	// _ "github.com/nurmeden/students-service/cmd/docs"
	"github.com/nurmeden/students-service/internal/app/model"
	"github.com/nurmeden/students-service/internal/app/usecase"
	"github.com/sirupsen/logrus"
)

const endpoint = "http://localhost:8080/api/courses/"

type StudentHandler struct {
	studentUsecase usecase.StudentUsecase
	logger         *logrus.Logger
}

func NewStudentHandler(studentUsecase usecase.StudentUsecase, logger *logrus.Logger) *StudentHandler {
	return &StudentHandler{
		studentUsecase: studentUsecase,
		logger:         logger,
	}
}

// CreateStudent godoc
// @Summary Create a new student
// @Description Create a new student with the input payload
// @Tags Students
// @Accept  json
// @Produce  json
// @Param student body model.Student true "Student data"
// @Success 201 {object} model.Student
// @Router /students/ [post]
// @Router /auth/sign-up/ [post]
func (h *StudentHandler) CreateStudent(c *gin.Context) {
	var student *model.Student
	err := c.ShouldBindJSON(&student)
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		h.logger.WithError(err).Error("Failed to decode request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode request body"})
		return
	}

	if student.Password == "" || student.Email == "" {
		h.logger.Error("Password or Email missing")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name or Email missing"})
		return
	}

	fmt.Printf("student.Email: %v\n", student.Email)

	exists, err := h.studentUsecase.CheckEmailExistence(context.Background(), student.Email)
	if err != nil {
		h.logger.WithError(err).Error("Failed to check email existence")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create student"})
		return
	}
	if exists {
		h.logger.Error("Email already exists")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}
	fmt.Printf("exists: %v\n", exists)
	createdStudent, err := h.studentUsecase.CreateStudent(context.Background(), student)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create student")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create student"})
		return
	}

	c.JSON(http.StatusCreated, createdStudent)
}

// GetStudentByID godoc
// @Summary Get student by ID
// @Description Get student by ID
// @Tags students
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} model.Student
// @Router /students/{id} [get]
func (h *StudentHandler) GetStudentByID(c *gin.Context) {
	fmt.Printf("c.Request.URL: %v\n", c.Request.URL)
	studentID := c.Param("id")
	fmt.Printf("studentID: %v\n", studentID)

	student, err := h.studentUsecase.GetStudentByID(context.Background(), studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("student get by id: %v\n", student)

	c.JSON(http.StatusOK, student)
}

// UpdateStudents godoc
// @Summary Update a student by ID
// @Description Update a student by ID
// @Tags students
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param student body model.Student true "Student object"
// @Success 200 {object} model.Student

// @Router /students/{id} [put]
func (h *StudentHandler) UpdateStudents(c *gin.Context) {
	studentID := c.Param("id")
	fmt.Printf("studentID: %v\n", studentID)
	var studentUpdateInput model.Student
	if err := c.ShouldBindJSON(&studentUpdateInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("studentUpdateInput: %v\n", studentUpdateInput)

	student, err := h.studentUsecase.UpdateStudent(context.Background(), studentID, &studentUpdateInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("studentUpdate: %v\n", student)
	c.JSON(http.StatusOK, gin.H{"data": student})
}

// DeleteStudent godoc
// @Summary Delete a student by ID
// @Description Delete a student by its ID
// @Tags students
// @Param id path int true "Student ID"
// @Success 200 {string} string "message: Студент успешно удален"

// @Router /students/{id} [delete]
func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	studentID := c.Param("id")

	err := h.studentUsecase.DeleteStudent(context.Background(), studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Студент успешно удален"})
}

func (h *StudentHandler) GetStudentsByCourseID(c *gin.Context) {
	courseID := c.Param("id")

	students, err := h.studentUsecase.GetStudentsByCourseID(context.Background(), courseID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get students by course ID")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get students by course ID"})
		return
	}

	if students == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No students found for course ID"})
		return
	}

	c.JSON(http.StatusOK, students)
}

// SignIn godoc
// @Summary Sign in a student
// @Description Authenticates a student using their email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param signInData body model.SignInData true "Sign in data"
// @Success 200 {object} TokenResponse

// @Router /api/sign-in [post]
func (h *StudentHandler) SignIn(c *gin.Context) {
	var signInData model.SignInData
	err := c.ShouldBindJSON(&signInData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode request body"})
		return
	}

	authResult, err := h.studentUsecase.SignIn(context.Background(), &signInData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate"})
		return
	}

	// refreshToken := uuid.New().String()
	// err = h.studentUsecase.SaveRefreshToken(authResult.UserID, refreshToken)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save refresh token"})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{"token": authResult.Token})
}

func (sc *StudentHandler) GetStudentCourses(c *gin.Context) {
	studentID := c.Param("id")
	fmt.Printf("studentID: %v\n", studentID)
	resp, err := http.Get(endpoint + studentID + "/courses")
	fmt.Printf("resp: %v\n", resp)
	if err != nil {
		// Обработка ошибки
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student courses"})
		return
	}
	defer resp.Body.Close()

	fmt.Printf("resp: %v\n", resp)

	var course *model.CourseResponse
	err = json.NewDecoder(resp.Body).Decode(&course)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode student courses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"courses": course})
}

func (h *StudentHandler) Logout(c *gin.Context) {
	userID := c.GetString("user_id")                   // получаем ID пользователя из контекста Gin
	err := h.studentUsecase.DeleteRefreshToken(userID) // удаляем refresh token из базы данных
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete refresh token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func (h *StudentHandler) RefreshToken(c *gin.Context) {
	// Получаем refresh token из запроса
	refreshToken := c.PostForm("refresh_token")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token is missing"})
		return
	}

	// Проверяем, что refresh token является действительным
	userID, err := h.studentUsecase.ValidateRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh_token"})
		return
	}

	// Генерируем новый access token
	token, err := h.studentUsecase.GenerateToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// Отправляем новый access token в ответе
	c.JSON(http.StatusOK, gin.H{"token": token})
}
