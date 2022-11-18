type FieldType = 'dropdown' | 'text' | 'timestamp' | 'hidden'
type ValidationTypes = 'isRequired' | 'maxCharLen' | 'minCharLen' | 'regex' | 'regexErrorMessage'
type ErrorComponents = 'LinkProjectModal' | 'CreateTaskModal' | 'SubscribeModal' | 'ConfirmationModal'

type LinkProjectModalFields = 'organization' | 'project' | 'timestamp'
type CreateTaskModalFields = 'organization' | 'project' | 'type' | 'title' | 'description' | 'areaPath' | 'timestamp'
type SubscriptionModalFields = 'organization' | 'project' | 'eventType' | 'channelID' | 'timestamp' | 'serviceType' | 'repository' | 'targetBranch' | 'repositoryName' | 'pullRequestCreatedBy' | 'pullRequestReviewersContains' | 'pullRequestCreatedByName' | 'pullRequestReviewersContainsName' | 'pushedBy' | 'mergeResult' | 'notificationType' | 'pushedByName' | 'mergeResultName' | 'notificationTypeName' | 'areaPath' | 'buildPipeline' | 'buildStatus' | 'releasePipeline' | 'stageName' | 'approvalType' | 'approvalStatus' | 'releaseStatus' | 'buildStatusName' | 'releasePipelineName' | 'stageNameValue' | 'approvalTypeName' | 'approvalStatusName' | 'releaseStatusName'

type ModalFormFieldConfig = {
    label: string
    value: string
    type: FieldType,
    optionsList?: LabelValuePair[]
    validations?: Partial<Record<ValidationTypes, string | number | boolean>>
}

type FormFieldNames = LinkProjectModalFields | CreateTaskModalFields | SubscriptionModalFields
