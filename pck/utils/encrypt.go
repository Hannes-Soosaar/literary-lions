package utils

// Add logic to encrypt pw
import(
	"golang.org/x/crypto/bcrypt"
)

func HashString(s string)string{
	hashedBytes,err :=bcrypt.GenerateFromPassword([]byte(s),bcrypt.DefaultCost)
	if err!= nil {
		return ""
	}
	hashedString:=string(hashedBytes)
	return hashedString
}

