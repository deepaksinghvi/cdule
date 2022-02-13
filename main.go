package main

import (
	"time"

	"github.com/deepaksinghvi/cdule/pkg/cdule"
)

/*
TODO This is for the development time debugging and to be removed later
*/
func main() {
	cdule := cdule.Cdule{}
	cdule.NewCdule()
	time.Sleep(5 * time.Minute)
	cdule.StopWatcher()
}
