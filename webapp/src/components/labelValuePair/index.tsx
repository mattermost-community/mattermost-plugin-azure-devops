import React from 'react';

import './styles.scss';

type LabelValuePairProps = {
    label?: string | JSX.Element;
    value: string
    labelIconClassName?: string
    labelExtraClassName?: string
}

const LabelValuePair = ({label, labelIconClassName, labelExtraClassName, value}: LabelValuePairProps) => {
    return (
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
                        <span className={`icon ${labelExtraClassName}`}>{label}</span>
                )
            }
            <span className='value text-truncate'>{value}</span>
        </p>
    );
};

export default LabelValuePair;
