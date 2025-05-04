# Login

The following demonstrates how to

## E-Mail + Password Login

Username + Password login was implemented.

```go
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/StarHack/go-internxt-drive/auth"
	"github.com/StarHack/go-internxt-drive/config"
	"github.com/StarHack/go-internxt-drive/users"
)

func main() {
	cfg := config.NewDefault("user@example.com", "super_secret_password_123")

	// /auth/login
	loginResp, err := auth.Login(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "login error: %v\n", err)
		os.Exit(1)
	}
	loginJSON, _ := json.MarshalIndent(loginResp, "", "  ")
	fmt.Println("Login response:")
	fmt.Println(string(loginJSON))

	// /auth/login/access
	accessResp, err := auth.AccessLogin(cfg, loginResp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "access error: %v\n", err)
		os.Exit(1)
	}
	accessJSON, _ := json.MarshalIndent(accessResp, "", "  ")
	fmt.Println("Access response:")
	fmt.Println(string(accessJSON))

	// Token is required for all other API calls:
	fmt.Println("Bearer token:")
	fmt.Println(accessResp.NewToken)

	usage, err := users.GetUsage(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "usage error: %v\n", err)
		os.Exit(1)
	}

	limit, err := users.GetLimit(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "limit error: %v\n", err)
	}

	fmt.Printf("Account usage: %d/%d bytes\n", usage.Drive, limit.MaxSpaceBytes)
}

```

## Token Login

If we already have a Bearer token from a previous login we may just use that.

```go
package main

import (
	"fmt"
	"os"

	"github.com/StarHack/go-internxt-drive/config"
	"github.com/StarHack/go-internxt-drive/users"
)

func main() {
	cfg := config.NewDefaultToken("eyXXXXXXXXXXXXXXXXXXXX.XXXXXXXXXXXXXXXXXXXX.XXXXXXXXXXXXXXXXXXXX")

	usage, err := users.GetUsage(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "usage error: %v\n", err)
		os.Exit(1)
	}

	limit, err := users.GetLimit(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "limit error: %v\n", err)
	}

	fmt.Printf("Account usage: %d/%d bytes\n", usage.Drive, limit.MaxSpaceBytes)
}

```
