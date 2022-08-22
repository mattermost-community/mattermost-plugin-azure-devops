import React from 'react';

import Dropdown from 'components/dropdown';

type Props = {
    fieldConfig: Pick<ModalFormFieldConfig, 'label' | 'type' | 'validations'>
    value: any
    optionsList: any
    onChange: (newValue: string) => void;
    error?: string
    isDisabled?: boolean
}

/**
 * A generic component to render form
 * you can add multiple input field types here
 */
const Form = ({fieldConfig: {label, type, validations}, value, optionsList, onChange, error, isDisabled}: Props): JSX.Element => {
    switch (type) {
    case 'dropdown' :
        return (
            <Dropdown
                placeholder={label}
                value={value}
                onChange={onChange}
                options={optionsList || []}
                required={validations?.isRequired as boolean}
                error={error}
                disabled={isDisabled}
            />
        );
    default:
        return <></>;
    }
};

export default Form;
