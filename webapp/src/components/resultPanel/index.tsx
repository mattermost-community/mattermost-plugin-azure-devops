import React, {forwardRef} from 'react';

import SVGWrapper from 'components/svgWrapper';
import plugin_constants from 'plugin_constants';

type ResultPanelProps = {
    iconClass?: string | null;
    header?: string | null;
    className?: string;
    primaryBtnText?: string;
    secondaryBtnText?: string;
    onPrimaryBtnClick?: () => void;
    onSecondaryBtnClick?: () => void;
};

const ResultPanel = forwardRef<HTMLDivElement, ResultPanelProps>(({
    header,
    className = '',
    primaryBtnText,
    secondaryBtnText,
    onPrimaryBtnClick,
    onSecondaryBtnClick,
    iconClass,
}: ResultPanelProps, feedAddedRef): JSX.Element => (
    <div
        className={`modal__body modal-body d-flex align-items-center justify-content-center flex-column secondary-panel ${className}`}
        ref={feedAddedRef}
    >
        <>
            {iconClass ? (
                <i className={`fa result-panel-icon ${iconClass || 'fa-check-circle-o'}`}/>
            ) : (
                <SVGWrapper
                    className='result-panel-icon'
                    width={58}
                    height={58}
                    viewBox='0 0 58 58'
                >
                    {plugin_constants.SVGIcons.check}
                </SVGWrapper>
            )}
            <h3 className='result-panel-text'>{header || 'Add new'}</h3>
            {onPrimaryBtnClick && (
                <button
                    className='btn btn-primary'
                    onClick={onPrimaryBtnClick}
                >
                    {primaryBtnText || 'Create new'}
                </button>
            )}
            {onSecondaryBtnClick && (
                <button
                    className='btn btn-link result-panel-close-btn'
                    onClick={onSecondaryBtnClick}
                >
                    {secondaryBtnText || 'Close'}
                </button>
            )}
        </>
    </div>
));

export default ResultPanel;
