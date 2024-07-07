package utils

import (
	"database/sql"
	"fmt"
	"time"

	"gitea.kood.tech/hannessoosaar/literary-lions/intenal/config"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/models"
)

func FindUserByUserName(userName string) models.User {
	var user models.User
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		fmt.Println("error opening DB", err)
		return user
	}
	defer db.Close()
	query := "SELECT id, username,email,password,role,created_at,modified_at,active,uuid FROM users  WHERE username = ?"
	row := db.QueryRow(query, userName)
	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.ModifiedAt, &user.Active, &user.UUID)
	if err != nil {
		if err == sql.ErrNoRows {
			models.GetInstance().SetError(err)
			return models.User{}
		}
		fmt.Printf("error scanning rows: %v", err)
		models.GetInstance().SetError(err)
		return models.User{}
	}
	return user
}
func FindUserByUserID(ID int) models.User {
	var user models.User
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		fmt.Println("error opening DB", err)
		models.GetInstance().SetError(err)
		return user
	}
	defer db.Close()
	query := "SELECT id, username,email,password,role,created_at,modified_at,active,uuid FROM users  WHERE ID = ?"
	row := db.QueryRow(query, ID)
	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.ModifiedAt, &user.Active, &user.UUID)
	if err != nil {
		if err == sql.ErrNoRows {
			models.GetInstance().SetError(err)
			return models.User{}
		}
		fmt.Printf("error scanning rows: %v", err)
		models.GetInstance().SetError(err)
		return models.User{}
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
		models.GetInstance().SetError(err)
		return err
	} else if UserWithUserNameExists(user.Username) {
		err := fmt.Errorf("there is User registered with %s this username, pleas register with another username", user.Username)
		models.GetInstance().SetError(err)
		return err
	} else {
		fmt.Println("Running the query to add the user.")
		query := "INSERT INTO users (username,email,password,role,created_at,modified_at,active,uuid) VALUES (?,?,?,?,?,?,?,?)"
		_, err := db.Exec(query, user.Username, user.Email, user.Password, user.Role, user.CreatedAt, user.ModifiedAt, user.Active, user.UUID)
		if err != nil {
			err := fmt.Errorf("error adding User: %v, Error: %v", user.Username, err)
			models.GetInstance().SetError(err)
			return err
		}
	}

	fmt.Println(err)
	return err
}

func FindUserByUUID(userUuid string) models.User {
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		fmt.Println("error opening DB", err)
		models.GetInstance().SetError(err)
	}
	defer db.Close()
	query := "SELECT id,username,email,password,role,created_at,modified_at,active,uuid FROM users  WHERE uuid = ?"
	row := db.QueryRow(query, userUuid)
	var user models.User
	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.ModifiedAt, &user.Active, &user.UUID)
	if err != nil {
		if err == sql.ErrNoRows {
			models.GetInstance().SetError(err)
			return models.User{}
		}
		fmt.Printf("error scanning rows: %v", err)
		models.GetInstance().SetError(err)
		return models.User{}
	}
	fmt.Printf("User %v \n:", user)
	return user
}

func AddNewUser(username string, email string, password string) error {
	var err error 
	var user models.User
	user.Username = username
	user.Email = email
	user.Password = password
	user.Role = "U"
	user.CreatedAt = time.Now().Format("02/01/06,15/04")
	user.ModifiedAt = time.Now().Format("02/01/06,15/04")
	user.Active = config.ACTIVE
	userUuid, err := GenerateUUID()
	if err != nil {
		return err
	}
	user.UUID = userUuid
	err = AddActiveUser(user)
	if err != nil {
		return err
	}
	return err
}

func ValidateUser(userName string, password string) (string, bool, error) {
	user := FindUserByUserName(userName)
	var uuid string
	if (user == models.User{}) {
		err := fmt.Errorf("%s is not a user, please register", userName)
		models.GetInstance().SetError(err)
		return uuid, false, err
	} else if user.Active == config.INACTIVE {
		err := fmt.Errorf("account blocked, please contact the admin")
		models.GetInstance().SetError(err)
		return uuid, false, err
	} else if ValidateUserCredential(user.Password, password) {
		uuid = user.UUID
		return uuid, true, nil
	} else {
		err := fmt.Errorf("the entered password is not correct")
		models.GetInstance().SetError(err)
	}
	return uuid, false, nil
}

func UserWithEmailExists(email string) bool {
	var exists bool
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		models.GetInstance().SetError(err)
	}
	defer db.Close()
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE email = ?);`
	err = db.QueryRow(query, email).Scan(&exists)
	if exists {
		return true
	}
	if err != nil {
		err = fmt.Errorf("error in the query to finding the mail from users %v ", err)
		models.GetInstance().SetError(err)
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
	query := `(SELECT * FROM users WHERE username = ? AND id != ?);`
	err = db.QueryRow(query, email, id).Scan(&exists)
	if exists {
		err = fmt.Errorf("the email %s has been taken, choose another", email)
		models.GetInstance().SetError(err)
		return true
	}
	if err != nil {
		err = fmt.Errorf("error in the query to finding the mail from users", err)
		models.GetInstance().SetError(err)
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
		err = fmt.Errorf("th username %s is taken please choose another", userName)
		models.GetInstance().SetError(err)
		return true
	}
	if err != nil {
		err = fmt.Errorf("error in the query to finding the username from users, %v", err)
		models.GetInstance().SetError(err)
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
	if err != nil {
		err = fmt.Errorf("error in the query to finding the username from users %v", err)
		models.GetInstance().SetError(err)
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
	}
	if updatedUser.Email != oldUser.Email {
		successMessage += "Your email has been updated. \n"
	}
	if !OtherUserWithUserNameExists(updatedUser.Username, oldUser.ID) {
		if updatedUser.Username != oldUser.Username {
			successMessage += "Your username has been updated.\n"
		}
	} else {
		errorMessage = fmt.Errorf("there is another user with the same username")
		models.GetInstance().SetError(errorMessage)
	}
	db, err := sql.Open("sqlite3", config.LION_DB)
	if err != nil {
		fmt.Println("error opening DB", err)
	}
	defer db.Close()
	query := "UPDATE users SET username=? ,email=?, password=?, role=?, active = 1, modified_at = CURRENT_TIMESTAMP WHERE id = ?"
	_, err = db.Exec(query, updatedUser.Username, updatedUser.Email, updatedUser.Password, updatedUser.Role, updatedUser.ID)
	if err != nil {
		errorMessage = fmt.Errorf("error, updating user %v ", err)
		models.GetInstance().SetError(errorMessage)
		return successMessage, errorMessage
	}
	models.GetInstance().SetSuccess(successMessage)
	models.GetInstance().SetError(errorMessage)
	return successMessage, errorMessage
}
