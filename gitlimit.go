package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: gitlimit repository [repository ...]\n")
	}
	commandLine := os.Getenv("SSH_ORIGINAL_COMMAND")
	args := strings.Fields(commandLine)
	if len(args) != 2 {
		log.Fatalf("Invalid git command %q\n", commandLine)
	}
	if args[0] != "git-upload-pack" && args[0] != "git-receive-pack" {
		log.Fatalf("Command not permitted %q\n", commandLine)
	}
	if len(args[1]) < 2 || args[1][0] != '\'' || args[1][len(args[1])-1] != '\'' {
		log.Fatalf("Invalid repository %q\n", args[1])
	}
	repository := args[1][1:len(args[1])-1]
	ok := false
	for _, permittedRepository := range os.Args[1:] {
		if repository == permittedRepository {
			ok = true
			break
		}
	}
	if !ok {
		log.Fatalf("Access to repository %q is denied", repository)
	}
	cmd := exec.Command("git-shell", "-c", fmt.Sprintf("%s '%s'", args[0], repository))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
