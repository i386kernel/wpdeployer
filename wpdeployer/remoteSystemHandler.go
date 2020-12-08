package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
)

// Login to NFS
func sshLogin() *ssh.Client{
	config := &ssh.ClientConfig{
		User: deployprops.NFSUsername,
		Auth: []ssh.AuthMethod{
			ssh.Password(deployprops.NFSPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", deployprops.NFSIPaddr, config)
	if err != nil {
		fmt.Println(err)
	}
	return client
}

//Create a sharable directory in NFS
func ExecuteCommand(s *ssh.Client, runcmd string){
	sess, err := s.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	var b bytes.Buffer
	sess.Stdout = &b
	if err := sess.Run(runcmd); err != nil{
		panic("Failed to Execute Command: " + err.Error())
	}
}

//Get combined output
func GetCombinedOutput(s *ssh.Client, runcmd string) []byte{
	sess, err := s.NewSession()
	combinedbs, err := sess.CombinedOutput(runcmd)
	if err != nil {
		fmt.Println(err)
	}
	return combinedbs
}
