import React from 'react';
import {FormControl} from 'react-bootstrap';

import './styles.scss';

type InputFieldProps = {
    type?: 'text' | 'password' | 'email';
    label?: string | JSX.Element;
    placeholder?: string;
    value: string | null;
    onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void;
    error?: boolean | string;
    disabled?: boolean;
    className?: string;
    required?: boolean;
}

const Input = ({type = 'text', label, placeholder = '', value = '', onChange, error, disabled = false, className = '', required}: InputFieldProps) => {
    return (
        <div className={`form-group ${className}`}>
            {label && <label className='form-group__label'>{label}</label>}
            <FormControl
                type={type}
                value={value ?? ''}
                onChange={(e) => onChange?.(e as React.ChangeEvent<HTMLInputElement>)}
                disabled={disabled}
                className={`form-group__control ${error && 'form-group__control--err'}`}

                // Don't remove the below placeholder, otherwise the label will overlap with input value when input isn't in focused state
                placeholder=' '
            />
            {placeholder && <label className='form-group__placeholder'>
                {placeholder}
                {required && '*'}
            </label>}
            {(error && typeof error === 'string') && <p className='form-group__err-text'>{error}</p>}
        </div>
    );
};

export default Input;
