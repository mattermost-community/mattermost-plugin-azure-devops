type FieldType = 'dropdown' | 'text' | 'timestamp'
type ValidationTypes = 'isRequired' | 'maxCharLen' | 'minCharLen' | 'regex' | 'regexErrorMessage'
type SubscriptionModalFields = 'organization' | 'project' | 'eventType' | 'channelID' | 'timestamp'
type ErrorComponents = 'SubscribeModal'

type ModalFormFieldConfig = {
    label: string
    value: string
    type: FieldType,
    optionsList?: LabelValuePair[]
    validations?: Partial<Record<ValidationTypes, string | number | boolean>>
}

type FormFieldNames = SubscriptionModalFields
