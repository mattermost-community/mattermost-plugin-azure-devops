import React from 'react';

import {Input, Select} from '@brightscout/mattermost-ui-library';

type Props = {
    fieldConfig: Pick<ModalFormFieldConfig, 'label' | 'type' | 'validations'>
    value: string | null
    optionsList?: LabelValuePair[]
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
            <Select
                label={label}
                options={optionsList || []}
                onSelectOptionHandler={onChange}
                className='form-input'
            />
        );
    case 'text' :
        return (
            <Input
                label={label}
                type='text'
                value={value ?? ''}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => onChange(e.target.value)}
                error={error}
                required={validations?.isRequired as boolean}
                disabled={isDisabled}
                fullWidth={true}
                className='form-input'
            />
        );
    default:
        return <></>;
    }
};

export default Form;
