import React from 'react';

import './styles.scss';

type BaseCardProps = {
    children: JSX.Element
}

const BaseCard = ({children}: BaseCardProps) => {
    return (
        <div className='wrapper'>{children}</div>
    );
};

export default BaseCard;
