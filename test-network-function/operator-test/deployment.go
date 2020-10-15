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
     checkDeployment ="oc get deployment %s -n %s -o json | jq -r '.status.conditions | sort_by(.lastUpdateTime) | .[-1].type'"
     //checkDeploymentReplicaStatus ="oc get deployment %s -n %s -o json | jq -r '.status.replicas'"
	//checkDeploymentReplicaStatus ="oc get deployment %s -n %s -o json | jq -r '.status.readyReplicas'"
     //outputs true or false
    checkDeploymentIsReady ="oc get deployment %s -n %s -o json | jq -r '.status.replicas==.status.readyReplicas'"
	CheckForDeploymentContainers ="oc get deployment %s -n %s -o json | jq -r '.spec.install.spec.deployments[].spec.template.spec.containers[].name'"

	// oc  get csv etcdoperator.v0.9.4  -n my-etcd  -o json  | jq -r '{name: .metadata.name ,namespace: .metadata.namespace, status:.status.phase ,deployments:{name: .spec.install.spec.deployments[].name, replicas: .spec.install.spec.deployments[].spec.replicas},crds: { owned:[{name:.spec.customresourcedefinitions.owned[].name,kind:.spec.customresourcedefinitions.owned[].kind}]  }}'

)

//Csv Cluster service version , manifests of the operator
type Deployment struct {
	result       int
	timeout      time.Duration
	args         []string
	Name         string
	Namespace    string
	ReplicaStatus  int
	ReadyReplicaStatus  int
	ReadyStatus bool
	ExpectedReplicaStatus string
	OperatorPod Pod
	deployedPods []Pod
}

// Args returns the command line args for the test.
func (d *Deployment) Args() []string {
	return d.args
}

// Timeout return the timeout for the test.
func (d *Deployment) Timeout() time.Duration {
	return d.timeout
}

// Result returns the test result.
func (d *Deployment) Result() int {
	return d.result
}

// ReelFirst returns a step which expects an csv status for the given csv.
func (d *Deployment) ReelFirst() *reel.Step {
	log.Info("At Reel first")
	return &reel.Step{
		Expect:  []string{d.ExpectedReplicaStatus},
		Timeout: d.timeout,
	}
}

// ReelMatch parses the csv status output and set the test result on match.
// Returns no step; the test is complete.
func (d *Deployment) ReelMatch(_ string, _ string, match string) *reel.Step {
	re := regexp.MustCompile(d.ExpectStatus)
	d.result = tnf.ERROR
	matched := re.MatchString(match)
	if matched {
		d.result = tnf.SUCCESS
	}
	return nil
}

// ReelTimeout does nothing;
func (d *Deployment) ReelTimeout() *reel.Step {
	return nil
}

// ReelEOF does nothing;
func (d *Deployment) ReelEOF() {
}

// NewCsv creates a new `NewCsv` test which runs the "csv" status.
func NewDeployment(name, namespace string, timeout time.Duration) *Csv {
	args := fmt.Sprintf(CheckCSVCommand, name, namespace)
	return &Csv{
		Name:      name,
		Namespace: namespace,
		result:    tnf.ERROR,
		timeout:   timeout,
		args:      []string{args},
	}
}
