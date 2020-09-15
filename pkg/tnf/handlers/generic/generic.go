package generic

import (
	"encoding/json"
	"github.com/redhat-nfvpe/test-network-function/internal/reel"
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf"
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf/handlers/generic/assertion"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"regexp"
	"time"
)

const (
	// ComposedAssertionsKey is the JSON key used to represent a payload of an array of assertion.Assertions.
	ComposedAssertionsKey = "composedAssertions"
	// DefaultResultKey is the JSON key used to represent the result payload.
	DefaultResultKey = "defaultResult"
	// NextResultContextsKey is the JSON key used to represent a payload of an array of future ResultContext.
	NextResultContextsKey = "nextResultContexts"
	// NextStepKey is the JSON key used to represent the next step payload.
	NextStepKey = "nextStep"
	// PatternKey is the JSON key used to represent a pattern payload.
	PatternKey = "pattern"
)

// Generic is a construct for defining an arbitrary simple test with prescriptive confines.  Essentially, the definition
// of the state machine for a Generic reel.Handler is restricted in this facade, since most common use cases do not need
// to perform too much heavy lifting that would otherwise require a Custom reel.Handler implementation.  Although
// Generic is exported for serialization reasons, it is recommended to instantiate new instances of Generic using either
// NewGenericFromJSONFile or NewGenericFromYAMLFile, which are tailored to properly initialize a Generic.
type Generic struct {

	// Arguments is the Unix command array.  Arguments is optional;  a command can also be issued using ReelFirstStep.
	Arguments []string `json:"arguments,omitempty" yaml:"arguments,omitempty"`

	// Description is a textual description of the overall functionality that is tested.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// FailureReason optionally stores extra information pertaining to why the test failed.

	FailureReason string `json:"failureReason,omitempty" yaml:"failureReason,omitempty"`

	// Matches contains an in order array of matches.
	Matches []Match `json:"matches,omitempty" yaml:"matches,omitempty"`

	// ReelFirstStep is the first Step returned by reel.ReelFirst().
	ReelFirstStep *reel.Step `json:"reelFirstStep,omitempty" yaml:"reelFirstStep,omitempty"`

	// ReelFirstStep is the first Step returned by reel.ReelFirst().
	ReelMatchStep *reel.Step `json:"reelMatchStep,omitempty" yaml:"reelMatchStep,omitempty"`

	// ResultContexts provides the ability to make assertion.Assertions based on the given pattern matched.
	ResultContexts []*ResultContext `json:"resultContexts,omitempty" yaml:"resultContexts,omitempty"`

	// reelMatchResultMap is an internal construct used to save time on lookups.  Since evaluation order of reel.Step
	// Expect regular expressions is important, the end user should define the order (ResultContexts) and realize that
	// the evaluating each regular expression is O(n).  However, when making lookups after the fact, the match pattern
	// has already been found, so ordering does not matter.  This solution duplicates data, but utilizing extra RAM on
	// the bastion server is not a concern.  Performance is favored over memory frugality.
	reelMatchResultMap map[string]int

	// ReelTimeoutStep is the reel.Step to take upon timeout.
	ReelTimeoutStep *reel.Step `json:"reelTimeoutStep,omitempty" yaml:"reelTimeoutStep,omitempty"`

	// TestResult is the result of running the tnf.Test.  0 indicates SUCCESS, 1 indicates FAILURE, 2 indicates ERROR.
	TestResult int `json:"testResult" yaml:"testResult"`

	// TestTimeout prevents the Test from running forever.
	TestTimeout time.Duration `json:"testTimeout,omitempty" yaml:"testTimeout,omitempty"`

	// currentReelMatchResultContexts is used to persist the current ResultContext over multiple invocations of ReelMatch.
	currentReelMatchResultContexts []*ResultContext
}

