package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/makersacademy/go-react-acebook-template/api/src/auth"
	"github.com/makersacademy/go-react-acebook-template/api/src/models"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func uploadFileToHostingService(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	client := resty.New()
	api_key := os.Getenv("IMGBB_API_KEY")
	client.SetFormData(map[string]string{
		"key": api_key,
	})

	// Open the file using the concrete type
	// src, err := fileHeader.Open()
	// if err != nil {
	// 	return "", fmt.Errorf("failed to open file: %v", err)
	// }
	// defer src.Close()

	resp, err := client.R().
		SetFileReader("image", fileHeader.Filename, file).
		Post("https://api.imgbb.com/1/upload")
	if err != nil {
		return "", err
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("failed to upload image: %s", resp.String())
	}

	var imgResponse struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	}
	err = json.Unmarshal(resp.Body(), &imgResponse)
	if err != nil {
		return "", err
	}

	return imgResponse.Data.URL, nil
}

func CreateUser(ctx *gin.Context) {
	var newUser models.User // Creates a variable called newUser with the User struct type User{gorm.Model(id,...), email, password}

	newUser = models.User{
		// Update user fields with file information
		Email:    ctx.PostForm("email"),
		Password: ctx.PostForm("password"),
		Username: ctx.PostForm("username"),
		PhotoURL: ctx.PostForm("image"),
	}

	if newUser.Email == "" || newUser.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "You must supply username and password"}) // Returns error if email and password are blank
		return
	}

	if len(newUser.Password) < 8 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Your password must be at least 8 characters"})
		return
	}

	var specialCharacters = []string{
		"!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "-", "_", "+", "=", "{", "}", "[", "]", "|", "\\", ":", ";", "'", "\"", "<", ">", ",", ".", "?", "/",
	}

	var containsSpecialCharacter = false
	for _, char := range newUser.Password {
		for _, specialChar := range specialCharacters {
			if string(char) == specialChar {
				containsSpecialCharacter = true
			}
		}
	}

	if containsSpecialCharacter != true {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Your password must have at least one special character"})
		return
	}

	if newUser.Email[0] == '@' {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email"})
		return
	}

	if !strings.Contains(newUser.Email, "@") {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email"})
		return
	}

	if strings.Count(newUser.Email, "@") > 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email"})
		return
	}

	if strings.Contains(newUser.Email, " ") {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email"})
		return
	}

	existingUser, err := models.FindUserByEmail(newUser.Email)
	if err != nil {
		SendInternalError(ctx, err)
		return
	}

	if existingUser != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "An account already exists with this email, try to login instead"})
		return
	}

	file, fileHeader, err := ctx.Request.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing image"})
		return
	}
	defer file.Close()

	// Upload the file to Imgbb
	photoURL, err := uploadFileToHostingService(file, fileHeader)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload photo"})
		return
	}

	newUser.PhotoURL = photoURL

	_, err = newUser.Save() // Adds newUser to database

	if err != nil {
		SendInternalError(ctx, err)
		return
	}

	userID := string(newUser.ID)
	token, _ := auth.GenerateToken(userID)

	userIDToken, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"ERROR": "USER ID NOT FOUND IN CONTEXT"})
		return
	}

	userIDString := userIDToken.(string)

	loggedUserID := strconv.Itoa(int([]byte(userIDString)[0]))

	ctx.JSON(http.StatusCreated, gin.H{"message": "OK", "token": token, "loggedUserID": loggedUserID}) //sends confirmation message back if successfully saved
}

func GetUser(ctx *gin.Context) {
	// userID := ctx.Param("id") // This is to check response in postman

	// The below two lines of code are to extract userID from token when that functionality becomes possible
	// val, _ := ctx.Get("userID")
	// userID := val.(string)
	userIDToken, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"ERROR": "USER ID NOT FOUND IN CONTEXT"})
		return
	}

	userIDString := userIDToken.(string)

	loggedUserID := strconv.Itoa(int([]byte(userIDString)[0]))

	token, _ := auth.GenerateToken(loggedUserID)
	user, err := models.FindUser(loggedUserID)
	if err != nil {
		SendInternalError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user, "token": token, "loggedUserID": loggedUserID})
}
