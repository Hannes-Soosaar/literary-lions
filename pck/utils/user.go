package utils

import (
	"database/sql"
	"fmt"
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
	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.ModifiedAt, &user.Active, &user.UUID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("%s, was not found \n", userName)
			return models.User{} // Consider returning an error too
		}
		fmt.Printf("error scanning rows: %v", err)
		return models.User{} // consider returning an error too
	}
	return user
}
func FindUserByUserID(ID int) models.User {
	var user models.User
	// See if we can make the open db into a separate function so we do not need to open close the DB for every request
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		fmt.Println("error opening DB", err)
		return user // should return an empty use
	}
	defer db.Close()
	query := "SELECT id, username,email,password,role,created_at,modified_at,active,uuid FROM users  WHERE ID = ?"
	row := db.QueryRow(query, ID)
	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.ModifiedAt, &user.Active, &user.UUID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("user with %d,  was not found \n", ID)
			return models.User{} // Consider returning an error too
		}
		fmt.Printf("error scanning rows: %v", err)
		return models.User{} // consider returning an error too
	}
	return user
}

func AddActiveUser(user models.User) error {
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		fmt.Println("error opening DB", err)
		return err
	}
	defer db.Close()
	if UserWithEmailExists(user.Email) {
		err := fmt.Errorf("there is User registered with %s this email, pleas register with another email", user.Email)
		return err
	} else if UserWithUserNameExists(user.Username) {
		err := fmt.Errorf("there is User registered with %s this username, pleas register with another username", user.Username)
		return err
	} else {
		query := "INSERT INTO users (username,email,password,role,created_at,modified_at,active,uuid) VALUES (?,?,?,?,?,?,?,?)"
		_, err := db.Exec(query, user.Username, user.Email, user.Password, user.Role, user.CreatedAt, user.ModifiedAt, user.Active, user.UUID)
		if err != nil {
			err := fmt.Errorf("error adding User: %v, Error: %v", user.Username, err)
			return err
		}
	}
	return err
}

// Sets the user to inactive does  not remove the user from the DB
func InactiveActiveUser(user models.User) error {
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		fmt.Println("error opening DB", err)
	}
	defer db.Close()
	user.ModifiedAt = time.Now().Format("02/01/06,15/04")
	user.Active = config.INACTIVE
	query := "UPDATE users SET Active = ?, modified_at = ? WHERE uuid = ?"
	_, err = db.Exec(query, user.Active, user.ModifiedAt)
	if err != nil {
		fmt.Println("Error, updating user ", err)
		return err
	}
	fmt.Printf("Use  update %v to Inactive  \n", user)
	return nil
}

// ! not implemented Admin level function
func ActivateUser(user models.User) error {
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		fmt.Println("error opening DB", err)
	}
	defer db.Close()
	return nil
}

// When we want to update a user we will find the user that is logged in by their UUID
func FindUserByUUID(userUuid string) models.User {
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		fmt.Println("error opening DB", err)
	}
	defer db.Close()
	query := "SELECT id,username,email,password,role,created_at,modified_at,active,uuid FROM users  WHERE uuid = ?" // can be simplified with using * instead of specfing the columns
	row := db.QueryRow(query, userUuid)
	var user models.User
	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.ModifiedAt, &user.Active, &user.UUID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("There is no user with the UUID  %s, was not found \n", userUuid)
			return models.User{} // Consider returning an error too
		}
		fmt.Printf("error scanning rows: %v", err)
		return models.User{} // consider returning an error too
	}
	fmt.Printf("User %v \n:", user)
	return user
}

func AddNewUser(username string, email string, password string) error {
	var user models.User
	user.Username = username
	user.Email = email
	user.Password = password
	user.Role = "U"                                      // ? should we have this with integers as well ?
	user.CreatedAt = time.Now().Format("02/01/06,15/04") //
	user.ModifiedAt = time.Now().Format("02/01/06,15/04")
	user.Active = config.ACTIVE
	userUuid, err := GenerateUUID()
	if err != nil {
		return err
	}
	user.UUID = userUuid
	err = AddActiveUser(user)
	fmt.Printf("Error from Add ActiveUser %v  \n", err)
	if err != nil {
		return err
	}
	return nil
}

