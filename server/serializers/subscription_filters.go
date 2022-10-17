package serializers

type GitRepositoriesResponse struct {
	Value []GitRepositories `json:"value"`
}

type GitRepositories struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GitBranchesResponse struct {
	Value []GitBranches `json:"value"`
}

type GitBranches struct {
	ID   string `json:"objectId"`
	Name string `json:"name"`
}
