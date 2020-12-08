package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type DBfile struct {
	SourceURL  string `yaml:"sourceURL"`
	DestURL    string `yaml:"destURL"`
	DbURL      string `yaml:"dbURL"`
	DbUsername string `yaml:"dbUsername"`
	DBPassword string `yaml:"dbPassword"`
	DBName     string `yaml:"dbName"`
	SSHKey string `yaml:"sshKey"`
	OpenShiftURL string `yaml:"openShiftURL"`
}

var filedata DBfile
var file = "/dbfile/dbdata.yml"

func fileops(filename string) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File does not exist. Trying again.....")
			time.Sleep(10 * time.Second)
			fileops(file)
		}
	}
	log.Printf("File %s, Proceeding", fileInfo.Name())
	rbs, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	err = yaml.Unmarshal(rbs, &filedata)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Current DB YAML Data: %s\n", filedata)
}

func changesiteURL() (string, error){
	var db *sql.DB
	dbConnect := fmt.Sprintf("%s:%s@(%s)/%s", filedata.DbUsername, filedata.DBPassword,
		filedata.DbURL, filedata.DBName)
	db, err := sql.Open("mysql", dbConnect)
	if err != nil {
		fmt.Println("DB Connection Failed", err)
	}
	db.SetConnMaxLifetime(10 * time.Second)
	db.SetMaxIdleConns(3)
	fmt.Println("Trying to connect to the specified DB......")
	err = db.Ping()
	if err != nil {
		fmt.Println("Unable to Contact the Specified DB/Tables in Database YAML file")
	}

	//querystmt := fmt.Sprintf("SHOW TABLES")
	//dbrw, err := db.Query(querystmt)
	//if err != nil {
	//	fmt.Println(err)
	//}
	////dbrw

	stmt1 := fmt.Sprintf("UPDATE wp_options SET option_value = replace(option_value, '%s', '%s') " +
		"WHERE option_name = 'home' OR option_name = 'siteurl'", filedata.SourceURL, filedata.DestURL)

	db.Query(stmt1)

	stmt2 := fmt.Sprintf("UPDATE wp_posts SET guid = replace(guid, '%s', '%s')", filedata.SourceURL,
		filedata.DestURL)
	stmt3 := fmt.Sprintf("UPDATE wp_posts SET post_content = replace(post_content, '%s', '%s')",
		filedata.SourceURL, filedata.DestURL)
	stmt4 := fmt.Sprintf("UPDATE wp_postmeta SET meta_value = replace(meta_value, '%s', '%s')",
		filedata.SourceURL, filedata.DestURL)

	statements := []string{stmt1, stmt2, stmt3, stmt4}

	fmt.Println("Trying to manipulate DB Entries")

	for _, item := range statements {
		res, err := db.Exec(item)
		if err != nil{
			db.Close()
			fmt.Println(err)
		}
		fmt.Println(res.RowsAffected())
	}
	fmt.Println("Closing Connection to DB")
	err = db.Close()
	if err != nil {
		fmt.Println(err)
	}
	return "Successfully Changed DB Entries", nil
}

func main() {
	for i := 0; i < 10; i++{
		fileops(file)
		reply, err := changesiteURL()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(reply)
		time.Sleep(1 * time.Minute)
	}
	os.Exit(0)
}
