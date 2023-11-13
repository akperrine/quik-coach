package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/akperrine/quik-coach/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var collection *mongo.Collection

func UserCollection(c *mongo.Database) {
	collection = c.Collection("users")
}

func HandleRequests() {
	http.HandleFunc("/users", getAllUsers)
	http.HandleFunc("/users/register", registerUser)
	http.HandleFunc("/users/login", loginUser)
	http.HandleFunc("/health_check", healthCheck)
	log.Fatal(http.ListenAndServe(":8000", nil))
}


func getAllUsers(w http.ResponseWriter, r *http.Request) {
	users := []models.User{}

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Println(err)
		handleError(w, http.StatusInternalServerError, "Failed to retrieve users from MongoDB")
		return
	}
	defer cursor.Close(context.Background())
	
	
	for cursor.Next(context.Background()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			log.Printf("Failed to decode user with error: %s", err)
			return
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Error during cursor iteration with error: %s", err)
		return
	}

	usersJSON, err := json.Marshal(users)
	if err != nil {
		log.Println("Failed to marshal users to JSON")
		return
	}

	writeJSONResponse(w, http.StatusOK, usersJSON)
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	json.NewDecoder(r.Body).Decode(user)

	encrptedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}

	user.ID = uuid.NewString()
	user.Password = string(encrptedPass)

    fmt.Printf("{ID: %s, FirstName: %s, LastName: %s, Email: %s}", user.ID, user.FirstName, user.LastName, user.Email)


	createdUser, insErr := collection.InsertOne(context.TODO(), user)

	if insErr != nil {
		fmt.Printf("Error creating new user: %s", insErr)
	}

	json.NewEncoder(w).Encode(createdUser)
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hi")
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user) 
	fmt.Printf("e: %s, p: %s \n", user.Email, user.Password)
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Invalid request"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	response := findOne(user.Email, user.Password)

	json.NewEncoder(w).Encode(response)

	// if foundUser.Email == "" || foundUser.Password == "" {
	// 	writeJSONResponse(w, http.StatusUnauthorized, []byte("Incorrect email or password"))
	// } else {
	// 	userJson, err := marshalWithoutPassword(foundUser)
	// 	if err != nil {
	// 		fmt.Println("Error: ", err)
	// 	}
	// 	writeJSONResponse(w, http.StatusOK, userJson)
	// }
}

func findOne(email, password string) map[string]interface{}{
	fmt.Println(password)
	user := &models.User{}


	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(user)

	if err == mongo.ErrNoDocuments {
		// No matching user found
		fmt.Println("User not found")
		var resp = map[string]interface{}{"status": false, "message": "Email address not found"}
		return resp
	} else if err != nil {
		log.Fatal(err)
		var resp = map[string]interface{}{"status": false, "message": "Invalid login credentials. Please try again"}
		return resp
	} else {
		fmt.Printf("Found user: %+v\n", user)
	}

	passwordErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if passwordErr != nil {
        // Passwords don't match
        fmt.Println("Invalid password")
		var resp = map[string]interface{}{"status": false, "message": "Invalid login credentials. Please try again"}
		return resp
    }

	expiresAt := jwt.NewNumericDate(time.Now().Add(time.Minute * 100000))



	claims := &models.Token{
		UserID: user.ID,
		Name:	user.FirstName,
		Email:	user.Email,
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}
	var resp = map[string]interface{}{"status": false, "message": "logged in"}
	resp["token"] = tokenString //Store the token in the response
	resp["user"] = user
	return resp
}


func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "healthy!")
}

func marshalWithoutPassword(u models.User) ([]byte, error) {
	type user models.User
	x := user(u)
	x.Password = ""

	return json.Marshal(x)
}
