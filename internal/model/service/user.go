package service

import (
	"log"
	"net/http"

	"github.com/Lutwidse/Techtrain-API/internal/model/data"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/gorm"
)

type UserService struct {
	Db   *gorm.DB
	User data.User
}

func (s *UserService) Create(c *gin.Context) {
	if err := c.BindJSON(&s.User); err != nil {
		log.Fatal(err)
	}
	token := uuid.New().String()
	userReq := data.User{Name: s.User.Name, XToken: token}

	result := s.Db.Exec("INSERT INTO `techtrain_db`.`users` (`name`, `x_token`) VALUES (?, ?)", userReq.Name, userReq.XToken)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Already Registered"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (s *UserService) Get(c *gin.Context) {
	token := c.GetHeader("x-token")
	if token == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Token Required"})
		return
	}

	user := s.User
	result := s.Db.First(&user, "x_token = ?", token)
	if result.Error != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "No User Found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"name": user.Name})
}

func (s *UserService) Update(c *gin.Context) {
	if err := c.BindJSON(&s.User); err != nil {
		log.Fatal(err)
	}
	token := c.GetHeader("x-token")
	if token == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Token Required"})
		return
	}

	userReq := data.User{Name: s.User.Name, XToken: token}

	user := s.User
	result := s.Db.First(&user, "x_token = ?", token)
	if result.Error != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "No User Found"})
		return
	}

	newName := userReq.Name
	oldName := user.Name

	s.Db.Exec("UPDATE `techtrain_db`.`users` SET `name` = ? WHERE (`name` = ?) and (`x_token` = ?)", newName, oldName, token)
	c.JSON(http.StatusOK, nil)
}
