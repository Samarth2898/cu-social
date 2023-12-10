package main

import (
	"context"
	helper "cu-social/lib"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
)

// Book ...
type userResponse struct {
	UserId    int32     `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type FeedObject struct {
	Description string `json:"description" db:"description"`
	VideoURL    string `json:"video_url" db:"video_url"`
	Title       string `json:"title" db:"title"`
	UserID      int    `json:"user_id" db:"user_id"`
	PostedBy    string `json:"username" db:"username"`
}

type UserSignupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSignupResponse struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

type UserInfoUpdate struct {
	ProfilePicture string `json:"profile_picture"`
	Biography      string `json:"biography"`
	UserID         int32  `json:"user_id"`
}

type UserInfoUpdateResponse struct {
	Success bool `json:"success"`
}

const (
	backendServer = "http://localhost:8080"
)

func main() {

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "sa.json")

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/styles", "./styles")
	r.Static("/js", "./js")
	r.Static("/images", "./images")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", gin.H{})
	})

	r.POST("/login-user", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		fmt.Println("username: ", username, "password: ", password)
		loginReqBody := UserLoginRequest{
			Username: username,
			Password: password,
		}

		marshalledLoginBody, _ := json.Marshal(loginReqBody)
		res, err := helper.PostReq("http://0.0.0.0:8080/users/login", marshalledLoginBody)
		if err != nil {
			fmt.Println("error sending POST request: ", err.Error())
		}
		defer res.Body.Close()

		postLoginResponse := &UserLoginResponse{}
		derr := json.NewDecoder(res.Body).Decode(postLoginResponse)
		if derr != nil {
			panic(derr)
		}
		fmt.Println("logged in: ", postLoginResponse.User.UserId, postLoginResponse.AccessToken)
		if res.Status == "200 OK" {
			fmt.Println("setting creds")
			c.Redirect(http.StatusSeeOther, "/feed?access_token="+postLoginResponse.AccessToken+"&user_id="+strconv.Itoa(int(postLoginResponse.User.UserId)))
		} else {
			c.Set("isAuthenticated", false)
			c.HTML(http.StatusBadRequest, "login.html", gin.H{})
		}
	})

	r.GET("/login-page", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})

	r.POST("/submit-signup-form", func(c *gin.Context) {
		username := c.PostForm("username")
		email := c.PostForm("email")
		password := c.PostForm("password")

		reqBody := UserSignupRequest{
			Username: username,
			Email:    email,
			Password: password,
		}
		marshalledBody, _ := json.Marshal(reqBody)
		res, err := helper.PostReq("http://0.0.0.0:8080/users", marshalledBody)
		if err != nil {
			fmt.Println("error sending POST request: ", err.Error())
		}
		defer res.Body.Close()

		postResponse := &UserSignupResponse{}
		derr := json.NewDecoder(res.Body).Decode(postResponse)
		if derr != nil {
			fmt.Println("error decoding signup response: ", derr.Error())
		}

		fmt.Println("signed user: ", postResponse.Username, res.StatusCode)
		loginReqBody := UserLoginRequest{
			Username: postResponse.Username,
			Password: password,
		}

		marshalledLoginBody, _ := json.Marshal(loginReqBody)
		res, err = helper.PostReq("http://0.0.0.0:8080/users/login", marshalledLoginBody)
		if err != nil {
			fmt.Println("error sending POST request: ", err.Error())
		}
		defer res.Body.Close()

		postLoginResponse := &UserLoginResponse{}
		derr = json.NewDecoder(res.Body).Decode(postLoginResponse)
		if derr != nil {
			panic(derr)
		}
		fmt.Println("logged in: ", postLoginResponse.User.UserId, postLoginResponse.AccessToken)
		if res.Status == "200 OK" {
			fmt.Println("setting creds")
			c.HTML(http.StatusOK, "create_user_profile.html", gin.H{
				"is_authenticated": true,
				"access_token":     postLoginResponse.AccessToken,
				"user_id":          postLoginResponse.User.UserId,
			})
		} else {
			c.HTML(http.StatusBadRequest, "create_user_profile.html", gin.H{
				"is_authenticated": false,
			})
		}
	})

	r.POST("/submit-profile", func(c *gin.Context) {
		userID, _ := strconv.Atoi(c.Query("user_id"))
		access_token := c.Query("access_token")
		url := uploadProfilePhoto(c)
		bio := c.PostForm("userBio")
		updateReqBody := UserInfoUpdate{
			UserID:         int32(userID),
			ProfilePicture: url,
			Biography:      bio,
		}
		marshalledBody, _ := json.Marshal(updateReqBody)
		res, err := helper.PostReq("http://0.0.0.0:8080/users/update", marshalledBody)
		if err != nil {
			fmt.Println("error sending POST request: ", err.Error())
		}
		defer res.Body.Close()

		userUpdateResponse := &UserInfoUpdateResponse{}
		derr := json.NewDecoder(res.Body).Decode(userUpdateResponse)
		if derr != nil {
			fmt.Println("error in decoding: ", derr.Error())
		}
		fmt.Println("Update status: ", userUpdateResponse.Success)
		c.Redirect(http.StatusSeeOther, "/feed?access_token="+access_token+"&user_id="+strconv.Itoa(userID))
	})

	r.GET("/feed", func(c *gin.Context) {
		userID, _ := strconv.Atoi(c.Query("user_id"))
		FeedObjectInstance := getFeedByUser(userID)
		c.HTML(http.StatusOK, "feed.html", gin.H{
			"FeedObjects": FeedObjectInstance,
		})
	})

	r.GET("/add_post", func(c *gin.Context) {
		c.HTML(http.StatusOK, "add_post.html", gin.H{})
	})

	r.GET("/profile", func(c *gin.Context) {
		userID := 4 // samy get userID from JWT token
		FeedObjectInstance := getProfileFeed(userID)
		c.HTML(http.StatusOK, "profile.html", gin.H{
			"PostObjects": FeedObjectInstance,
		})
	})

	r.POST("/uploadVideo", func(c *gin.Context) {
		uploadVideoFunc(c)
	})

	log.Fatal(r.Run(":3000"))
}

func getFeedByUser(userID int) *[]FeedObject {

	url := fmt.Sprintf("%s/users/%d/following/posts", backendServer, userID)

	// Send GET request to fetch user posts
	response, err := http.Get(url)
	if err != nil {
		log.Fatal("Error fetching posts:", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("Failed to fetch posts. Status code: %d", response.StatusCode)
	}

	// Read response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error reading response:", err)
	}

	// Unmarshal JSON response into Post slice
	var posts []FeedObject
	err = json.Unmarshal(body, &posts)
	if err != nil {
		log.Fatal("Error unmarshaling JSON:", err)
	}

	return &posts

}

func getProfileFeed(userID int) *[]FeedObject {

	url := fmt.Sprintf("%s/users/%d/posts", backendServer, userID)

	// Send GET request to fetch user posts
	response, err := http.Get(url)
	if err != nil {
		log.Fatal("Error fetching posts:", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("Failed to fetch posts. Status code: %d", response.StatusCode)
	}

	// Read response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error reading response:", err)
	}

	// Unmarshal JSON response into Post slice
	var posts []FeedObject
	err = json.Unmarshal(body, &posts)
	if err != nil {
		log.Fatal("Error unmarshaling JSON:", err)
	}

	return &posts

}

func uploadVideoFunc(c *gin.Context) {

	bucketName := "cusocial"

	ctx := context.Background()

	file, handler, err := c.Request.FormFile("videoFile")
	if err != nil {
		c.String(http.StatusBadRequest, "Error retrieving file")
		return
	}
	defer file.Close()

	client, err := storage.NewClient(ctx)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to create GCS client")
		return
	}
	defer client.Close()

	objectName := handler.Filename
	obj := client.Bucket(bucketName).Object(objectName)
	wObj := obj.NewWriter(ctx)
	defer wObj.Close()

	if _, err := io.Copy(wObj, file); err != nil {
		c.String(http.StatusInternalServerError, "Error uploading file to GCS")
		return
	}

	if err := wObj.Close(); err != nil {
		c.String(http.StatusInternalServerError, "Error closing GCS writer")
		return
	}

	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName)

	c.String(http.StatusOK, "File uploaded successfully! URL: %s", url)
	c.Redirect(http.StatusMovedPermanently, "/feed")
}

func uploadProfilePhoto(c *gin.Context) string {

	bucketName := "cusocial"

	ctx := context.Background()

	file, handler, err := c.Request.FormFile("profileImage")
	if err != nil {
		c.String(http.StatusBadRequest, "Error retrieving file")
		return ""
	}
	defer file.Close()

	client, err := storage.NewClient(ctx)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to create GCS client")
		return ""
	}
	defer client.Close()

	objectName := handler.Filename
	obj := client.Bucket(bucketName).Object(objectName)
	wObj := obj.NewWriter(ctx)
	defer wObj.Close()

	if _, err := io.Copy(wObj, file); err != nil {
		c.String(http.StatusInternalServerError, "Error uploading file to GCS")
		return ""
	}

	if err := wObj.Close(); err != nil {
		c.String(http.StatusInternalServerError, "Error closing GCS writer")
		return ""
	}

	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName)

	return url
}
