package main

import (
	handler "HelloChenHZ/pvcbackup/handler"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

const backupImageName = "backupimage:v1.0"
const recoveryImageName = "recoveryimage:v1.0"

func initDB(dbPath string) {
	// 打开或创建 LevelDB 数据库
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		log.Fatalf("无法打开数据库: %v", err)
	}
	defer db.Close()

	// 插入值
	key := []byte("hello")
	value := []byte("world")
	err = db.Put(key, value, nil)
	if err != nil {
		log.Fatalf("无法插入值: %v", err)
	}

	fmt.Println("值插入成功:", string(key), string(value))
}

func insertDB(dbPath string, key []byte, value []byte) {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		log.Fatalf("无法打开数据库: %v", err)
	}
	defer db.Close()

	// 插入值
	err = db.Put(key, value, nil)
	if err != nil {
		log.Fatalf("无法插入值: %v", err)
	}

	fmt.Println("值插入成功:", string(key), string(value))
}

func traverseDB(dbPath string) {
	// 打开 LevelDB 数据库
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		log.Fatalf("无法打开数据库: %v", err)
	}
	defer db.Close()

	// 创建迭代器
	iter := db.NewIterator(nil, nil)
	defer iter.Release()

	// 遍历迭代器
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		fmt.Printf("键：%s，值：%s\n", key, value)
	}
	if err := iter.Error(); err != nil {
		log.Fatalf("迭代器错误: %v", err)
	}
}

func snapshotDB(dbPath, backupPath string) {
	// Open LevelDB database
	db, err := leveldb.OpenFile(dbPath, &opt.Options{})
	if err != nil {
		log.Fatalf("Unable to open database: %v", err)
	}
	defer db.Close()

	// Create a snapshot
	snapshot, err := db.GetSnapshot()
	if err != nil {
		log.Fatalf("Unable to create snapshot: %v", err)
	}
	defer snapshot.Release()

	// Create a file to write backup
	backupFile, err := os.Create(backupPath)
	if err != nil {
		log.Fatalf("Unable to create backup file: %v", err)
	}
	defer backupFile.Close()

	// Iterate over the snapshot and write data to the file
	iter := snapshot.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		_, err := fmt.Fprintf(backupFile, "%s:%s\n", key, value)
		if err != nil {
			log.Fatalf("Error writing to backup file: %v", err)
		}
	}
	iter.Release()

	if err := iter.Error(); err != nil {
		log.Fatalf("Iterator error: %v", err)
	}

	fmt.Println("Backup completed successfully.")
}

func backupDB(dbPath, backupPath string) {
	// Open LevelDB database
	db, err := leveldb.OpenFile(dbPath, &opt.Options{})
	if err != nil {
		log.Fatalf("Unable to open database: %v", err)
	}
	defer db.Close()

	initDB(backupPath)
	// Create a snapshot
	snapshot, err := db.GetSnapshot()
	if err != nil {
		log.Fatalf("Unable to create snapshot: %v", err)
	}
	defer snapshot.Release()

	// Iterate over the snapshot and write data to the file
	iter := snapshot.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		insertDB(backupPath, key, value)
		//_, err := insertDB(backupPath, key, value)
		//if err != nil {
		//	log.Fatalf("Error writing to backup file: %v", err)
		//}
	}
	iter.Release()

	if err := iter.Error(); err != nil {
		log.Fatalf("Iterator error: %v", err)
	}

	fmt.Println("Backup completed successfully.")
}

func main() {
	//initDB("./test")
	//insertDB("./test", []byte("HELLO"), []byte("WORLD"))
	//traverseDB("./test2")
	//snapshotDB("./test", "./backup")
	//backupDB("./test", "./test2")

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
		handler.CreateJob(*pvcName, nodeName, *dataPath, *s3Path, backupImageName)
	}

	if *action == "recovery" {
		handler.CreateJob(*pvcName, nodeName, *dataPath, *s3Path, recoveryImageName)
	}
}
