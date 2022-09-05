import React from 'react';
import {Button} from 'react-bootstrap';

import Tooltip from 'components/tooltip';

import {onPressingEnterKey} from 'utils';

import './styles.scss';

type IconColor = 'danger'

type IconButtonProps = {
    iconClassName: string
    tooltipText: string
    extraClass?: string
    onClick?: () => void
}

const IconButton = ({tooltipText, iconClassName, extraClass = '', onClick}: IconButtonProps) => (
    <Tooltip tooltipContent={tooltipText}>
        <Button
            variant='outline-danger'
            className={`plugin-btn button-wrapper btn-icon ${extraClass}`}
            onClick={onClick}
            aria-label={tooltipText}
            role='button'
            tabIndex={0}
            onKeyDown={(event) => onPressingEnterKey(event, () => onClick?.())}
        >
            <i
                className={iconClassName}
                aria-hidden='true'
            />
        </Button>
    </Tooltip>
);

export default IconButton;
