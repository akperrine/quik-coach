package models

import "fmt"


type User struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"`
 }

 func (u User) userToString() string {
	return fmt.Sprintf("ID: %s, FirstName: %s, LastName: %s, Email: %s", u.ID, u.FirstName, u.LastName, u.Email)
}

