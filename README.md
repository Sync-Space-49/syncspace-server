# SyncSpace Server
## What Does it Do 🤔
Syncspace backend is the REST API that handles all communication between the frontend and backend services. It handles all data storage and retrieval for the project. It also handles authentication and authorization for the project ultilizing JSON Web Tokens.


## 👩‍💻 Local Development
### Installation ⚒️
To run this project, you will need to have [Go 1.20](https://go.dev/dl/) installed. This also hooks into a [postgres](https://www.postgresql.org/download/) database to store data.

To install this project, simply clone the repository using git:

For SSH (Recommended):

`git clone git@github.com:Sync-Space-49/syncspace-server.git`

For HTTPS:

`git clone https://github.com/Sync-Space-49/syncspace-server.git`


### Configuration ⚙️
When initally getting setup you will need to create a `.env` file in the root directory based on the `.env.sample` file. The main portion you will need to setup for local development is your Postgres password and DB name.


### Running 🚀
You can download the project's Go dependencies using the `go get` command. To run the project, use `go run main.go`; this will spin up a server on the url specified in `API_HOST`. Whenever you make changes to the code, you will need to restart the server (ctrl+c in the terminal kills the current process) to see the changes. When making changes to dependencies, you will need to run `go mod tidy` to update the `go.mod` file then use `go get -u` to fetch the latest versions of the dependencies listed in the `go.mod` file.

### Compiling 🏗️
To compile the project to a binary executable, you can run `go build` in the root directory. This is useful for deploying the project but not usually needed for local development.


### Tests 🧪
TODO


### Deployment 📦
TODO