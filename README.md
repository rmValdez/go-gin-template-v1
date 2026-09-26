# Go Gin Template V1

Here is the complete Windows + IDE integrated terminal setup from installing Go to running your first Go project.

## 1. Install Go
Go to the official Go website:
[Go Downloads](https://go.dev/dl/?utm_source=chatgpt.com)
Download the Windows AMD64 MSI installer.
For a typical Windows 11 PC, choose something like:
`go1.xx.x.windows-amd64.msi`
Run the installer.
Use the default installation location:
`C:\Program Files\Go`
Click through the installer until it finishes.

## 2. Verify Go installation
Close your IDE/terminal if it was already open.
Then reopen your IDE.
Open the IDE's integrated terminal:
**Terminal → New Terminal**
You should get PowerShell, for example:
`PS C:\Users\My PC>`
Run:
```bash
go version
```
Expected:
`go version go1.27.1 windows/amd64`
*(Your machine already reached this stage successfully.)*

You can also verify where Go is installed:
```bash
where.exe go
```
Expected:
`C:\Program Files\Go\bin\go.exe`

## 3. Create your project folder
For example:
`C:\Users\My PC\Documents\Github\personal\template\integrated-repo\go-gin-template-v1`
You can create/open this folder directly in your IDE.
For VS Code:
**File → Open Folder**
Select:
`go-gin-template-v1`

## 4. Open the IDE terminal
In VS Code:
**Terminal → New Terminal**
You should see something like:
`PS C:\Users\My PC\Documents\Github\personal\template\integrated-repo\go-gin-template-v1>`
You can verify your current directory:
```bash
pwd
```

## 5. Initialize the Go module
Run:
```bash
go mod init github.com/rmValdez/go-gin-template-v1
```
Go creates:
`go.mod`
Your project now looks like:
```
go-gin-template-v1/
└── go.mod
```
The `go.mod` contains your module name:
`module github.com/rmValdez/go-gin-template-v1`

## 6. Create main.go
In the IDE Explorer:
Right-click `go-gin-template-v1` → **New File** → `main.go`
Put this inside:
```go
package main

import "fmt"

func main() {
    fmt.Println("go-gin-template-v1")
}
```
Save:
`Ctrl + S`
Your project now looks like:
```
go-gin-template-v1/
├── go.mod
└── main.go
```

## 7. Run the Go application
In the same IDE integrated terminal:
```bash
go run .
```
You should get:
`go-gin-template-v1`

That's the basic Go development cycle:
Edit code ↓ Save ↓ `go run .` ↓ Application runs

## 8. Run go mod tidy
Now that you have Go code:
```bash
go mod tidy
```
This makes sure your module dependencies are correct.
For this simple project, there may be nothing to download because `fmt` is part of Go's standard library.

## 9. Build the application
Run:
```bash
go build .
```
This compiles your application.
On Windows you'll typically get:
`go-gin-template-v1.exe`
So:
```
go-gin-template-v1/
├── go.mod
├── main.go
└── go-gin-template-v1.exe
```
You can run it:
```bash
.\go-gin-template-v1.exe
```
Output:
`go-gin-template-v1`

## 10. Recommended Git setup
Since this is going into your GitHub template repository:
```bash
git init
```
Create `.gitignore`:
```
go-gin-template-v1.exe
```
Then:
```bash
git add .
git commit -m "chore: initialize Go template"
```

## Your complete workflow
From now on, when starting a Go project, the basic process is:
```bash
# Open project in IDE
# Initialize module
go mod init github.com/rmValdez/project-name

# Download/clean dependencies
go mod tidy

# Run
go run .

# Build
go build .

# Test
go test ./...
```

And your IDE handles the actual editing:
```
IDE
 │
 ├── main.go
 ├── go.mod
 ├── internal/
 ├── cmd/
 └── ...
 │
 ▼ Integrated PowerShell
 │
 ├── go run .
 ├── go test ./...
 ├── go build .
 └── go mod tidy
```
You do not need a separate Command Prompt or PowerShell window. The integrated terminal in your IDE is enough. Your current `go-gin-template-v1` project has already successfully completed the installation, module initialization, and `go run .` stages.