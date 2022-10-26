import React from 'react';

import Tooltip from 'components/tooltip';

import './styles.scss';

type ChipProps = {
    text: string
    extraClass?: string
}

const Chip = ({text, extraClass = ''}: ChipProps) => (
    <Tooltip tooltipContent={text}>
        <div className={`chip Badge__box text-truncate ${extraClass}`}>{text}</div>
    </Tooltip>
);

export default Chip;
