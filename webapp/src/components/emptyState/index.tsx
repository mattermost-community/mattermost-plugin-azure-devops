import React from 'react';

import './styles.scss';

type EmptyStatePropTypes = {
    title: string,
    subTitle?: {
        text: string
        slashCommand?: string
    },
    buttonText?: string,
    buttonAction?: (event: React.SyntheticEvent) => void;
}

const EmptyState = ({title, subTitle, buttonText, buttonAction}: EmptyStatePropTypes) => (
    <div className='no-data d-flex'>
        <div className='d-flex flex-column align-items-center'>
            <div className='no-data__icon d-flex justify-content-center align-items-center'>
                <svg
                    width='32px'
                    height='32px'
                    viewBox='0 0 32 32'
                >
                    <path d='M0 11.865l2.995-3.953 11.208-4.557v-3.292l9.828 7.188-20.078 3.896v10.969l-3.953-1.141zM32 5.932v19.536l-7.672 6.531-12.401-4.073v4.073l-7.974-9.885 20.078 2.396v-17.26z'/>
                </svg>
            </div>
            <p className='no-data__title'>{title}</p>
            {subTitle && (
                <>
                    <p className='no-data__subtitle'>{subTitle.text}</p>
                    {
                        subTitle.slashCommand && <p className='slash-command'>{subTitle.slashCommand}</p>
                    }

                </>
            )}
            {
                buttonText && buttonAction && (
                    <button
                        onClick={buttonAction}
                        className='plugin-btn no-data__btn btn btn-primary'
                    >
                        {buttonText}
                    </button>
                )
            }
        </div>
    </div>
);

export default EmptyState;
