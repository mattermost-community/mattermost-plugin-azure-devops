import Tooltip from 'components/tooltip';
import React from 'react';

import {onPressingEnterKey} from 'utils';

import './styles.scss';

type LabelValuePairProps = {
    onClickValue?: () => void
    label?: string | JSX.Element;
    value: string
    labelIconClassName?: string
    labelExtraClassName?: string
}

const LabelValuePair = ({label, labelIconClassName, labelExtraClassName, value, onClickValue}: LabelValuePairProps) => (
    <p className='margin-bottom-10 d-flex align-item-center'>
        {
            labelIconClassName && (
                <i
                    className={`${labelIconClassName} ${labelExtraClassName} icon-mm`}
                />
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
