import React from 'react';

import './styles.scss';

const LinearLoader = ({extraClass = ''}: {extraClass?: string}) : JSX.Element => (
    <div className={`linear-loader ${extraClass}`}>
        <div className='linear-loader__bar'/>
    </div>
);

export default LinearLoader;
