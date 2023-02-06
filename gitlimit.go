package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	mRead = 1 << iota
	mWrite
)

var modeParse = map[string]int{
	"r": mRead, "w": mWrite,
	"rw": mRead | mWrite, "wr": mRead | mWrite,
}

var mFlag = flag.String("m", "rw", "Mode (r / w / rw)")
var cFlag = flag.Bool("c", false, "Check access / dry run mode")

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		log.Fatalf("Usage: gitlimit [flags] repository [repository ...]\n")
	}
	mode, ok := modeParse[*mFlag]
	if !ok {
		log.Fatalf("Mode (-m) flag must be r / w / rw")
	}
	commandLine := os.Getenv("SSH_ORIGINAL_COMMAND")
	args := strings.Fields(commandLine)
	if len(args) != 2 {
		log.Fatalf("Invalid git command %q\n", commandLine)
	}
	switch args[0] {
	case "git-upload-pack":
		if (mode & mRead) == 0 {
			log.Fatalf("Download not permitted %q\n", commandLine)
		}
	case "git-receive-pack":
		if (mode & mWrite) == 0 {
			log.Fatalf("Upload not permitted %q\n", commandLine)
		}
	default:
		log.Fatalf("Command not permitted %q\n", commandLine)
	}
	if len(args[1]) < 2 || args[1][0] != '\'' || args[1][len(args[1])-1] != '\'' || args[1][1] == '-' {
		log.Fatalf("Invalid repository %q\n", args[1])
	}
	repository := args[1][1 : len(args[1])-1]
	ok = false
	for _, permittedRepository := range flag.Args() {
		if repository == permittedRepository {
			ok = true
			break
		}
	}
	if !ok {
		log.Fatalf("Access to repository %q is denied", repository)
	}
	if *cFlag {
		log.Printf("%s to %s is permitted", args[0], repository)
	} else {
		cmd := exec.Command("git-shell", "-c", fmt.Sprintf("%s '%s'", args[0], repository))
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}
