package serializers

type GitRepositoriesResponse struct {
	Value []*GitRepository `json:"value"`
}

type GitRepository struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GitBranchesResponse struct {
	Value []*GitBranch `json:"value"`
}

type GitBranch struct {
	ID   string `json:"objectId"`
	Name string `json:"name"`
}

type PipelineApproveRequest struct {
	Status   string `json:"status"`
	Comments string `json:"comments"`
}
