import React from 'react';

import './styles.scss';

const Loader = () => {
    return (
        <div className='loader-container d-flex align-items-center justify-content-center'>
            <div className='loader-container__spinner'/>
        </div>
    );
};

export default Loader;