func ValidateUser(userName string, password string) (string, bool, error) {
	user := FindUserByUserName(userName)
	var uuid string
	if (user == models.User{}) {
		err := fmt.Errorf("user not found")
		return uuid, false, err
	} else if user.Active == config.INACTIVE {
		err := fmt.Errorf("account blocked, please contact the admin")
		return uuid, false, err
	} else if ValidateUserCredential(user.Password, password) {
		uuid = user.UUID
		return uuid, true, nil
	}
	return uuid, false, nil
}

func ValidateRegistrationOfUser(userName string, email string) {

}

func UserWithEmailExists(email string) bool {
	var exists bool
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		fmt.Println("error opening DB", err)
	}
	defer db.Close()
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE email = ?);`
	err = db.QueryRow(query, email).Scan(&exists)
	if exists {
		return true
	}
	fmt.Printf("The query gives us the following result %v ", exists)
	if err != nil {
		fmt.Errorf("error in the query to finding the mail from users", err)
		return true
	}
	return false
}
func OtherUserWithEmailExists(email string, id int) bool {
	var exists bool
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		fmt.Println("error opening DB", err)
	}
	defer db.Close()
	query := `SELECT * FROM users WHERE username = ? AND id != ?);`
	err = db.QueryRow(query, email, id).Scan(&exists)
	if exists {
		return true
	}
	fmt.Printf("The query gives us the following result %v ", exists)
	if err != nil {
		fmt.Errorf("error in the query to finding the mail from users", err)
		return true
	}
	return false
}

func UserWithUserNameExists(userName string) bool {
	var exists bool
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		fmt.Println("error opening DB", err)
	}
	defer db.Close()
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE username = ?);`
	err = db.QueryRow(query, userName).Scan(&exists)
	if exists {
		return true
	}
	fmt.Printf("The query gives us the following result %v ", exists)
	if err != nil {
		fmt.Errorf("error in the query to finding the username from users", err)
		return true
	}
	return false
}

func OtherUserWithUserNameExists(userName string, id int) bool {
	var exists bool
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		fmt.Println("error opening DB", err)
	}
	defer db.Close()
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE username = ? and id !=?);`
	err = db.QueryRow(query, userName, id).Scan(&exists)
	if exists {
		return true
	}
	fmt.Printf("The query gives us the following result %v ", exists)
	if err != nil {
		fmt.Errorf("error in the query to finding the username from users", err)
		return true
	}
	return false
}

func UpdateUserProfile(updatedUser models.User) (string, error) {
	var oldUser models.User
	var successMessage string
	var errorMessage error
	oldUser = FindUserByUserID(updatedUser.ID)
	if updatedUser.Password == "" {
		updatedUser.Password = oldUser.Password
	} else {
		updatedUser.Password = HashString(updatedUser.Password)
		successMessage += "Your password has been updated. \n"
		fmt.Println(updatedUser.Password)
	}
	// fmt.Printf("The old user %v \n", oldUser)
	// fmt.Printf("The new data %v \n ", updatedUser)
	if !OtherUserWithEmailExists(updatedUser.Email, oldUser.ID) {
		fmt.Println("User with Email does not exists \n")
		if updatedUser.Email != oldUser.Email {
			successMessage += "Your email has been updated. \n"
		}
	} else {
		errorMessage = fmt.Errorf("there is a user with the same email")
	}
	if !OtherUserWithUserNameExists(updatedUser.Username, oldUser.ID) {
		fmt.Println("User with same name does not  Exists ")
		if updatedUser.Username != oldUser.Username {
			successMessage += "Your username has been updated. \n"
		}
	} else {
		errorMessage = fmt.Errorf("there is another user with the same username")
	}
	if updatedUser.Role != oldUser.Role{
		successMessage += "Your Role has been updated. \n"
	} 
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		fmt.Println("error opening DB", err)
	}
	defer db.Close()
	updatedUser.Active = 1 // sets account to active can hardcode
	fmt.Println(updatedUser.Password)
	query := "UPDATE users SET username=? ,email=?, password=?, role=?, active = ?, modified_at = CURRENT_TIMESTAMP WHERE id = ?"
	_, err = db.Exec(query, updatedUser.Username, updatedUser.Email, updatedUser.Password, updatedUser.Role,updatedUser.Active,updatedUser.ID)
	if err != nil {
		fmt.Println("Error, updating user ", err)
		errorMessage = fmt.Errorf(err.Error()) 
		return successMessage , errorMessage
	}
	fmt.Println(successMessage)
	return successMessage, errorMessage
}
