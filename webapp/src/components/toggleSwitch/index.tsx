import React from 'react';

import './styles.scss';

type ToggleSwitchLabelPositioning = 'left' | 'right'
type ToggleSwitchProps = {
    active: boolean;
    onChange: (active: boolean) => void;
    label?: string;
    labelPositioning?: ToggleSwitchLabelPositioning;
}

const ToggleSwitch = ({
    active,
    onChange,
    label,
    labelPositioning = 'left',
}: ToggleSwitchProps): JSX.Element => (
    <div className={`toggle-switch-container d-flex align-items-center ${labelPositioning === 'right' && 'flex-row-reverse justify-content-end'}`}>
        {label && <span className={labelPositioning === 'left' ? 'toggle-switch-label--left' : 'toggle-switch-label--right'}>{label}</span>}
        <label className='toggle-switch cursor-pointer'>
            <input
                type='checkbox'
                className='toggle-switch__input'
                checked={active}
                onChange={(e) => onChange(e.target.checked)}
            />
            <span className='toggle-switch__slider cursor-pointer'/>
        </label>
    </div>
);

export default ToggleSwitch;
