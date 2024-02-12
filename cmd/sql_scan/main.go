package main

import (
	"fmt"
	"github.com/SergeyAkaikin/vulnscan/internal/sqli/union_injection"
	"net/http"
)

func main() {
	url := "http://127.0.0.1:8000/searchproducts.php"
	cookie := http.Cookie{Name: "PHPSESSID", Value: "9d24c1a2db981724fdc266fc327805fd"}
	search := "searchitem=sss"
	u := union_injection.New(url, http.MethodPost, search)
	fmt.Println(u.Inject(&cookie))
	fmt.Println()

}
