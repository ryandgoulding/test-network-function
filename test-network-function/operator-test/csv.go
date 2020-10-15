package operator

import (
	"fmt"
	"github.com/redhat-nfvpe/test-network-function/internal/reel"
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf"
	log "github.com/sirupsen/logrus"
	"regexp"
	"time"
)

const (
	// CheckCSVCommand is the OC command for checking for CSV.
	CheckCSVCommand = "oc  get csv %s  -n %s -o json  | jq -r '.status.phase'"
	// CheckDeploymentName is the OC command for checking for deployment name from CSV.
	CheckDeploymentName = "oc  get csv %s  -n %s -o json  | jq -r '.spec.install.spec.deployments[].name'"
	CheckDeploymentNameRegEx=".+"
    ChkDepExepectedReplicas="oc  get csv %s  -n %s  -o json  | jq -r '.spec.install.spec.deployments[].spec.replicas'"
	ChkDepReplicasRegEx="\\d"
	//chkDepReplicasRegEx="^(true|fale){1}"
)

//Csv Cluster service version , manifests of the operator
type Csv struct {
	result       int
	timeout      time.Duration
	args         []string
	Name         string
	Status       string
	Namespace    string
	ExpectStatus string
	Deployments []Deployment
}

// Args returns the command line args for the test.
func (c *Csv) Args() []string {
	return c.args
}

// Timeout return the timeout for the test.
func (c *Csv) Timeout() time.Duration {
	return c.timeout
}

// Result returns the test result.
func (c *Csv) Result() int {
	return c.result
}

// ReelFirst returns a step which expects an csv status for the given csv.
func (c *Csv) ReelFirst() *reel.Step {
	log.Info("At Reel first")
	return &reel.Step{
		// this is not the place, the array is used to validate different return
		Expect:  []string{"CSV","DEPLOYMENT","POD"},
		Timeout: c.timeout,
	}
}

// ReelMatch parses the csv status output and set the test result on match.
// Returns no step; the test is complete.
func (c *Csv) ReelMatch(pattern string, _ string, match string) *reel.Step {
	if pattern == "CSV"{
		re := regexp.MustCompile(c.ExpectStatus)
		c.result = tnf.ERROR
		matched := re.MatchString(match)
		if matched {
			c.result = tnf.SUCCESS
		}
	}

	return nil
}

// ReelTimeout does nothing;
func (c *Csv) ReelTimeout() *reel.Step {
	return nil
}

// ReelEOF does nothing;
func (c *Csv) ReelEOF() {
}

// NewCsv creates a new `NewCsv` test which runs the "csv" status.
func NewCsv(name, namespace string, timeout time.Duration) *Csv {
	var cmdArgs []string
	cmdArgs=append(cmdArgs,fmt.Sprintf(CheckCSVCommand, name, namespace))
	cmdArgs=append(cmdArgs,fmt.Sprintf(CheckDeploymentName, name, namespace))
	cmdArgs=append(cmdArgs,fmt.Sprintf(ChkDepReplicasRegEx, name, namespace))

	return &Csv{
		Name:      name,
		Namespace: namespace,
		result:    tnf.ERROR,
		timeout:   timeout,
		args:      cmdArgs,
	}
}
