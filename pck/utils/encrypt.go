package utils

// Add logic to encrypt pw
import(

	"fmt"
	"crypto/sha256"
)

func HashString(s string)string{
	hash:= sha256.New()
	hash.Write([]byte(s))
	hashedBytes :=hash.Sum(nil)
	hashedString :=fmt.Sprintf("%x",hashedBytes)
	fmt.Printf("the hashed String is %s", hashedString )
	return hashedString
}

