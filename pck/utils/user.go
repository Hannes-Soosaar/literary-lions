package utils

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	// "gitea.kood.tech/hannessoosaar/literary-lions/intenal/config"
	"gitea.kood.tech/hannessoosaar/literary-lions/intenal/config"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/models"
)

// ? perhaps we are getting too much information
func FindUserByUserName(userName string) models.User {
	var user models.User
	// See if we can make the open db into a separate function so we do not need to open close the DB for every request
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		fmt.Println("error opening DB", err)
		return user // should return an empty use
	}
	defer db.Close()
	query := "SELECT id, username,email,password,role,created_at,modified_at,active,uuid FROM users  WHERE username = ?"
	row := db.QueryRow(query, userName)
	err = row.Scan(&user.Username, &user.Email, &user.Role, &user.Password, &user.Role, &user.CreatedAt, &user.ModifiedAt, &user.Active, &user.UUID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("%s, was not found \n", userName)
			return models.User{} // Consider returning an error too
		}
		fmt.Printf("error scanning rows: %v", err)
		return models.User{} // consider returning an error too
	}
	fmt.Printf("User %v \n:", user)
	return user
}

func AddActiveUser(user models.User) error {
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		fmt.Println("error opening DB", err)
		return err
	}
	defer db.Close()
	dbUser := FindUserByUserName(user.Username)
	fmt.Println("finding user!")
	fmt.Println(user)
	if (dbUser == models.User{}) {
		query := "INSERT INTO users (username,email,password,role,created_at,modified_at,active,uuid) VALUES (?,?,?,?,?,?,?,?)"
		_, err := db.Exec(query, user.Username, user.Email, user.Password, user.Role, user.CreatedAt, user.ModifiedAt, user.Active, user.UUID)
		if err != nil {
			fmt.Printf("Error adding User: %v, Error: %v", user, err)
		}
	} else {
		fmt.Printf("There is User %s by this name!", dbUser.Username)
		return nil
	}
	//? Should we check if there is the user is also active and inactive
	fmt.Printf("User %s added", user.Username)
	return nil
}

func DeleteActiveUser(userID int) error {

	// Find User by ID and changes the user to not active
	return nil
}

func FindUserByID(userId int ) error {

return nil
}


func AddUserTest() {

	var user models.User

	user.Username = "NewName"
	user.Email = "new@Email.com"
	user.Password = "123"
	user.Role = "U"                                      // ? should we have this with integers as well ?
	user.CreatedAt = time.Now().Format("02/01/06,15/04") //
	user.ModifiedAt = time.Now().Format("02/01/06,15/04")
	user.Active = config.ACTIVE
	userUuid, err := GenerateUUID()
	if err != nil {
		fmt.Printf("there was an error generating a UUID")
		log.Panic(err) // ! a bit over board with this!
	}
	user.UUID = userUuid
	err = AddActiveUser(user)
	if err != nil {
		fmt.Printf("The user %v was not added %v", user, err )
	}
}

func ValidateUser(userName string) error {

	user := FindUserByUserName(userName)
	if user.Active == 1 {
		return nil
	}

	return nil
}
