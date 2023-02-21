import React from 'react';
import {Modal as RBModal} from 'react-bootstrap';

import './styles.scss';

type ModalHeaderProps = {
    title?: string | JSX.Element;
    onHide: () => void;
    showCloseIconInHeader?: boolean;
}

const ModalHeader = ({title, showCloseIconInHeader, onHide}: ModalHeaderProps) : JSX.Element => (
    <>{(title || showCloseIconInHeader) && (
        <RBModal.Header className='azd-modal__header azd-modal__header-icon'>
            <div className='azd-modal__title d-flex align-items-center justify-content-between'>
                {title && <p className='azd-modal__title'>{title}</p>}
                {showCloseIconInHeader && (
                    <button
                        className='style--none'
                        onClick={onHide}
                    >
                        <i className='icon icon-close azd-modal__close-icon'/>
                    </button>
                )}
            </div>
        </RBModal.Header>
    )}
    </>
);

export default ModalHeader;
