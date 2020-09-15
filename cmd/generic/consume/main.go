package main

import (
	"encoding/json"
	"fmt"
	expect "github.com/google/goexpect"
	"github.com/google/goterm/term"
	"github.com/redhat-nfvpe/test-network-function/internal/reel"
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf"
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf/handlers/generic"
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf/interactive"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

var timeout = time.Second * 20

func runTest(expecter *expect.Expecter, tester *tnf.Tester, handlers []reel.Handler, ch <-chan error) {
	t, err := tnf.NewTest(expecter, *tester, handlers, ch)
	if err != nil {
		panic(err)
	}
	result, err := t.Run()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Result: %d\n", result)
	y, err := json.MarshalIndent(tester, "", "    ")
	fmt.Printf("%s\n", y)
}

func setupAndRunTest(file string, expecter *expect.Expecter, errorChannel <-chan error) {
	tester, handlers, err := generic.NewGenericFromJSONFile(file)
	if err != nil {
		log.Errorf("Error: %v", err)
		panic(err)
	}
	runTest(expecter, tester, handlers, errorChannel)
}

func runAgainstRaspberryPi(file string) {
	goExpectSpawner := interactive.NewGoExpectSpawner()
	var spawnContext interactive.Spawner = goExpectSpawner
	context, err := interactive.SpawnSSH(&spawnContext, "pi", "192.168.1.5", timeout, expect.Verbose(true))

	if err != nil {
		panic(err)
	}
	log.Info(term.Bluef("Running %s against my raspberry pi print / sous vide server", file))
	setupAndRunTest(file, context.GetExpecter(), context.GetErrorChannel())
}

func main() {
	file := os.Args[1]
	goExpectSpawner := interactive.NewGoExpectSpawner()
	var spawnContext interactive.Spawner = goExpectSpawner
	oc, ch, err := interactive.SpawnOc(&spawnContext, "test", "test", "default", timeout, expect.Verbose(true))
	if err != nil {
		panic(err)
	}
	log.Info(term.Bluef("Running %s from oc test(test)", file))
	setupAndRunTest(file, oc.GetExpecter(), ch)

	log.Info(term.Bluef("Running %s from development machine", file))
	goExpectSpawner = interactive.NewGoExpectSpawner()
	var spawnContext2 interactive.Spawner = goExpectSpawner
	context, err := interactive.SpawnShell(&spawnContext2, timeout, expect.Verbose(true))
	setupAndRunTest(file, context.GetExpecter(), context.GetErrorChannel())

	if len(os.Args) >= 3 {
		runAgainstRaspberryPi(file)
	}
}
