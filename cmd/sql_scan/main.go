package main

import (
	"fmt"
	"github.com/SergeyAkaikin/vulnscan/internal/sqli_scan_app/blind_injection"
	"net/http"
)

func main() {
	url := "http://127.0.0.1:8000/login1.php"
	qParams := []string{"uid", "password"}
	method := http.MethodPost

	inj := blind_injection.New(url, method, qParams...)
	injectable, err := inj.InjectParam()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(injectable)
}
