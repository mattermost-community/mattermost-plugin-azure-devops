import React from 'react';

import './styles.scss';

type LabelValuePairProps = {
    label: string;
    value: string
}

const LabelValuePair = ({label, value}: LabelValuePairProps) => (
    <p className='margin-bottom-10'>
        <strong>{`${label}: `}</strong>
        <span className='value'>{value}</span>
    </p>
);

export default LabelValuePair;
