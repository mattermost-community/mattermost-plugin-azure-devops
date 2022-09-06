type FieldType = 'dropdown' | 'text' | 'timestamp'
type ValidationTypes = 'isRequired' | 'maxCharLen' | 'minCharLen' | 'regex' | 'regexErrorMessage'
type ErrorComponents = 'LinkProjectModal' | 'CreateTaskModal' | 'SubscribeModal'

type SubscriptionModalFields = 'organization' | 'project' | 'eventType' | 'channelID' | 'timestamp'
type CreateTaskModalFields = 'organization' | 'project' | 'type' | 'title' | 'description' | 'areaPath' | 'timestamp';
type LinkProjectModalFields = 'organization' | 'project' | 'timestamp'

type ModalFormFieldConfig = {
    label: string
    value: string
    type: FieldType,
    optionsList?: LabelValuePair[]
    validations?: Partial<Record<ValidationTypes, string | number | boolean>>
}

type FormFieldNames = LinkProjectModalFields | CreateTaskModalFields | SubscriptionModalFields
