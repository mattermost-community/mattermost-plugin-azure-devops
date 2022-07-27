import React from 'react';
import {Button} from 'react-bootstrap';

import Tooltip from 'components/tooltip';

import './styles.scss';

type IconButtonProps = {
    iconClassName: string
    tooltipText: string
    extraClass?: string
    onClick?: () => void
}

const IconButton = ({tooltipText, iconClassName, extraClass, onClick}: IconButtonProps) => {
    return (
        <Tooltip tooltipContent={tooltipText}>
            <Button
                variant='outline-danger'
                className={`plugin-btn button-wrapper btn-icon ${extraClass}`}
            >
                <i
                    className={`${iconClassName} icon`}
                    aria-hidden='true'
                />
            </Button>
        </Tooltip>
    );
};

export default IconButton;
