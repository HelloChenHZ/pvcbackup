package main

import (
	handler "HelloChenHZ/pvcbackup/handler"
	"flag"
	"fmt"
)

const backupImageName = "backupimage:v1.0"
const recoveryImageName = "recoveryimage:v1.0"

func main() {
	handler.Init()

	action := flag.String("a", "", "action")
	pvcName := flag.String("p", "", "PVC Name")
	dataPath := flag.String("d", "", "Data Path")
	s3Path := flag.String("d", "", "s3 bucket Path")
	flag.Parse()
	fmt.Println(*action, *pvcName, *s3Path)
	// get node path by pvc name
	nodeName := handler.GetNodeName(*pvcName)

	if *action == "backup" {
		// create job
		handler.CreateJob(*pvcName, nodeName, *dataPath, *s3Path, backupImageName, "-d "+*dataPath+"-s "+*s3Path)
	}

	if *action == "recovery" {
		handler.CreateJob(*pvcName, nodeName, *dataPath, *s3Path, recoveryImageName, "-d "+*dataPath+"-s "+*s3Path)
	}
}
