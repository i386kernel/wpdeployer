package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"strings"
)


type NFSService struct {
	sshsess   *ssh.Client
	WpDir     string
	MysqlDir  string
	sourceURL string
	destURL   string
}

func (n *NFSService) CheckIfProjectExists(){
	combinedbs := GetCombinedOutput(n.sshsess, "cd /mnt; ls;")
	strbs := string(combinedbs)
	strsplit := strings.Split(strbs, "\n")

	for _, v := range strsplit{
		if v == "mysql-"+projectname+"-PR"{
			fmt.Printf("The project %s already exists, EXITING....", projectname)
			fmt.Println("RUN THE PROGRAM AGAIN!!!!")
			os.Exit(1)
		}
	}
}

func (n *NFSService) NFSDirOperations(){
	ExecuteCommand(n.sshsess, "cd /mnt/;" + "mkdir " + n.WpDir + "; chmod -R 777 *")
	ExecuteCommand(n.sshsess, "cd /mnt/;" + "mkdir " + n.MysqlDir + "; chmod -R 777 *")
}

func (n *NFSService) DBudpdaterops(){
	yamlrun1 := fmt.Sprintf("echo '/mnt/%s/ %s %s' >> /etc/yamlupdater/yamlprop.txt", n.WpDir, n.sourceURL, n.destURL)
	ExecuteCommand(n.sshsess, yamlrun1)
}


func (n *NFSService) NFSServiceOperations() {
	// Exporting Directories
	nfsexport1 := fmt.Sprintf("echo '/mnt/%s *(rw,sync,no_root_squash,insecure)' >> /etc/exports", n.WpDir)
	ExecuteCommand(n.sshsess, nfsexport1)
	nfsexport2 := fmt.Sprintf("echo '/mnt/%s *(rw,sync,no_root_squash,insecure)' >> /etc/exports", n.MysqlDir)
	ExecuteCommand(n.sshsess, nfsexport2)
	// Restarting NFS Service
	ExecuteCommand(n.sshsess, "systemctl restart nfs")
}
