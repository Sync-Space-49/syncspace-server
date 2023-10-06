# SyncSpace Server
## What Does it Do 🤔
Syncspace backend is the REST API that handles all communication between the frontend and backend services. It handles all data storage and retrieval for the project. It also handles authentication and authorization for the project ultilizing JSON Web Tokens.
## 👩‍💻 Usage
### Installation ⚒️
To run this project, you will need to have [go 1.20](https://go.dev/dl/) installed. This also hooks into a [postgres](https://www.postgresql.org/download/) database to store data.

To install this project, simply clone the repository using ```git clone git@github.com:Sync-Space-49/syncspace-server.git``` for SSH (recommended) or ```git clone https://github.com/Sync-Space-49/syncspace-server.git``` for HTTPS.

### Configuration ⚙️
When initally getting setup you will need to create a `.env` file in the root directory based on the `.env.sample` file. The main portion you will need to setup for local development is your Postgres password and DB name.
### Running 🚀
After installing the project, dependancies, and setting up enviroment variables, you can run the project using `go run main.go` in the root directory. This will start the server on the url specified in `API_HOST`. Whenever you make changes to the code, you will need to restart the server (ctrl+c in the terminal kills the current process) to see the changes.

### Tests 🧪
TODO
### Deployment 📦
TODO