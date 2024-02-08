package main

import (
	"fmt"
	"github.com/SergeyAkaikin/vulnscan/internal/sqli_scan_app/blind_injection"
	"net/http"
)

func main() {
	url := "http://127.0.0.1:8000/login1.php"
	uid := "uid=test"
	pass := "password=test"
	b := blind_injection.New(url, http.MethodPost, uid, pass)
	fmt.Println(b.TimeBased())
	fmt.Println(b.BooleanBased())

}
