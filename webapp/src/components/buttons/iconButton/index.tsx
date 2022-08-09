import React from 'react';
import {Button} from 'react-bootstrap';

import Tooltip from 'components/tooltip';

import './styles.scss';

type IconColor = 'danger'

type IconButtonProps = {
    iconClassName: string
    tooltipText: string
    iconColor?: IconColor
    extraClass?: string
    onClick?: () => void
}

const IconButton = ({tooltipText, iconClassName, extraClass = '', iconColor, onClick}: IconButtonProps) => {
    return (
        <Tooltip tooltipContent={tooltipText}>
            <Button
                variant='outline-danger'
                className={`button-wrapper  ${extraClass} ${iconColor === 'danger' && 'danger'}`}
            >
                <i
                    className={iconClassName}
                    aria-hidden='true'
                />
            </Button>
        </Tooltip>
    );
};

export default IconButton;
