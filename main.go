package main

import (
	"fmt"

	"github.com/coreos/go-oidc"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	fmt.Println(provider, err)
	fmt.Println("yea")
}
