package operator

import (
	"fmt"
	expect "github.com/google/goexpect"
	"github.com/redhat-nfvpe/test-network-function/internal/reel"
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf"
	log "github.com/sirupsen/logrus"
	"regexp"
	"time"
)

const (

	checkOperators ="oc get pod %s -n | grep %s"
	checkOperatorStatus ="oc get pod %s -n %s -o json | jq -r '.status'"
)

//Csv Cluster service version , manifests of the operator
type Pod struct {
	result       int
	timeout      time.Duration
	args         []string
	Name         string
	Namespace    string
	Status       string
	ExpectedStatus string
}

// Args returns the command line args for the test.
func (p *Pod) Args() []string {
	return p.args
}

// Timeout return the timeout for the test.
func (d *Pod) Timeout() time.Duration {
	return d.timeout
}

// Result returns the test result.
func (p *Pod) Result() int {
	return p.result
}

// ReelFirst returns a step which expects an csv status for the given csv.
func (p *Pod) ReelFirst() *reel.Step {
	log.Info("At Reel first")
	return &reel.Step{
		Expect:  []string{p.ExpectedReplicaStatus},
		Timeout: p.timeout,
	}
}

// ReelMatch parses the csv status output and set the test result on match.
// Returns no step; the test is complete.
func (p *Pod) ReelMatch(_ string, _ string, match string) *reel.Step {
	re := regexp.MustCompile(p.ExpectStatus)
	p.result = tnf.ERROR
	matched := re.MatchString(match)
	if matched {
		p.result = tnf.SUCCESS
	}
	return nil
}

// ReelTimeout does nothing;
func (p  *Pod) ReelTimeout() *reel.Step {
	return nil
}

// ReelEOF does nothing;
func (p *Pod) ReelEOF() {
}

// NewCsv creates a new `NewCsv` test which runs the "csv" status.
func NewPod(name, namespace string, timeout time.Duration) *Csv {
	args := fmt.Sprintf(CheckCSVCommand, name, namespace)
	return &Csv{
		Name:      name,
		Namespace: namespace,
		result:    tnf.ERROR,
		timeout:   timeout,
		args:      []string{args},
	}
}
