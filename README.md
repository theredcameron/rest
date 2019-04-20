Go REST
=======

## Summary
Go REST is designed for providing a library for creating a REST API client in go.  
It includes support for easy session and cookie management.

## Example
```go
package main

import (
	"fmt"

	"github.com/theredcameron/rest"
)

type Runner struct {
	Run  string `json:"run"`
	Away string `json:"away"`
	Fast string `json:"fast"`
}

func main() {
	endpoints := rest.Endpoints{
		{
			Description: "Test description",
			Method:      "GET",
			Path:        "/api/{id}",
			F:           Test,
		},
	}

	meta := &rest.CookieMeta{
		File:      "test",
		Path:      "./",
		MaxAge:    3000,
		Key:       []byte("key"),
		StoreName: "test",
	}

	router, err := rest.NewRouter(endpoints, meta)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Starting on port 4000")
	err = router.Start("4000")
	if err != nil {
		fmt.Println(err)
	}
}

func Test(r *rest.Request) (interface{}, error) {
	run := Runner{
		Run:  "running",
		Away: r.Params["away"],
		Fast: r.Vars["id"],
	}

	if r.GetCookieValue("testing") == nil {
		r.SetCookieValue("testing", "run")
	}

	return run, nil
}

```