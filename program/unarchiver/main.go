package main

import (
	"fmt"
	"github.com/naxiemolv/disk_spider"
)

var verbose = true

func main() {

	parseCommandLine()
	r,_ := disk_spider.NewUnArchiver("arc.x")
	r.UnArchive()

}

func parseCommandLine() {

}

func Println(v ... interface{}) {
	if verbose {
		fmt.Println(v)
	}
}