package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf"
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf/handlers/generic"
	condition2 "github.com/redhat-nfvpe/test-network-function/pkg/tnf/handlers/generic/assertion"
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf/handlers/generic/condition"
	int2 "github.com/redhat-nfvpe/test-network-function/pkg/tnf/handlers/generic/condition/intcondition"
	string2 "github.com/redhat-nfvpe/test-network-function/pkg/tnf/handlers/generic/condition/stringcondition"
	"gopkg.in/yaml.v2"
	"strings"
	"time"
)

var testTimeout = time.Second * 20

var output *string

func parseArgs() {
	output = flag.String("output", "json", "\"json\" (default) or \"yaml\"")
	flag.Parse()
}

func main() {
	parseArgs()

	// set up the command
	//strCommand := "cd /tmp && curl -X GET https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm > epel-release-latest-7.noarch.rpm && yum -y install ./epel-release-latest-7.noarch.rpm 2>&1"
	strCommand := "yum install -y iperf 2>&1"
	description := "This test attempts to install iperf in the underlying container."
	command := strings.Split(strCommand, " ")

	// set up matches and assertions
	installMatch := `(?m).*Complete!`
	alreadyInstalledMatch := `(?m).*Nothing to do`
	notAvailable := `(?m)No package (\w+) available.*`
	abl := condition2.NewOrBooleanLogic()
	var bl condition2.BooleanLogic = abl

	ic := int2.NewComparisonCondition(27, int2.Equal)
	var icc condition.Condition = ic

	ec := string2.NewEqualsCondition("iperf")
	var c condition.Condition = ec
	ec2 := string2.NewEqualsCondition("unexpected")
	var falseC condition.Condition = ec2
	ec3 := int2.NewIsIntCondition()
	var falseC2 condition.Condition = ec3

	reelMatchResults := []*generic.ResultContext{
		{
			Pattern: notAvailable,
			ComposedAssertions: []condition2.Assertions{
				{
					Assertions: []condition2.Assertion{
						{GroupIdx: 0, Condition: &icc},
						{GroupIdx: 1, Condition: &falseC},
						{GroupIdx: 0, Condition: &falseC2},
						{GroupIdx: 1, Condition: &c},
					},
					Logic: &bl,
				},
			},
			DefaultResult: tnf.ERROR,
		},
		{Pattern: installMatch, DefaultResult: tnf.SUCCESS},
		{Pattern: alreadyInstalledMatch, DefaultResult: tnf.SUCCESS},
	}
	g := generic.GenerateGeneric(command, reelMatchResults, description, testTimeout)
	if *output == "json" {
		j, _ := json.MarshalIndent(g, "", "    ")
		fmt.Println(string(j))
	} else if *output == "yaml" {
		y, _ := yaml.Marshal(g)
		fmt.Println(string(y))
	}
}
