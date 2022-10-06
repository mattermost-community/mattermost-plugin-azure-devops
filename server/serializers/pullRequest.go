package serializers

type PullRequestValue struct {
	PullRequestID int        `json:"pullRequestId"`
	Title         string     `json:"title"`
	SourceRefName string     `json:"sourceRefName"`
	TargetRefName string     `json:"targetRefName"`
	Reviewers     []Reviewer `json:"reviewers"`
}

type Reviewer struct {
	DisplayName string `json:"displayName"`
}
