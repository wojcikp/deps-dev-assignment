package dependenciesloader

type VersionKey struct {
	System  string `json:"system"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Node struct {
	VersionKey VersionKey `json:"versionKey"`
	Bundled    bool       `json:"bundled"`
	Relation   string     `json:"relation"`
	Errors     []string   `json:"errors"`
}

type Edge struct {
	FromNode    int    `json:"fromNode"`
	ToNode      int    `json:"toNode"`
	Requirement string `json:"requirement"`
}

type Dependencies struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
	Error string `json:"error"`
}

type ProjectKey struct {
	ID string `json:"id"`
}

type Documentation struct {
	ShortDescription string `json:"shortDescription"`
	URL              string `json:"url"`
}

type Check struct {
	Name          string        `json:"name"`
	Documentation Documentation `json:"documentation"`
	Score         int           `json:"score"`
	Reason        string        `json:"reason"`
	Details       []string      `json:"details"`
}

type Repository struct {
	Name   string `json:"name"`
	Commit string `json:"commit"`
}

type ScorecardInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
}

type Scorecard struct {
	Date         string        `json:"date"`
	Repository   Repository    `json:"repository"`
	Scorecard    ScorecardInfo `json:"scorecard"`
	Checks       []Check       `json:"checks"`
	OverallScore float64       `json:"overallScore"`
	Metadata     []string      `json:"metadata"`
}

type DependencyDetails struct {
	ProjectKey      ProjectKey `json:"projectKey"`
	OpenIssuesCount int        `json:"openIssuesCount"`
	StarsCount      int        `json:"starsCount"`
	ForksCount      int        `json:"forksCount"`
	License         string     `json:"license"`
	Description     string     `json:"description"`
	Homepage        string     `json:"homepage"`
	Scorecard       Scorecard  `json:"scorecard"`
}
