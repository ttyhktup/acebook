package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/makersacademy/go-react-acebook-template/api/src/auth"
	"github.com/makersacademy/go-react-acebook-template/api/src/models"
)

type JSONComment struct {
	ID      uint   `json:"_id"`
	Message string `json:"message"`
	Likes   int    `json:"likes"`
	PostId  int    `json:"post_id"`
	User    JSONUser
}

type createCommentRequestBody struct {
	Message string
}

func CreateComment(ctx *gin.Context) {
	var requestBody createCommentRequestBody
	err := ctx.BindJSON(&requestBody)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	if len(requestBody.Message) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Post message empty"})
		return
	}

	postID := ctx.Param("id")
	postIdInt, err := strconv.Atoi(postID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	CommentTime := time.Now()

	userIDToken, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"ERROR": "USER ID NOT FOUND IN CONTEXT"})
		return
	}

	userIDString := userIDToken.(string)

	newComment := models.Comment{
		UserID:    strconv.Itoa(int([]byte(userIDString)[0])),
		Message:   requestBody.Message,
		CreatedAt: CommentTime,
		Likes:     0,
		PostId:    postIdInt,
	}

	_, err = newComment.Save()
	if err != nil {
		SendInternalError(ctx, err)
		return
	}

	comments, err := models.FetchAllCommentsByPostId(postIdInt)

	if err != nil {
		SendInternalError(ctx, err)
		return
	}

	val, _ := ctx.Get("userID")
	userID := val.(string)
	token, _ := auth.GenerateToken(userID)

	loggedUserID := strconv.Itoa(int([]byte(userID)[0]))

	// Convert comments to JSON Structs
	jsonComments := make([]JSONComment, 0)
	for _, comment := range *comments {
		user, err := models.FindUser(comment.UserID)
		if err != nil {
			fmt.Println("FindUser error in CreateComment: ", err)
			user.ID = 0
			user.Username = ""
		}
		jsonComments = append(jsonComments, JSONComment{
			Message: comment.Message,
			ID:      comment.ID,
			Likes:   comment.Likes,
			PostId:  comment.PostId,
			User: JSONUser{
				UserID:   user.ID,
				Username: user.Username,
				PhotoURL: user.PhotoURL,
			},
		})
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Comment created", "comments": jsonComments, "token": token, "loggedUserID": loggedUserID})
}

func GetAllCommentsByPostId(ctx *gin.Context) {
	// retrieving a parameter from the request URL here below
	// needing route to be structured like this: /posts/:post_id
	postID := ctx.Param("id")

	postIdInt, err := strconv.Atoi(postID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	comments, err := models.FetchAllCommentsByPostId(postIdInt)

	if err != nil {
		SendInternalError(ctx, err)
		return
	}

	val, _ := ctx.Get("userID")
	userID := val.(string)
	token, _ := auth.GenerateToken(userID)

	loggedUserID := strconv.Itoa(int([]byte(userID)[0]))

	// Convert comments to JSON Structs
	jsonComments := make([]JSONComment, 0)
	for _, comment := range *comments {
		user, err := models.FindUser(comment.UserID)
		if err != nil {
			fmt.Println("FindUser error in CreateComment: ", err)
			user.ID = 0
			user.Username = ""
		}
		jsonComments = append(jsonComments, JSONComment{
			Message: comment.Message,
			ID:      comment.ID,
			Likes:   comment.Likes,
			PostId:  comment.PostId,
			User: JSONUser{
				UserID:   user.ID,
				Username: user.Username,
				PhotoURL: user.PhotoURL,
			},
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"comments": jsonComments, "token": token, "loggedUserID": loggedUserID})
}

func GetSpecificComment(ctx *gin.Context) {
	postID := ctx.Param("id")

	postIdInt, err := strconv.Atoi(postID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	commentID := ctx.Param("comment_id")

	commentIdInt, err := strconv.Atoi(commentID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	comment, err := models.FetchSpecificComment(postIdInt, commentIdInt)
	if err != nil {
		SendInternalError(ctx, err)
		return
	}

	user, err := models.FindUser(comment.UserID)
	if err != nil {
		fmt.Println("FindUser error in CreateComment: ", err)
		user.ID = 0
		user.Username = ""
	}

	val, _ := ctx.Get("userID")
	userID := val.(string)
	token, _ := auth.GenerateToken(userID)

	loggedUserID := strconv.Itoa(int([]byte(userID)[0]))

	jsonComment := JSONComment{
		Message: comment.Message,
		ID:      comment.ID,
		Likes:   comment.Likes,
		PostId:  comment.PostId,
		User: JSONUser{
			UserID:   user.ID,
			Username: user.Username,
			PhotoURL: user.PhotoURL,
		},
	}

	ctx.JSON(http.StatusOK, gin.H{"comment": jsonComment, "token": token, "loggedUserID": loggedUserID})
}

func DeleteComment(ctx *gin.Context) {
	postID := ctx.Param("id")

	postIdInt, err := strconv.Atoi(postID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	commentID := ctx.Param("comment_id")

	commentIdInt, err := strconv.Atoi(commentID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	comment, err := models.FetchSpecificComment(postIdInt, commentIdInt)
	if err != nil {
		SendInternalError(ctx, err)
		return
	}

	userIDToken, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"ERROR": "USER ID NOT FOUND IN CONTEXT"})
		return
	}

	userIDString := userIDToken.(string)

	token, _ := auth.GenerateToken(userIDString)

	loggedUserID := strconv.Itoa(int([]byte(userIDString)[0]))

	if comment.UserID != strconv.Itoa(int([]byte(userIDString)[0])) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User ID can only delete own post"})
		return
	}

	// Check if the comment is nil
	if comment == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Delete post from database
	DeletedComment, err := comment.Delete()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Comment deleted successfully", "deleted comment": DeletedComment, "token": token, "loggedUserID": loggedUserID})
}

func UpdateCommentLikes(ctx *gin.Context) {
	// Get the post ID from the URL path parameter
	postID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Get the comment ID from the URL path parameter
	commentID, err := strconv.ParseUint(ctx.Param("comment_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	// Fetch the comment from the database
	comment, err := models.FetchSpecificComment(int(postID), int(commentID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Check if the post is nil
	if comment == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Increment the likes count
	likedComment, err := comment.SaveLike()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save like in comment"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Like added successfully", "liked_comment": likedComment})
}
