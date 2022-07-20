import React from 'react';

import './styles.scss';

type ModalSubTitleAndErrorProps = {
    subTitle?: string | JSX.Element;
    error?: string | JSX.Element;
}

const ModalSubTitleAndError = ({subTitle, error}: ModalSubTitleAndErrorProps) : JSX.Element => (
    <>
        {subTitle && <p className='modal__subtitle'>{subTitle}</p>}
        {error && <p className='modal__error'>{error}</p>}
    </>
);

export default ModalSubTitleAndError;