// ResultContext evaluates the Result for a given Match.  If ComposedAssertions is not supplied, then Result is assigned
// to the reel.Handler result.  If ComposedAssertions is supplied, then the ComposedAssertions are evaluated against the
// match.  The result of ComposedAssertions evaluation is assigned to the reel.Handler's result.
type ResultContext struct {

	// Pattern is the pattern causing a match in reel.Handler ReelMatch.
	Pattern string `json:"pattern,omitempty" yaml:"pattern,omitempty"`

	// ComposedAssertions is a means of making many assertion.Assertion claims about the match.
	ComposedAssertions []assertion.Assertions `json:"composedAssertions,omitempty" yaml:"composedAssertions,omitempty"`

	// DefaultResult is the result of the test.  This is only used if ComposedAssertions is not provided.
	DefaultResult int `json:"defaultResult,omitempty" yaml:"defaultResult,omitempty"`

	// NextStep is an optional next step to take after an initial ReelMatch.
	NextStep *reel.Step `json:"nextStep,omitempty" yaml:"nextStep,omitempty"`

	// NextResultContexts is an optional array which provides the ability to make assertion.Assertions based on the next pattern match.
	NextResultContexts []*ResultContext `json:"nextResultContexts,omitempty" yaml:"nextResultContexts,omitempty"`
}

// MarshalJSON is a shim provided over the default implementation that omits empty NextResultContexts slices.  This
// custom MarshallJSON implementation is needed due to a recursive definition (type ResultContext has a property of type
// ResultContext).
func (r *ResultContext) MarshalJSON() ([]byte, error) {
	if len(r.NextResultContexts) <= 0 {
		return json.Marshal(&struct {
			Pattern            string                 `json:"pattern,omitempty"`
			ComposedAssertions []assertion.Assertions `json:"composedAssertions,omitempty"`
			DefaultResult      int                    `json:"defaultResult"`
			NextStep           *reel.Step             `json:"nextStep,omitempty"`
		}{
			Pattern:            r.Pattern,
			ComposedAssertions: r.ComposedAssertions,
			DefaultResult:      r.DefaultResult,
			NextStep:           r.NextStep,
		})
	}

	// Normally, you would just augment the struct here by adding the missing NextResultContexts field.  However, since
	// NextResultContexts is recursive (i.e., it is a ResultContext), doing so causes a loop.  Thus, this requires a
	// more robust definition.
	return json.Marshal(&struct {
		Pattern            string                 `json:"pattern,omitempty"`
		ComposedAssertions []assertion.Assertions `json:"composedAssertions,omitempty"`
		DefaultResult      int                    `json:"defaultResult"`
		NextStep           *reel.Step             `json:"nextStep,omitempty"`
		NextResultContexts []*ResultContext       `json:"nextResultContexts,omitempty"`
	}{
		Pattern:            r.Pattern,
		ComposedAssertions: r.ComposedAssertions,
		DefaultResult:      r.DefaultResult,
		NextStep:           r.NextStep,
		NextResultContexts: r.NextResultContexts,
	})
}

// Match follows the Container design pattern, and is used to store the arguments to a reel.Handler's ReelMatch
// function in a single data transfer object.
type Match struct {

	// Pattern is the pattern causing a match in reel.Handler ReelMatch.
	Pattern string `json:"pattern,omitempty" yaml:"pattern,omitempty"`

	// Before contains the text before the Match.
	Before string `json:"before,omitempty" yaml:"before,omitempty"`

	// Match is the matched string.
	Match string `json:"match,omitempty" yaml:"match,omitempty"`
}

// init initializes a Generic, including building up the reelMatchResultMap.  reelMatchResultMap is pre-built for
// performance reasons.
func (g *Generic) init() {
	g.currentReelMatchResultContexts = g.ResultContexts
	g.reelMatchResultMap = map[string]int{}
	for _, resultContext := range g.currentReelMatchResultContexts {
		g.reelMatchResultMap[resultContext.Pattern] = resultContext.DefaultResult
	}
}

// Args returns the command line arguments as an array of type string.
func (g *Generic) Args() []string {
	return g.Arguments
}

// Timeout returns the test timeout.
func (g *Generic) Timeout() time.Duration {
	return g.TestTimeout
}

// Result returns the test result.
func (g *Generic) Result() int {
	return g.TestResult
}

// ReelFirst returns the first step to perform.
func (g *Generic) ReelFirst() *reel.Step {
	return g.ReelFirstStep
}

// findResultContext is an internal helper function used to search an array of ResultContext instances for a given
// pattern.  Since order of ResultContext is important, this operation is O(n).
func (g *Generic) findResultContext(pattern string) *ResultContext {
	for _, context := range g.currentReelMatchResultContexts {
		if context.Pattern == pattern {
			return context
		}
	}
	return nil
}

