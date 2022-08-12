import React from 'react';

import './styles.scss';

type DisplayIcon = 'folder' | 'azure'

type EmptyStatePropTypes = {
    title: string,
    subTitle?: {
        text: string
        slashCommand?: string
    },
    buttonText?: string,
    buttonAction?: (event: React.SyntheticEvent) => void;
    icon?: DisplayIcon;
}

// TODO: UI to be changed
const EmptyState = ({title, subTitle, buttonText, buttonAction, icon = 'folder'}: EmptyStatePropTypes) => {
    return (
        <div className='no-data d-flex'>
            <div className='d-flex flex-column align-items-center'>
                <div className='no-data__icon d-flex justify-content-center align-items-center'>
                    {
                        icon === 'azure' && (
                            <svg
                                width='36'
                                height='36'
                                viewBox='0 0 36 36'
                                fill='none'
                                xmlns='http://www.w3.org/2000/svg'
                            >
                                <path
                                    d='M4.449 12.4965L27.033 8.1045L15.978 0V3.7155L3.3705 8.8485L0 13.3065V23.571L4.449 24.855V12.4965ZM13.416 31.407L27.3705 36L36 28.638V6.618L27.0345 8.1045V27.5565L4.449 24.855L13.416 36V31.407Z'
                                    fill='#8E8E8E'
                                />
                            </svg>
                        )
                    }
                    {
                        icon === 'folder' && (
                            <svg
                                width='48'
                                height='40'
                                viewBox='0 0 48 40'
                                fill='none'
                                xmlns='http://www.w3.org/2000/svg'
                            >
                                <path
                                    d='M2 21.1112V35.2223C2 36.8792 3.34314 38.2223 5 38.2223H36.3333C37.9902 38.2223 39.3333 36.8792 39.3333 35.2223V21.1112M2 21.1112V11.6667C2 10.0099 3.34315 8.66675 5 8.66675H13.5361C14.1284 8.66675 14.7074 8.84206 15.2002 9.1706L20.6887 12.8296C21.1815 13.1581 21.7605 13.3334 22.3528 13.3334H36.3333C37.9902 13.3334 39.3333 14.6766 39.3333 16.3334V21.1112M2 21.1112H39.3333'
                                    stroke='#8E8E8E'
                                    strokeWidth='3.5'
                                    strokeLinecap='round'
                                    strokeLinejoin='round'
                                />
                                <path
                                    d='M9.77783 2H16.9917C18.3737 2 19.7248 2.40907 20.8746 3.17565L24.3477 5.49102C25.4976 6.2576 26.8486 6.66667 28.2306 6.66667H39.0001C42.8661 6.66667 46.0001 9.80067 46.0001 13.6667V28.4444'
                                    stroke='#8E8E8E'
                                    strokeWidth='3.5'
                                    strokeLinecap='round'
                                    strokeLinejoin='round'
                                />
                            </svg>
                        )
                    }
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
};

export default EmptyState;
