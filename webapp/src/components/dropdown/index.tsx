import React, {useState} from 'react';

import './styles.scss';

type DropdownProps = {
    value: string | null;
    placeholder: string;
    onChange: (newValue: string) => void;
    options: {
        label?: string,
        value: string,
    }[];
    customOption?: {
        label?: string;
        value: string;
        onClick: (customOptionValue: string) => void;
    }
}

const Dropdown = ({value, placeholder, options, onChange, customOption}: DropdownProps): JSX.Element => {
    const [open, setOpen] = useState(false);

    // Handles closing the popover and updating the value when someone selects an option
    const handleInputChange = (newValue: { label?: string, value: string }) => {
        setOpen(false);
        onChange(newValue.value);
    };

    // Handles when someone clicks on the custom option
    const handleCustomOptionClick = () => {
        // Update the value on the input to indicate custom options has been chosen
        handleInputChange({
            label: customOption?.label,
            value: customOption?.value as string,
        });

        // Take the action that need to be taken to handle when the user chooses custom option
        if (customOption?.onClick) {
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
        <div className='dropdown'>
            <div
                className={`dropdown__field d-flex align-items-center justify-content-between ${open && 'dropdown__field--open'}`}
            >
                <p className={`dropdown__field-text ${!value && 'dropdown__field-text--placeholder'}`}>{getLabel(value)?.label || getLabel(value)?.value || placeholder}</p>
                <i className={`fa fa-angle-down dropdown__field-angle ${open && 'dropdown__field-angle--rotated'}`}/>
                <input
                    className='dropdown__field-input'
                    onFocus={() => setOpen(true)}
                    onBlur={handleInputBlur}
                />
            </div>
            <ul className={`dropdown__options-list ${open && 'dropdown__options-list--open'}`}>
                {
                    options.map((option) => (
                        <li
                            key={option.value}
                            onClick={() => handleInputChange(option)}
                            className='dropdown__option-item'
                        >{option.label || option.value}</li>
                    ))
                }
                {customOption && (
                    <li
                        onClick={handleCustomOptionClick}
                        className='dropdown__option-item dropdown__custom-option'
                    >
                        {customOption.label || customOption.value}
                    </li>
                )}
            </ul>
        </div>
    );
};

export default Dropdown;
