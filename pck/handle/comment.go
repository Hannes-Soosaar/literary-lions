package handle

import (
	"fmt"
	"net/http"
)


func CommentHandler( w http.ResponseWriter, r *http.Request){

fmt.Println("Comment handler started")

}