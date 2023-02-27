import React from 'react';

import LinearLoader from 'components/loader/linear';

import './styles.scss';

type ModalLoaderProps = {
    loading?: boolean;
}

const ModalLoader = ({loading}: ModalLoaderProps): JSX.Element => <div className='azd-modal__loader-container'>{loading && <LinearLoader/>}</div>;

export default ModalLoader;
