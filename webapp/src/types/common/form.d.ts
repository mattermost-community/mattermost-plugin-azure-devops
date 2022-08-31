type FieldType = 'dropdown' | 'text' | 'timestamp'
type ValidationTypes = 'isRequired' | 'maxCharLen' | 'minCharLen' | 'regex' | 'regexErrorMessage'
type ErrorComponents = 'LinkProjectModal' | 'CreateTaskModal' | 'SubscribeModal'

type LinkProjectModalFields = 'organization' | 'project' | 'timestamp'
type CreateTaskModalFields = 'organization' | 'project' | 'type' | 'title' | 'description' | 'areaPath' | 'timestamp'
type SubscriptionModalFields = 'organization' | 'project' | 'eventType' | 'channelID' | 'timestamp'
type ErrorComponents = 'SubscribeModal' | 'CreateTaskModal'

type ModalFormFieldConfig = {
    label: string
    value: string
    type: FieldType,
    optionsList?: LabelValuePair[]
    validations?: Partial<Record<ValidationTypes, string | number | boolean>>
}

type FormFields = LinkProjectModalFields | CreateTaskModalFields | SubscriptionModalFields
