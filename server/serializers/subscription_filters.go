package serializers

type PipelineApproveRequest struct {
	Status     string `json:"status"`
	Comments   string `json:"comments"`
	Comment    string `json:"comment"`
	ApprovalID string `json:"approvalId"`
}
