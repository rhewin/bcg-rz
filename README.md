# Golang Project
Implementation of simple shopping cart with Promotions using Golang.

Folder Structure:
- build: contains ready-to-run binary
- schema-app: GraphQL schema
- src: raw project's source code 
  - assets: contains static files like image, css and js
  - db: contains json files that act as a data supplier serve from API
  - views: contains html templates


## How to Run
Please make sure your system has golang installed. To do so you can refer to this link https://golang.org/doc/install. 

  - Open terminal
  - Go to [project_folder]/src and run in terminal: $ go run main.go
  
When running, the project will take port 8081 to run the server. Open the browser and locate url http://localhost:8081/
(*Make sure there is no other application that running/use port 8081. If so, you can temporarily shutdown the apps that run/use port 8081)


## How to Test
Go to project folder, run from terminal: $ go test -v
