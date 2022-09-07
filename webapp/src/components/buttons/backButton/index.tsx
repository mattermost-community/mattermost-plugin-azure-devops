import React from 'react';

import Tooltip from 'components/tooltip';

import {onPressingEnterKey} from 'utils';

type BackButtonProps = {
    onClick: () => void
}

const BackButton = ({onClick}: BackButtonProps) => (
    <Tooltip tooltipContent={'Back'}>
        <svg
            tabIndex={0}
            className='link-title margin-right-8'
            onKeyDown={(event) => onPressingEnterKey(event, onClick)}
            onClick={onClick}
            width='9'
            height='15'
            viewBox='0 0 9 15'
            fill='none'
            xmlns='http://www.w3.org/2000/svg'
        >
            <path
                d='M8 1L1.49999 7.50001L8 14'
                stroke='#1C58D9'
                strokeWidth='2'
            />
        </svg>
    </Tooltip>
);

export default BackButton;
