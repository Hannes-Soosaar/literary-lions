package utils

import (
	"database/sql"
	"fmt"

	// "gitea.kood.tech/hannessoosaar/literary-lions/intenal/config"
	"gitea.kood.tech/hannessoosaar/literary-lions/intenal/config"
	"gitea.kood.tech/hannessoosaar/literary-lions/pck/models"
)

// ? perhaps we are getting too much information
func FindUserByUserName(userName string) models.User {

	var user models.User
	// See if we can make the open db into a separate function so we do not need to open close the DB for every request
	db, err := sql.Open("sqlite3", config.LION_DB)  // ! the path is different here than it is for the other 
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
		fmt.Printf("User %v \n:",user)
	return user
}

func AddActiveUser( user models.User)error{

	db, err := sql.Open("sqlite3", config.LION_DB)  // ! the path is different here than it is for the other 
	if err != nil {
		fmt.Println("error opening DB", err)
		return err // should return an empty use
	}
	defer db.Close()
	
	dbUser := FindUserByUserName(user.Username) // retuns a model.User instance

	if (dbUser == models.User{}) {
	

	} else{
		fmt.Printf("There is User %s by this name!", dbUser.Username)
		return nil
	}
	//? Should we check if there is the user is also active and inactive
	fmt.Printf("User %s added", user.Username)

	return nil
}

func DeleteActiveUser(userID int)error{
// Find User by ID and changes the user to not active
return nil
}




func  ValidateUser(userName string) error{

	user := FindUserByUserName(userName)
	if user.Active == 1 {
		return nil
	}

	return nil
}
