import React from 'react';

import './styles.scss';

type DisplayIcon = 'folder' | 'azure' | 'subscriptions'

type EmptyStatePropTypes = {
    title: string,
    subTitle?: {
        text: string
        slashCommand?: string
    },
    buttonText?: string,
    buttonAction?: (event: React.SyntheticEvent) => void;
    icon?: DisplayIcon;
    wrapperExtraClass?: string;
}

// TODO: UI to be changed
const EmptyState = ({title, subTitle, buttonText, buttonAction, icon = 'folder', wrapperExtraClass}: EmptyStatePropTypes) => {
    return (
        <div className={`no-data d-flex ${wrapperExtraClass}`}>
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
                    {icon === 'subscriptions' && (
                        <svg
                            width='120'
                            height='120'
                            viewBox='0 0 120 120'
                            fill='none'
                            xmlns='http://www.w3.org/2000/svg'
                        >
                            <rect
                                width='120'
                                height='120'
                                rx='60'
                                fill='#F0F0F1'
                            />
                            <path
                                d='M77.6914 76.5548L43.1814 42.0273L40.9414 44.2498L49.2014 52.5098C48.2564 54.2773 47.7489 56.2548 47.7489 58.2498V68.7498L44.2489 72.2498V73.9998H70.6914L75.4689 78.7773L77.6914 76.5548ZM51.2489 70.4998V58.2498C51.2489 57.1823 51.4414 56.1323 51.8439 55.1523L67.1914 70.4998H51.2489ZM56.4989 75.7498H63.4989C63.4989 76.6781 63.1302 77.5683 62.4738 78.2247C61.8174 78.8811 60.9272 79.2498 59.9989 79.2498C59.0706 79.2498 58.1804 78.8811 57.524 78.2247C56.8677 77.5683 56.4989 76.6781 56.4989 75.7498ZM53.5064 47.9073C54.4339 47.3123 55.4489 46.8748 56.4989 46.5073C56.4989 46.3323 56.4989 46.1748 56.4989 45.9998C56.4989 45.0716 56.8677 44.1813 57.524 43.525C58.1804 42.8686 59.0706 42.4998 59.9989 42.4998C60.9272 42.4998 61.8174 42.8686 62.4738 43.525C63.1302 44.1813 63.4989 45.0716 63.4989 45.9998C63.4989 46.1748 63.4989 46.3323 63.4989 46.5073C68.6964 48.0473 72.2489 52.8248 72.2489 58.2498V66.6498L68.7489 63.1498V58.2498C68.7489 55.9292 67.827 53.7036 66.1861 52.0627C64.5451 50.4217 62.3196 49.4998 59.9989 49.4998C58.6339 49.4998 57.2864 49.8498 56.0789 50.4798L53.5064 47.9073Z'
                                fill='#8E8E8E'
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
