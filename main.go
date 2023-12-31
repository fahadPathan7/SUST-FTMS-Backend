package main

import (
	"fmt"
	"ftms/controller"
	"ftms/router"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	controller.CreateDbConnection() // creating database connection

	fmt.Println("Server is running at port 5050...") // shows that server is running

	r := router.Router() // create router. it will be used to register routes.

	// Create a CORS handler with desired options.
	// it will allow api to be accessed from any origin
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // All origins
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	})

	// wrapping router with the CORS handler.
	// wrapping is done to allow api to be accessed from any origin
	handler := c.Handler(r)

	//connect  my frontend dir G:\SUST_FTMS_Frontend with my backend dir G:\SUST_FTMS_Backend also handle handler
	http.Handle("/", http.FileServer(http.Dir("G:\\SUST_FTMS_Frontend"))) // registering router with http Handle.
	http.Handle("/api/", handler) // registering router with http Handle.
	// it will handle all the incoming requests. "/" means all incoming requests.
	// second parameter is the router. here it is wrapped with CORS handler.

	http.ListenAndServe(":5050", nil) // this will start the server.
	// second parameter is the handler. nil means use default handler.
	// default handler is router. so it will use router to handle all the incoming requests.

	fmt.Println("Server is stopped!...") // shows that server is stopped
}