package main

import (
	"flag"
	"log"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	StartUserSession()

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}
}
