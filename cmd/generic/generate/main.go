package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf/handlers/generic"
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf/handlers/generic/assertion"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type CreateTestApplication struct {
	command            string
	reader             *bufio.Reader
	timeout            time.Duration
	reelResultContexts []*generic.ResultContext
}

func (c *CreateTestApplication) init(file *os.File) {
	c.reader = bufio.NewReader(file)
}

func (c *CreateTestApplication) pollForCommand() error {
	fmt.Print("Please enter your command, then press return: ")
	command, err := c.reader.ReadString('\n')
	c.command = command[:len(command)-1]
	return err
}

func (c *CreateTestApplication) pollForTimeoutSeconds() error {
	fmt.Print("Please enter the test timeout in seconds: ")
	timeoutSeconds, err := c.reader.ReadString('\n')
	if err != nil {
		return err
	}
	timeSecondsInt, err := strconv.Atoi(timeoutSeconds[:len(timeoutSeconds)-1])
	if err != nil {
		return err
	}
	c.timeout = time.Duration(timeSecondsInt) * time.Second
	return nil
}

func (c *CreateTestApplication) pollForResultContext() ([]assertion.Assertions, error) {
	var assertionsArray []assertion.Assertions
	fmt.Print("Would you like to add special logic to distinguish if the match results in pass/fail (y/n): ")
	pattern, err := c.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	pattern = pattern[:len(pattern)-1]

	if pattern == "y" || pattern == "Y" {
		fmt.Print("What sort of boolean logic do you want to use across assertionsArray (and/or)? ")
		logicString, err := c.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		logicString = logicString[:len(logicString)-1]
		var booleanLogic assertion.BooleanLogic
		if logicString == "and" {
			andBooleanLogic := assertion.NewAndBooleanLogic()
			booleanLogic = andBooleanLogic
		} else if logicString == "or" {
			orBooleanLogic := assertion.NewOrBooleanLogic()
			booleanLogic = orBooleanLogic
		} else {
			return nil, errors.New(fmt.Sprintf("unknown boolean logic string: %s", logicString))
		}

		fmt.Print("Which match index do you want to use? ")

		assertions := assertion.Assertions{
			Assertions: nil,
			Logic:      &booleanLogic,
		}
		assertionsArray = append(assertionsArray, assertions)

	}
	return assertionsArray, nil
}

func (c *CreateTestApplication) pollForPattern() error {
	fmt.Print("Please enter a pattern regex for your Handler's ReelFirst(): ")
	pattern, err := c.reader.ReadString('\n')
	if err != nil {
		return err
	}
	pattern = pattern[:len(pattern)-1]
	fmt.Print("What is the default test result for this pattern (0=Success, 1=Failure, 2=Error): ")
	defaultResult, err := c.reader.ReadString('\n')
	if err != nil {
		return err
	}
	defaultResult = defaultResult[:len(defaultResult)-1]
	defaultResultInt, err := strconv.Atoi(defaultResult)
	if err != nil {
		return err
	}

	assertionsArray, err := c.pollForResultContext()
	if err != nil {
		return err
	}

	reelResultContext := &generic.ResultContext{
		Pattern:            pattern,
		ComposedAssertions: assertionsArray,
		DefaultResult:      defaultResultInt,
	}
	c.reelResultContexts = append(c.reelResultContexts, reelResultContext)

	return nil
}

func (c *CreateTestApplication) pollForOutputFilename() (string, error) {
	fmt.Print("What is the output test file name: ")
	defaultResult, err := c.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	defaultResult = defaultResult[:len(defaultResult)-1]
	return defaultResult, nil
}

func (c *CreateTestApplication) Run() error {
	err := c.pollForCommand()
	if err != nil {
		return err
	}

	err = c.pollForTimeoutSeconds()
	if err != nil {
		return err
	}

	err = c.pollForPattern()
	return err
}

func (c *CreateTestApplication) GetCommand() string {
	return c.command
}

func (c *CreateTestApplication) GetTimeout() time.Duration {
	return c.timeout
}

func (c *CreateTestApplication) GetResultContexts() []*generic.ResultContext {
	return c.reelResultContexts
}

func NewCreateTestApplication(file *os.File) *CreateTestApplication {
	c := &CreateTestApplication{}
	c.init(file)
	return c
}

func main() {
	app := NewCreateTestApplication(os.Stdin)
	err := app.Run()
	if err != nil {
		panic(err)
	}
	g := generic.GenerateGeneric(strings.Split(app.GetCommand(), " "), app.GetResultContexts(), "generated command", app.GetTimeout())

	outputFile, err := app.pollForOutputFilename()
	if err != nil {
		panic(err)
	}

	j, err := json.MarshalIndent(&g, "", "    ")
	if err != nil {
		panic(err)
	}
	_ = ioutil.WriteFile(outputFile, []byte(string(j)), 0644)
	fmt.Println(string(j))
}
