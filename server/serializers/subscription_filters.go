package serializers

type PipelineApproveRequest struct {
	Status     string `json:"status"`
	Comments   string `json:"comments"`
	ApprovalID string `json:"approvalId"`
}
