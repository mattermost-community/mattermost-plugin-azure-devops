import React from 'react';
import {FormControl} from 'react-bootstrap';

import './styles.scss';

type InputFieldProps = {
    type?: 'text' | 'password' | 'email';
    label?: string | JSX.Element;
    placeholder?: string;
    value?: string;
    onChange?: (e: React.ChangeEvent) => void;
    error?: boolean | string;
    disabled?: boolean;
    required?: boolean;
}

const Input = ({type = 'text', label, placeholder = '', value = '', onChange, error, disabled = false, required = false}: InputFieldProps) => {
    return (
        <div className='form-group'>
            {label && <label className='form-group__label'>{label}</label>}
            {required && (
                <span
                    className='error-text'
                    style={{marginLeft: '3px'}}
                >{'*'}</span>)
            }
            <FormControl
                type={type}
                value={value}
                onChange={(e) => onChange?.(e)}
                placeholder={placeholder}
                disabled={disabled}
                className={`form-group__control ${error && 'form-group__control--err'}`}
            />
            {(error && typeof error === 'string') && <p className='form-group__err-text'>{error}</p>}
        </div>
    );
};

export default Input;
