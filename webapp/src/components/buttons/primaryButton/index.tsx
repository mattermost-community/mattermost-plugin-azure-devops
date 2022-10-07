import React from 'react';

type PrimaryButtonProps = {
    text: string
    iconName?: string
    extraClass?: string
    onClick?: () => void
    isSecondaryButton?: boolean
    isDisabled?: boolean
}

const PrimaryButton = ({text, iconName, extraClass = '', onClick, isSecondaryButton, isDisabled = false}: PrimaryButtonProps) => (
    <button
        disabled={isDisabled}
        onClick={onClick}
        className={`plugin-btn btn ${isSecondaryButton ? 'btn-link' : 'btn-primary'} ${extraClass}`}
    >
        {
            iconName && (
                <i
                    className={`${iconName} margin-left-8`}
                    aria-hidden='true'
                />
            )
        }
        {text}
    </button>
);

export default PrimaryButton;
