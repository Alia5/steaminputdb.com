# Building

The building instructions are grouped for the Site itself (Backend/Frontend) and the Buddy-App.  

## Site

Building/Deving the site is only done in a Linux environment, so if you are running Microslops OS, just install the WSL environment.

### Requirements

- Go 1.26.2+
- NodeJS 25+
- Git (obviously)
- Protocol Buffers Compiler (protoc)  
  https://protobuf.dev/installation/
- [**Either/or**] PostgreSQL 18+ or Docker (for running PostgreSQL)
  - Use the compose.dev.yml file in the root of the repo to quickly spin up a compatible PostgreSQL instance.
- Optional but reccomended:  
  - air (for hot reloading)  
    go install github.com/air-verse/air@latest
  - just (for task running)  
    https://github.com/casey/just

### General

clone the repository and checkout submodules:  

```bash
git clone git@github.com:Alia5/steaminputdb.com.git
cd steaminputdb.com
git submodule update --init --recursive
```

#### Backend

Assuming you have all the dependencies installed.  

Spin up PostgreSQL (docker here)

```bash
docker compose -f compose.dev.yml up
```

Setup your .env file (copy .env.example and fill in the values):  

You will **need** a Steam API key to run the backend  
You can get one from Steam on: https://steamcommunity.com/dev/apikey

```bash
cd backend
cp .env.example .env
# Edit the .env file and insert your API KEY
```

Generate/Build the protobuf files (from the backend subdirectory):

```bash
chmod +x ./scripts/gen.sh
./scripts/gen.sh
```

You can then run the backend with (from the backend subdirectory):

```bash
go run ./cmd/steaminputdb/   
```

**Alternatively**, you can use air for hot reloading from the backend subdirectory by just running

```bash
air
```

#### Frontend

From the frontend subdirectory, install dependencies:  

```bash
cd frontend
npm i
```

You can then run the frontend (including hot reload) with: 

```bash
npm run dev
```

## Buddy-App

The Buddy-App does build on Windows, but it is recommended to use a linux environment for dev and just use GOOS=windows to build the app for windows.  
**If** you are running Microslops OS, just install the WSL environment.  

### Requirements

- Go 1.26.2+
- NodeJS 25+
- Optional but reccomended:  
  - air (for hot reloading)  
    go install github.com/air-verse/air@latest
  - just (for task running)  
    https://github.com/casey/just

### General

clone the repository and checkout submodules:  

```bash
git clone git@github.com:Alia5/steaminputdb.com.git
cd steaminputdb.com
git submodule update --init --recursive
```

The buddy-app has Javascript and Go components.  
The JS-Parts get injected into a running Steam client and are embedded in the Go binary, thus they will need to be built first!  

### Building the JS-Parts

```bash
cd buddy-app/steam/steam_cef/steam_js/templates
npm i
npm run build
```

### Building the Go Binary

Generate/Build the protobuf files (from the backend subdirectory):

```bash
chmod +x ./scripts/gen.sh
./scripts/gen.sh
```

Create a .env file in the buddy-app subdirectory with the following content:

```env
DEV=1
STEAMINPUTDB_BUDDY_CORS_ORIGINS=*
LOG_LEVEL=debug
```

You need to enable _CEF Remote Debugging_ in the Steam Client to allow the buddy-app to connect to it.  
To do that, create a file named `.cef-enable-remote-debugging` in your Steam directory (where the Steam binary is located), then restart Steam.  
(Note that Steam uses TCP port 8080 on localhost, **with no easy way to change it**, so make sure that port is already used by another application...)

You can then run the buddy-app with (from the buddy-app subdirectory):

```bash
go run ./cmd/steaminputdb-buddy 
```

**alternatively**, and reccomended you can use air for hot reloading from the buddy-app subdirectory by just running

```bash
air
```

To build the buddy-app binary without running it use

```bash
go build -o dist/steaminputdb-buddy ./cmd/steaminputdb-buddy
```

If you want to build the buddy-app binary for windows from a linux environment (eg. WSL), run

```bash
GOOS=windows GOARCH=amd64 go build -o dist/steaminputdb-buddy.exe ./cmd/steaminputdb-buddy
```

#### Windows users

The air script will even work when run from Windows Powershell, even if the repo is checked out inside WSL.  
Assuming you have `just` installed on Windows.

```pwsh
cd \\wsl.localhost\<DISTRO_NAME>\<REPO_DIRECTORY>\buddy-app
air
```
