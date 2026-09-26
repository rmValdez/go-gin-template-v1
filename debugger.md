# Go + Gin Debugging Steps

Here is the debugging sequence we went through for `go-gin-template-v1`, in order.

### 1. Check whether Go was installed
```powershell
go version
```
Initially, `go` was not recognized.

Then:
```powershell
where.exe go
```
Found: `C:\Program Files\Go\bin\go.exe`

After reopening the terminal:
```powershell
go version
```
Result: `go version go1.27.1 windows/amd64`

### 2. Make sure we were in the correct project directory
We accidentally ran Go commands from: `C:\WINDOWS\system32` which caused `Access is denied`.

Moved to:
```powershell
cd "C:\Users\My PC\Documents\Github\personal\template\integrated-repo\go\go-gin-template-v1"
```
Then verified with `Get-Location`.

### 3. Check the Go module
```powershell
Get-ChildItem
```
Confirmed: `README.md`, `go.mod`, `main.go`.

Then:
```powershell
go mod tidy
```
Initially, it said: `"all" matched no packages` because `main.go` didn't yet contain a real Go package/import.

### 4. Verify basic Go execution
We created:
```go
package main

import "fmt"

func main() {
    fmt.Println("go-gin-template-v1")
}
```

Then:
```powershell
go run .
```
Output: `go-gin-template-v1`
This confirmed the Go installation and module were working.

### 5. Add Gin
The Gin code imported:
```go
import (
    "net/http"
    "github.com/gin-gonic/gin"
)
```

Running `go run .` gave:
`no required module provides package github.com/gin-gonic/gin`

So we installed Gin:
```powershell
go get github.com/gin-gonic/gin
```

Then:
```powershell
go mod tidy
```
This installed Gin v1.12.0 and its dependencies.

### 6. First Gin execution problem — Windows Application Control
Running:
```powershell
go run .
```
produced: `An Application Control policy has blocked this file.`

**Important conclusion:** This was not a Gin or Go compilation error. Windows was blocking the temporary `.exe` that `go run` creates.

We checked generated executables with:
```powershell
Get-ChildItem *.exe
```
and later successfully ran a compiled executable, confirming the Go/Gin application itself could compile and execute.

### 7. Gin started but port 8080 was already occupied
Gin eventually reached:
`[GIN-debug] Listening and serving HTTP on :8080`

but then failed with:
`listen tcp :8080: bind: Only one usage of each socket address is permitted.`
Meaning another process was already listening on port 8080.

### 8. Find what was using port 8080
First:
```powershell
Get-NetTCPConnection -LocalPort 8080 -State Listen
```
Then:
```powershell
netstat -ano | findstr :8080
```
We found:
```text
TCP    0.0.0.0:8080    0.0.0.0:0    LISTENING    3500
TCP    [::]:8080       [::]:0       LISTENING    3500
```
So the process ID was `3500`.

Then:
```powershell
Get-Process -Id 3500
```
to identify the process.

### 9. Check for leftover IDE processes
Because Antigravity/VS Code is Electron-based, we also checked:
```powershell
Get-Process | Where-Object {
    $_.ProcessName -match "antigravity|code|electron"
} | Select-Object ProcessName, Id
```
There were many Code processes. We also checked/attempted to stop remaining Code processes when they were holding project files.

### 10. Final successful Gin startup
Eventually, the server successfully started:
```text
[GIN-debug] GET    /health
[GIN-debug] Listening and serving HTTP on :8080
```
So the backend was confirmed working.

The current minimal server is essentially:
```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()

    router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status": "ok",
        })
    })

    router.Run(":8080")
}
```

## Current status
| Check | Result |
|-------|--------|
| Go installed | ✅ |
| Go executable found | ✅ |
| Correct project directory | ✅ |
| go.mod | ✅ |
| go run basic app | ✅ |
| Gin installed | ✅ |
| go mod tidy | ✅ |
| Windows execution issue investigated | ✅ |
| Port 8080 conflict investigated | ✅ |
| Gin server | ✅ |
| /health endpoint | ✅ |

## Remaining Gin warnings
The server still showed:
`Running in "debug" mode`
and:
`You trusted all proxies`

Those are warnings, not failures.

For the reusable template, the next hardening steps should be:
- Configure `GIN_MODE`
- Configure trusted proxies explicitly
- Move `main.go` into `cmd/api`
- Add `internal/handlers`
- Add `internal/routes`
- Add configuration/env handling
- Add graceful shutdown
- Add tests

That would turn the current proof-of-concept into the actual reusable Go + Gin backend template.
