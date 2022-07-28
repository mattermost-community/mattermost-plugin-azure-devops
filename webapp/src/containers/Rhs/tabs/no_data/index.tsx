import React from 'react';

import './styles.scss';

type NoDataPropTypes = {
    title: string,
    subTitle: string,
    buttonText: string,
    buttonAction: (event: React.SyntheticEvent) => void;
}

const NoData = ({title, subTitle, buttonText, buttonAction}: NoDataPropTypes) => {
    return (
        <div className='no-data d-flex align-items-center'>
            <div className='d-flex flex-column align-items-center'>
                <div className='no-data__icon d-flex justify-content-center align-items-center'>
                    <svg
                        width='32px'
                        height='32px'
                        viewBox='0 0 32 32'
                        xmlns='http://www.w3.org/2000/svg'
                    >
                        <path d='M0 11.865l2.995-3.953 11.208-4.557v-3.292l9.828 7.188-20.078 3.896v10.969l-3.953-1.141zM32 5.932v19.536l-7.672 6.531-12.401-4.073v4.073l-7.974-9.885 20.078 2.396v-17.26z'/>
                    </svg>
                </div>
                <p className='no-data__title'>{title}</p>
                {subTitle && <p className='no-data__subtitle'>{subTitle}</p>}
                <button
                    onClick={buttonAction}
                    className='no-data__btn btn btn-primary'
                >
                    {buttonText}
                </button>
            </div>
        </div>
    );
};

export default NoData;
