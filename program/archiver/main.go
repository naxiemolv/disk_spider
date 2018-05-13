package main

import (
	"fmt"
	"os"
	"os/signal"

	"io/ioutil"
	"encoding/json"

	"github.com/naxiemolv/disk_spider"
	"path/filepath"
	"log"
	"strings"
)

var DirPath = make([]string, 0)
var Suffixs = make([]string, 0)
var verbose = true

func main() {
	loadConfig("config.json")
	parseCommandLine()

	fc := make(chan *disk_spider.File, 100)
	exit := make(chan os.Signal, 1)

	if len(DirPath) > 0 == false {
		Println("[no target dir]")
		os.Exit(0)
	}

	go func() {
		for _, v := range DirPath {
			_,err:=disk_spider.WalkDirToChan(v, Suffixs, fc)
			if err!=nil {
				Println(err.Error())
			}

		}

	}()

	go func() {
		arch, err := disk_spider.NewArchiver("arc.x")

		if err != nil {
			Println("[can not create archive file]")
			os.Exit(1)
		}

		for {
			f := <-fc

			if f == nil {
				break
			}
			Println(f.FileName)

			err:= arch.Archive(f)
			if err != nil {
				fmt.Println(err.Error())
			}

		}
		arch.Finish()
		Println("[process finish]")
		exit<-os.Interrupt


	}()
	signal.Notify(exit, os.Interrupt, os.Kill)
	<-exit


}

func Println(v ... interface{}) {
	if verbose {
		fmt.Println(v)
	}
}

func loadConfig(jsonPath string) {

	b, err := ioutil.ReadFile(getCurrentDirectory() + "/" + jsonPath)
	if err != nil {
		Println("[config JSON Error]:", err.Error())
		os.Exit(0)
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(b, &m)
	if err != nil {
		Println("[config JSON Error]")
		os.Exit(0)
	}

	if dirPaths, ok := m["dir_path"].([]interface{}); ok {
		Println(dirPaths)
		for _, v := range dirPaths {
			DirPath = append(DirPath, v.(string))
		}
	}

	if suffixs, ok := m["suffix"].([]interface{}); ok {
		for _, v := range suffixs {
			Suffixs = append(Suffixs, v.(string))
		}
	}

}

func parseCommandLine() {

}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
