# The "generic" `reel.Handler`

Defining test cases is burdensome as each test must implement `reel.Handler` and `tnf.Test` interfaces.  Thus, each test
must define seven different GoLang functions.  Much of the time, the robust `reel.Handler` state machine can be greatly
simplified given the test context.  For example, if a test is just checking for file existence, then defining a State
Machine implementation is cumbersome and overly complex.  The `generic Handler` provides a mechanism for users to more
rapidly create test implementations by simplifying definition of a test.

## Design Aspects Not Considered Yet

* Variable substitution.  Often times, tests require configuration.  For example, if a test pings a particular host,
  perhaps that host is configurable.
  - use golang template, render, validate json schema, include test in catalog
* URI for a test
* Saving information for later use.
* match -> matches
* next should combine nextStep and nextResultContext
error next step, fail next step, success next step?  yes... probably...
* append \n to reel.Step.Execute... not currently working!