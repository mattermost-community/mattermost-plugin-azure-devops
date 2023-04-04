import React from 'react';

import Tooltip from 'components/tooltip';

import {onPressingEnterKey} from 'utils';

import './styles.scss';

type LabelValuePairProps = {
    onClickValue?: () => void
    label?: string | JSX.Element;
    value: string
    icon?: LabelIconProps
    labelExtraClassName?: string
}

const LabelValuePair = ({label, icon, value, labelExtraClassName, onClickValue}: LabelValuePairProps) => (
    <p className='margin-bottom-10 d-flex align-item-center'>
        {
                icon?.className && (
                <Tooltip tooltipContent={icon?.tooltipText}>
                    <i
                        className={`${icon?.className} ${icon?.extraClassName} icon-mm`}
                        aria-hidden='true'
                    />
                </Tooltip>
            )
        }
        {
            label && (
                typeof (label) === 'string' ?
                    <strong className={labelExtraClassName ?? ''}>{`${label}: `}</strong> :
                    <span className={`icon label-icon ${labelExtraClassName}`}>{label}</span>
            )
        }
        <Tooltip tooltipContent={value}>
            {
                onClickValue ? (
                    <span
                        aria-label={value}
                        role='button'
                        tabIndex={0}
                        className='value font-size-14 font-bold link-title margin-0 text-truncate'
                        onKeyDown={(event) => onPressingEnterKey(event, onClickValue)}
                        onClick={onClickValue}
                    >
                        {value}
                    </span>
                ) : <span className='value text-truncate'>{value}</span>
            }
        </Tooltip>
    </p>
);

export default LabelValuePair;
