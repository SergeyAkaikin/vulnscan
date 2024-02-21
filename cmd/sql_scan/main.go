package main

import (
	"fmt"
	"github.com/SergeyAkaikin/vulnscan/internal/sqli/error_injection"
	"net/http"
)

func main() {
	url := "http://127.0.0.1:8000/login1.php"
	//cookie := http.Cookie{Name: "PHPSESSID", Value: "9d24c1a2db981724fdc266fc327805fd"}
	uid := "uid=sss"
	password := "password=test"
	u := error_injection.New(url, http.MethodPost, uid, password)
	fmt.Println(u.Inject(nil))
	fmt.Println()

}
