import React, {useState} from 'react';

import './styles.scss';

type DropdownProps = {
    value: string | null;
    placeholder: string;
    onChange: (newValue: string) => void;
    options:LabelValuePair[];
    customOption?: LabelValuePair & {
        onClick: (customOptionValue: string) => void;
    }
    loadingOptions?: boolean;
    disabled?: boolean;
    required?: boolean;
    error?: boolean | string;
}

const Dropdown = ({value, placeholder, options, onChange, customOption, loadingOptions, disabled, error, required}: DropdownProps): JSX.Element => {
    const [open, setOpen] = useState(false);

    // Handles closing the popover and updating the value when someone selects an option
    const handleInputChange = (newOption: LabelValuePair) => {
        setOpen(false);

        // Trigger onChange only if there is a change in the dropdown value
        if (newOption.value !== value) {
            onChange(newOption.value);
        }
    };

    // Handles when someone clicks on the custom option
    const handleCustomOptionClick = () => {
        // Update the value on the input to indicate custom options has been chosen
        handleInputChange({
            label: customOption?.label,
            value: customOption?.value as string,
        });

        // Take the action that need to be taken(only if not already taken) to handle when the user chooses custom option
        if (customOption?.onClick && customOption.value !== value) {
            customOption.onClick(customOption.value);
        }
    };

    const getOptions = () => (customOption ? [...options, {label: customOption.label, value: customOption.value}] : options);

    const getLabel = (optionValue: string | null) => getOptions().find((option) => option.value === optionValue);

    const handleInputBlur = () => {
        // Delaying the closing of the option menu so that when someone chooses any option, the function to update the value is getting called and the updates are happening
        setTimeout(() => {
            setOpen(false);
        }, 250);
    };

    return (
        <div className={`dropdown ${error && 'dropdown--error'}`}>
            <div
                className={`dropdown__field d-flex align-items-center justify-content-between ${open && 'dropdown__field--open'} ${disabled && 'dropdown__field--disabled'}`}
            >
                {placeholder && <label className={`dropdown__field-text dropdown__field-placeholder ${value && 'dropdown__field-placeholder--shifted'}`}>
                    {placeholder}
                    {required && '*'}
                </label>}
                {value && <p className='dropdown__field-text text-ellipses'>
                    {getLabel(value)?.label || getLabel(value)?.value}
                </p>}
                {!loadingOptions && <i className={`fa fa-angle-down dropdown__field-angle ${open && 'dropdown__field-angle--rotated'}`}/>}
                {loadingOptions && <div className='dropdown__loader'/>}
                <input
                    className='dropdown__field-input'
                    onFocus={() => setOpen(true)}
                    onBlur={handleInputBlur}
                    disabled={disabled}
                />
            </div>
            <ul className={`dropdown__options-list ${open && 'dropdown__options-list--open'}`}>
                {
                    options.map((option) => (
                        <li
                            key={option.value}
                            onClick={() => !disabled && handleInputChange(option)}
                            className='dropdown__option-item text-ellipses'
                        >{option.label || option.value}</li>
                    ))
                }
                {customOption && (
                    <li
                        onClick={() => !disabled && handleCustomOptionClick()}
                        className='dropdown__option-item dropdown__custom-option text-ellipses'
                    >
                        {customOption.label || customOption.value}
                    </li>
                )}
            </ul>
            {typeof error === 'string' && <p className='dropdown__err-text'>{error}</p>}
        </div>
    );
};

export default Dropdown;
