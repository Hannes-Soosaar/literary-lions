package handle

import (
	"fmt"
	"net/http"

	"gitea.kood.tech/hannessoosaar/literary-lions/pck/utils"
)

func verifySession(r *http.Request) string {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		err = fmt.Errorf("error %v \n", err)
		return err.Error()
	}
	sessionUser := utils.FindUserByUUID(cookie.Value)
	if sessionUser.Username != "" {
		return sessionUser.Username
	}
	return ""
}

func verifyPostMethod(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Request method", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func verifyGetMethod(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid Request method", http.StatusMethodNotAllowed)
		return false
	}
	return true
}
