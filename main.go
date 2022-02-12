package main

import (
	"time"

	"github.com/deepaksinghvi/cdule/pkg/cdule"
)

func main() {
	cdule := cdule.Cdule{}
	cdule.NewCdule()
	time.Sleep(5 * time.Minute)
	cdule.StopWatcher()
}
