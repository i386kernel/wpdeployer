package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Openshift struct{
	Workloadname string
	URL string
	Token string
	DIRName string
}

func genrandom()int{
	rand.Seed(time.Now().Unix())
	//Generate a random number x where x is in range 30000<=x<=32767
	rangeLower := 30010
	rangeUpper := 32760
	randomNum := rangeLower + rand.Intn(rangeUpper-rangeLower+1)
	return randomNum
}

var randint = genrandom()
var projectname  = "wp-auto-"+strconv.Itoa(randint)

type wpprops struct{
	NFSIPaddr string
	NFSUsername string
	NFSPassword string
	PRURL string
	DRURL string
	PRToken string
	DRToken string
	PRsvc string
	DRsvc string
}

var PRWPdirname = projectname+"-PR"
var PRMysqlDirName = "mysql-"+projectname+"-PR"
var DRWPdirname = projectname+"-DR"
var DRMysqlDirName = "mysql-"+projectname+"-DR"

//var inidir = "/etc/wpdeployer/deployprop.ini"
var inidir = "C:\\Users\\LakshyaNanjangud\\go\\src\\github.com\\i386kernel\\GoExper\\Work-Load-Automation\\wpdeployer\\clouddeploy.ini"

var deployprops = wpprops{}

func init(){
	cfg, err := ini.Load(inidir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	deployprops.NFSIPaddr = cfg.Section("NFS").Key("ipaddr").String()+":"+"22"
	deployprops.NFSUsername = cfg.Section("NFS").Key("username").String()
	deployprops.NFSPassword = cfg.Section("NFS").Key("password").String()
	deployprops.PRURL = cfg.Section("OCP_Cluster").Key("prURL").String()
	deployprops.DRURL = cfg.Section("OCP_Cluster").Key("drURL").String()
	deployprops.PRToken = cfg.Section("OCP_Cluster").Key("prToken").String()
	deployprops.DRToken = cfg.Section("OCP_Cluster").Key("drToken").String()
	deployprops.PRsvc = cfg.Section("WP_URL").Key("prsvc").String()+":"+strconv.Itoa(randint)
	deployprops.DRsvc = cfg.Section("WP_URL").Key("drsvc").String()+":"+strconv.Itoa(randint)
}


func main(){
	tn := time.Now()

	sshsess := sshLogin()
	defer sshsess.Close()

	// NFS Directory and Service Ops

	fmt.Println("------------Performing NFS Operations----------")
	nfsPROps := NFSService{
		sshsess:   sshsess,
		WpDir:   PRWPdirname,
		MysqlDir: PRMysqlDirName,
		sourceURL: deployprops.DRsvc,
		destURL:   deployprops.PRsvc,
	}
	fmt.Println("Checking if the project exists")
	nfsPROps.CheckIfProjectExists()

	fmt.Printf("Proceeding with a new ---%s--- project\n", projectname)
	fmt.Println("Performing NFS Directory Operations")
	nfsPROps.NFSDirOperations()

	fmt.Println("Performing NFS Service Operations")
	nfsPROps.NFSServiceOperations()

	fmt.Println("Performing Wordpress helper service operations")
	nfsPROps.DBudpdaterops()

	nfsDROps := NFSService{
		sshsess:   sshsess,
		WpDir:   DRWPdirname,
		MysqlDir: DRMysqlDirName,
		sourceURL: deployprops.PRsvc,
		destURL:   deployprops.DRsvc,
	}
	nfsDROps.NFSDirOperations()
	nfsDROps.NFSServiceOperations()
	nfsDROps.DBudpdaterops()

	fmt.Println("---------Performing OCP operations----------")
	createProject(deployprops.PRURL, deployprops.PRToken)

	fmt.Println("MySql Operations")
	mySqlPR := Openshift{
		Workloadname: "mysql-" + projectname,
		URL:          deployprops.PRURL,
		Token:        deployprops.PRToken,
		DIRName:      PRMysqlDirName,
	}

	fmt.Println("Creating mysql Persistent Volume")
	mySqlPR.createPersistantVolume()
	fmt.Println("Creating Mysql Persistent volume claim")
	mySqlPR.createPersistantVolumeClaim()
	fmt.Println("Creating Mysql Deployment")
	mySqlPR.createMySqlDep()
	fmt.Println("Creating Mysql Service")
	mySqlPR.createMySqlService()

	wpPR := Openshift{
		Workloadname: "wordpress-" + projectname,
		URL:          deployprops.PRURL,
		Token:        deployprops.PRToken,
		DIRName:      PRWPdirname,
	}
	fmt.Println("Creating wordpress Persistent Volume")
	wpPR.createPersistantVolume()
	fmt.Println("Creating wordpress Persistent volume Volume")
	wpPR.createPersistantVolumeClaim()
	fmt.Println("Creating wordpress deployment")
	wpPR.createWPDep()
	fmt.Println("Creating wordpress service")
	wpPR.createWPService()

	// DR Volume Operations
	fmt.Println("Performing DR Volume Operations")
	fmt.Println("Creating MySql Persistent Volume")
	mySqlDR := Openshift{
		Workloadname: "mysql-" + projectname,
		URL:          deployprops.DRURL,
		Token:        deployprops.DRToken,
		DIRName:      DRMysqlDirName,
	}
	mySqlDR.createPersistantVolume()

	fmt.Println("Creating Wordpress Persistent Volume")
	wpDR := Openshift{
		Workloadname: "wordpress-" + projectname,
		URL:          deployprops.DRURL,
		Token:        deployprops.DRToken,
		DIRName:      DRWPdirname,
	}
	wpDR.createPersistantVolume()
	fmt.Println("Time Elapsed: " + (time.Since(tn)).String())
	fmt.Printf("Wordpress PR URL: %s\nWordpress DR URL: %s\n", deployprops.PRsvc, deployprops.DRsvc)
}


