package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)


func checkpropfile(){

	proppath := "/etc/wphelper/"

	if _, err := os.Stat(proppath+"wphelperprop"); os.IsNotExist(err){
		err = os.Mkdir(proppath, 777)
		if err != nil {
			log.Fatalln("Unable to create directory in '/etc/':", err)
		}
		_, err = os.Create(proppath+"wphelperpop")
		if err != nil {
			log.Fatal("Error during creation:", err)
		}
		fmt.Println("wphelperprop created")
	}
	fmt.Println("wphelperprop file exists, proceeding...")
}


func dbDataManager(){
		// Open the file if it exists
		file, err := os.Open("/etc/wphelper/wphelperprops")
		if err != nil {
			log.Fatalln("File does not exist", err)
		}
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		var text []string
		for scanner.Scan() {
			text = append(text, scanner.Text())
		}
		file.Close()
		for _, v := range text {
			spv := strings.Split(v, " ")
			dbdata := fmt.Sprintf(`sourceURL: %s
destURL: %s
dbURL: mysqlsvc
dbUsername: root
dbPassword: password
dbName: wordpress`, spv[1], spv[2])

			//check if the directory exists
			_, err := os.Stat(spv[0])
			if err == nil {
				//Write Yaml file using dbdata
				err := ioutil.WriteFile(spv[0]+"dbdata.yml", []byte(dbdata), 0666)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
}

func main() {
	checkpropfile()
	for {
		dbDataManager()
		time.Sleep(10 * time.Second)
	}

}