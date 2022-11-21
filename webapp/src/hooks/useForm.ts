import {useState} from 'react';

import {subscriptionModal} from 'pluginConstants/form';
import subscribeModal from 'reducers/subscribeModal';

// Set initial value of form fields
const getInitialFieldValues = (
    formFields: Partial<Record<FormFieldNames, ModalFormFieldConfig>>,
): Partial<Record<FormFieldNames, string>> => {
    let fields = {};
    Object.keys(formFields).forEach((field) => {
        fields = {
            ...fields,
            [field as FormFieldNames]:
                formFields[field as FormFieldNames]?.value ||
                (field as FormFieldNames === 'timestamp' ? Date.now().toString() : ''),
        };
    });

    return fields as unknown as Partial<Record<FormFieldNames, string>>;
};

/**
 * Filter out all the fields for which validation check is required
 * and set an empty string as the default error message
 */
const getFieldsWhereErrorCheckRequired = (
    formFields: Partial<Record<FormFieldNames, ModalFormFieldConfig>>,
): Partial<Record<FormFieldNames, string>> => {
    let fields = {};
    Object.keys(formFields).forEach((field) => {
        if (formFields[field as FormFieldNames]?.validations) {
            fields = {
                ...fields,
                [field as FormFieldNames]: '',
            };
        }
    });

    return fields as unknown as Partial<Record<FormFieldNames, string>>;
};

// Check each type of validation and return the required error message
const getValidationErrorMessage = (
    formFields: Partial<Record<FormFieldNames, string>>,
    fieldName: FormFieldNames,
    fieldLabel: string,
    validationType: ValidationTypes,
): string => {
    switch (validationType) {
    case 'isRequired':
        return formFields[fieldName] ? '' : `${fieldLabel} is required`;
    default:
        return '';
    }
};

// Generic hook to handle form fields
function useForm(initialFormFields: Partial<Record<FormFieldNames, ModalFormFieldConfig>>) {
    // Form field values
    const [formFields, setFormFields] = useState(getInitialFieldValues(initialFormFields));

    // Form field error state
    const [errorState, setErrorState] = useState<Partial<Record<FormFieldNames, string>>>(
        getFieldsWhereErrorCheckRequired(initialFormFields),
    );

    /**
     * Set new field value on change
     * and reset field error state
     */
    const onChangeFormField = (fieldName: FormFieldNames, value: string) => {
        if (fieldName === 'eventType') {
            setErrorState({...errorState, [fieldName]: ''});
            setFormFields({...getInitialFieldValues(initialFormFields), [fieldName]: value, organization: formFields.organization, project: formFields.project, channelID: formFields.channelID, serviceType: formFields.serviceType});
        } else if (fieldName === 'serviceType') {
            setErrorState({...errorState, [fieldName]: ''});
            setFormFields({...getInitialFieldValues(initialFormFields), [fieldName]: value, organization: formFields.organization, project: formFields.project, channelID: formFields.channelID});
        } else {
            setErrorState({...errorState, [fieldName]: ''});
            setFormFields({...formFields, [fieldName]: value});
        }
    };

    // Validate all form fields and set error if any
    const isErrorInFormValidation = (): boolean => {
        let errorFields = {};
        Object.keys(initialFormFields).forEach((field) => {
            if (initialFormFields[field as FormFieldNames]?.validations) {
                Object.keys(initialFormFields[field as FormFieldNames]?.validations ?? '').forEach((validation) => {
                    const validationMessage = getValidationErrorMessage(
                        formFields,
                        field as FormFieldNames,
                        initialFormFields[field as FormFieldNames]?.label || '',
                        validation as ValidationTypes,
                    );
                    if (validationMessage) {
                        errorFields = {
                            ...errorFields,
                            [field]: validationMessage,
                        };
                    }
                });
            }
        });

        if (!Object.keys(errorFields).length) {
            return false;
        }

        setErrorState(errorFields);
        return true;
    };

    // Reset form field values and error states
    const resetFormFields = () => {
        setFormFields(getInitialFieldValues(initialFormFields));
        setErrorState(getFieldsWhereErrorCheckRequired(initialFormFields));
    };

    // Set value for a specific form field
    const setSpecificFieldValue = (modifiedFormFields: Partial<Record<FormFieldNames, string>>) => {
        setFormFields(modifiedFormFields);
    };

    return {formFields, errorState, setSpecificFieldValue, onChangeFormField, isErrorInFormValidation, resetFormFields};
}

export default useForm;
