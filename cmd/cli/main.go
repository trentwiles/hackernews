package main

import (
	"flag"
	"fmt"
	//"log"
	"os"
)

var USAGE_TEST string = fmt.Sprintf(`
NAME: hackernews admin CLI

USAGE: %s [options]

COMMANDS:
	--username		REQUIRED - username to grant/revoke admin privileges on
	--remarks		OPTIONAL - comments on user that will be entered into the database
	--revoke		OPTIONAL - set to true to remove an existing admin's privledges
`, os.Args[0])

func main() {
	username := flag.String("username", "", "Username to grant/revoke admin privledges from")
	// remarks := flag.String("remarks", "", "Optional comments to insert into database?")
	// revoke := flag.Bool("revoke", false, "Revoke privledges?")

	flag.Parse()

	if *username == "" {
		fmt.Println("error: `--username` is required")
		fmt.Println()
		fmt.Println(USAGE_TEST)
	}
}
