import React from 'react';

type PrimaryButtonProps = {
    text: string
    iconName?: string
    extraClass?: string
    onClick?: () => void
}

const PrimaryButton = ({text, iconName, extraClass = '', onClick}: PrimaryButtonProps) => (
    <button
        onClick={onClick}
        className={`plugin-btn btn btn-primary ${extraClass}`}
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
