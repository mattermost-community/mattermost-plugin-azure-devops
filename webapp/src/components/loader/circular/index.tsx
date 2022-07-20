import React from 'react';

import './styles.scss';

const CircularLoader = () => {
    return (
        <div className='loader-container d-flex align-items-center justify-content-center'>
            <div className='loader-container__spinner'/>
        </div>
    );
};

export default CircularLoader;
