package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"

	"golang.org/x/crypto/ssh"
)

func main() {
	homePath := os.Getenv("HOME")
	if homePath == "" {
		log.Fatal("please set your $HOME path first")
	}
	sshPrivatePath := path.Join(homePath, ".ssh/id_rsa")

	var auth []ssh.AuthMethod
	auth = []ssh.AuthMethod{PublicKeyFile(sshPrivatePath)}

	sshConfig := &ssh.ClientConfig{
		User: "aw160502",
		Auth: auth,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	connection, err := ssh.Dial("tcp", "47.74.179.121:22", sshConfig)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to dial ssh server %v", err))
	}

	session, err := connection.NewSession()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create new session %v", err))
	}

	out, err := session.CombinedOutput("ls -l")
	if err != nil {
		log.Fatal(fmt.Errorf("failed to exec cmd %v", err))
	}
	log.Println(string(out))

	session.Close()
	log.Println("ALL SUCCESS")
}

// PublicKeyFile function
func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}
