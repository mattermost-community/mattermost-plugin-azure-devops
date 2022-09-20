import React from 'react';

import SVGWrapper from 'components/svgWrapper';

import plugin_constants from 'plugin_constants';

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
    isLoading?: boolean;
}

// TODO: UI to be changed
const EmptyState = ({title, subTitle, buttonText, buttonAction, icon = 'folder', wrapperExtraClass, isLoading}: EmptyStatePropTypes) => {
    if (isLoading) {
        return null;
    }

    return (
        <div className={`no-data d-flex ${wrapperExtraClass}`}>
            <div className='d-flex flex-column align-items-center'>
                <div className='no-data__icon d-flex justify-content-center align-items-center'>
                    {
                        icon === 'azure' && (
                            <SVGWrapper
                                width={36}
                                height={36}
                                viewBox=' 0 0 36 36'
                            >
                                {plugin_constants.SVGIcons.azure}
                            </SVGWrapper>
                        )
                    }
                    {
                        icon === 'folder' && (
                            <SVGWrapper
                                width={48}
                                height={40}
                                viewBox=' 0 0 48 40'
                            >
                                {plugin_constants.SVGIcons.folder}
                            </SVGWrapper>
                        )
                    }
                    {icon === 'subscriptions' && (
                        <SVGWrapper
                            width={120}
                            height={120}
                            viewBox=' 0 0 120 120'
                        >
                            {plugin_constants.SVGIcons.subscriptions}
                        </SVGWrapper>
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
                {buttonText && buttonAction && (
                    <button
                        onClick={buttonAction}
                        className='plugin-btn no-data__btn btn btn-primary'
                    >
                        {buttonText}
                    </button>
                )}
            </div>
        </div>
    );
};

export default EmptyState;
