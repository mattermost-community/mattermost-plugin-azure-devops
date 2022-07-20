import React from 'react';
import {Modal as RBModal} from 'react-bootstrap';

import './styles.scss';

type ModalFooterProps = {
    onConfirm?: () => void;
    confirmBtnText?: string;
    cancelBtnText?: string;
    onHide?: () => void;
    className?: string;
    confirmDisabled?: boolean;
    cancelDisabled?: boolean;
}

const ModalFooter = ({onConfirm, onHide, cancelBtnText, confirmBtnText, className = '', confirmDisabled, cancelDisabled}: ModalFooterProps) : JSX.Element => (
    <RBModal.Footer className={`modal__footer d-flex flex-column justify-content-center align-items-center ${className}`}>
        {onConfirm && (
            <button
                className='btn btn-primary modal__confirm-btn'
                onClick={onConfirm}
                disabled={confirmDisabled}
            >
                {confirmBtnText || 'Confirm'}
            </button>
        )}
        {onHide && (
            <button
                className='btn btn-link modal__cancel-btn'
                onClick={onHide}
                disabled={cancelDisabled}
            >
                {cancelBtnText || 'Cancel'}
            </button>
        )}
    </RBModal.Footer>
);

export default ModalFooter;
