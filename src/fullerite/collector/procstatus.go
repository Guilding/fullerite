package collector

import (
	"fullerite/config"
	"fullerite/metric"

	"regexp"

	l "github.com/Sirupsen/logrus"
)

// ProcStatus collector type
type ProcStatus struct {
	baseCollector
	compiledRegex    map[string]*regexp.Regexp
	query            *regexp.Regexp
	matchCommandLine bool
}

// Query returns ProcStatus collectors query
func (ps ProcStatus) Query() *regexp.Regexp {
	return ps.query
}

// MatchCommandLine returns ProcStatus collectors matches command line
func (ps ProcStatus) MatchCommandLine() bool {
	return ps.matchCommandLine
}

// NewProcStatus creates a new Test collector.
func NewProcStatus(channel chan metric.Metric, initialInterval int, log *l.Entry) *ProcStatus {
	ps := new(ProcStatus)

	ps.log = log
	ps.channel = channel
	ps.interval = initialInterval

	ps.name = "ProcStatus"
	ps.query = regexp.MustCompile("")
	ps.matchCommandLine = true
	ps.compiledRegex = make(map[string]*regexp.Regexp)

	return ps
}

// Configure this takes a dictionary of values with which the handler can configure itself
func (ps *ProcStatus) Configure(configMap map[string]interface{}) {
	if query, exists := configMap["query"]; exists {
		re, err := regexp.Compile(query.(string))
		if err != nil {
			ps.log.Warn("Failed to compile regex: ", err)
		} else {
			ps.query = re
		}
	}

	if matchCommandLine, exists := configMap["matchCommandLine"]; exists {
		ps.matchCommandLine = matchCommandLine.(bool)
	}

	if generatedDimensions, exists := configMap["generatedDimensions"]; exists {
		for dimension, generator := range config.GetAsMap(generatedDimensions) {
			//don't use MustCompile otherwise program will panic due to misformated regex
			re, err := regexp.Compile(generator)
			if err != nil {
				ps.log.Warn("Failed to compile regex: ", generator, err)
			} else {
				ps.compiledRegex[dimension] = re
			}
		}
	}

	ps.configureCommonParams(configMap)
}
