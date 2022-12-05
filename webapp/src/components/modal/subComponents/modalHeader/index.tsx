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
        <RBModal.Header className='modal__header modal__header-icon'>
            <div className='modal__title d-flex align-items-center justify-content-between'>
                {title && <p className='modal__title'>{title}</p>}
                {showCloseIconInHeader && (
                    <button
                        className='style--none'
                        onClick={onHide}
                    >
                        <i className='icon icon-close modal__close-icon'/>
                    </button>
                )}
            </div>
        </RBModal.Header>
    )}
    </>
);

export default ModalHeader;
