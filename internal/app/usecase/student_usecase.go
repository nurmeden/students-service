package usecase

import (
	"context"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/nurmeden/students-service/internal/app/model"
	"github.com/nurmeden/students-service/internal/app/repository"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type StudentUsecase interface {
	CreateStudent(ctx context.Context, student *model.Student) (*model.Student, error)
	GetStudentByID(ctx context.Context, id string) (*model.Student, error)
	GetStudentByCoursesID(ctx context.Context, id string) (*model.Student, error)
	UpdateStudent(ctx context.Context, student_id string, student *model.Student) (*model.Student, error)
	DeleteStudent(ctx context.Context, id string) error
	SignIn(signInData *model.SignInData) (*model.AuthToken, error)
}

type studentUsecase struct {
	studentRepo repository.StudentRepository
	logger      *logrus.Logger
}

func NewStudentUsecase(studentRepo repository.StudentRepository, logger *logrus.Logger) StudentUsecase {
	return &studentUsecase{
		studentRepo: studentRepo,
		logger:      logger,
	}
}

func (u *studentUsecase) CreateStudent(ctx context.Context, student *model.Student) (*model.Student, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(student.Password), bcrypt.DefaultCost)
	if err != nil {
		u.logger.Errorf("Error generating password hash: %v", err)
		return nil, err
	}
	student.Password = string(hashedPassword)

	return u.studentRepo.Create(ctx, student)
}

func (u *studentUsecase) GetStudentByID(ctx context.Context, id string) (*model.Student, error) {
	return u.studentRepo.Read(ctx, id)
}

func (u *studentUsecase) GetStudentByCoursesID(ctx context.Context, id string) (*model.Student, error) {
	return u.studentRepo.GetStudentByCoursesID(ctx, id)
}

func (u *studentUsecase) UpdateStudent(ctx context.Context, student_id string, student *model.Student) (*model.Student, error) {
	student, err := u.studentRepo.Read(ctx, student_id)
	if err != nil {
		return nil, err
	}

	return u.studentRepo.Update(ctx, student)
}

func (u *studentUsecase) DeleteStudent(ctx context.Context, id string) error {
	return u.studentRepo.Delete(ctx, id)
}

func (u *studentUsecase) SignIn(signInData *model.SignInData) (*model.AuthToken, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = signInData.UserID
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return nil, err
	}

	authToken := &model.AuthToken{
		UserID:    signInData.UserID,
		Token:     tokenString,
		ExpiresAt: time.Now().Add(time.Hour * 1),
	}
	return authToken, nil
}
