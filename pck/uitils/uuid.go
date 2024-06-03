package utils

import (
	"crypto/rand"
	"fmt"
)


// This is the first iteration of the UUID  
func GenerateUUID() (string,error){
    u := make([]byte, 16)
    if _, err := rand.Read(u); err != nil {
        return "", err
    }

    // Set the version to 4
    u[6] = (u[6] & 0x0f) | 0x40
    // Set the variant to RFC 4122
    u[8] = (u[8] & 0x3f) | 0x80

    return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
        u[0:4], u[4:6], u[6:8], u[8:10], u[10:]), nil
}

