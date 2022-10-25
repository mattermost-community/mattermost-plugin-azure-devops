import React from 'react';

import './styles.scss';

type ChipProps = {
    text: string
    extraClass?: string
}

const Chip = ({text, extraClass = ''}: ChipProps) => <div className={`chip Badge__box text-truncate ${extraClass}`}>{text}</div>;

export default Chip;