// ReelMatch informs of a match event, returning the next step to perform.
func (g *Generic) ReelMatch(pattern string, before string, match string) *reel.Step {
	m := &Match{Pattern: pattern, Before: before, Match: match}
	g.Matches = append(g.Matches, *m)

	resultContext := g.findResultContext(pattern)
	composedAssertions := resultContext.ComposedAssertions
	if len(composedAssertions) > 0 {
		for _, composedAssertion := range composedAssertions {
			regex := regexp.MustCompile(pattern)
			success, err := (*composedAssertion.Logic).Evaluate(composedAssertion.Assertions, match, *regex)
			if err != nil {
				// exit early on a test error.
				g.FailureReason = err.Error()
				g.TestResult = tnf.ERROR
				return nil
			} else if success == false {
				// exit early on failure
				g.TestResult = tnf.FAILURE
				return nil
			}
			// only report success if nothing else is left
			if resultContext.NextStep == nil {
				g.TestResult = tnf.SUCCESS
				return nil
			}
		}
	}

	// Else, see if we have more work to do.  If not, return defaultResult.
	if resultContext.NextStep == nil {
		g.TestResult = resultContext.DefaultResult
		return nil
	}

	g.currentReelMatchResultContexts = resultContext.NextResultContexts
	return resultContext.NextStep
}

// ReelTimeout informs of a timeout event, returning the next step to perform.
func (g *Generic) ReelTimeout() *reel.Step {
	return g.ReelTimeoutStep
}

// ReelEOF informs of the eof event.
func (g *Generic) ReelEOF() {
	// do nothing.
}

// dataUnmarshaler is a shim interface abstracted over the Unmarshal function, which has an identical signature for json
// and yaml.  This allows for code reuse in createGeneric.
type dataUnmarshaler interface {
	Unmarshal(data []byte, v interface{}) error
}

// yamlUnmarshaler is a YAML implementation of the dataUnmarshaler interface.
type yamlUnmarshaler struct{}

// Unmarshal delegates to yaml.Unmarshal.
func (y *yamlUnmarshaler) Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

// jsonUnmarshaler is a JSON implementation of the dataUnmarshaler interface.
type jsonUnmarshaler struct{}

// Unmarshal delegates to json.Unmarshal.
func (j *jsonUnmarshaler) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// NewGenericFromYAMLFile instantiates and initializes a Generic from a YAML-serialized file.
func NewGenericFromYAMLFile(filename string) (*tnf.Tester, []reel.Handler, error) {
	y := &yamlUnmarshaler{}
	var unmarshaler dataUnmarshaler = y
	g, err := createGeneric(filename, &unmarshaler)
	if err != nil {
		return nil, nil, err
	}
	// poor man's polymorphism
	var tester tnf.Tester = g
	var handler reel.Handler = g
	return &tester, []reel.Handler{handler}, nil
}

// NewGenericFromJSONFile instantiates and initializes a Generic from a JSON-serialized file.
func NewGenericFromJSONFile(filename string) (*tnf.Tester, []reel.Handler, error) {
	j := &jsonUnmarshaler{}
	var unmarshaler dataUnmarshaler = j
	g, err := createGeneric(filename, &unmarshaler)
	if err != nil {
		return nil, nil, err
	}
	// poor man's polymorphism
	var tester tnf.Tester = g
	var handler reel.Handler = g
	return &tester, []reel.Handler{handler}, nil
}

// createGeneric is a helper function which stores common code for instantiating and initializing a Generic using any
// unmarshaller implementation.
func createGeneric(filename string, unmarshaller *dataUnmarshaler) (*Generic, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	g := &Generic{}
	err = (*unmarshaller).Unmarshal(bytes, g)
	if err != nil {
		return nil, err
	}
	g.init()
	return g, nil
}

// GenerateGeneric is a factory method for creating a Generic version of tnf.Test.
func GenerateGeneric(args []string, reelMatchResults []*ResultContext, description string, testTimeout time.Duration) *Generic {
	var keys []string
	for _, matchResult := range reelMatchResults {
		keys = append(keys, matchResult.Pattern)
	}
	return &Generic{
		Arguments:       args,
		Description:     description,
		ReelFirstStep:   &reel.Step{Expect: keys, Timeout: testTimeout},
		ReelMatchStep:   nil,
		ResultContexts:  reelMatchResults,
		ReelTimeoutStep: nil,
		TestResult:      tnf.ERROR,
		TestTimeout:     testTimeout,
	}
}
