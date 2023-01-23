import React from 'react';
import {Modal as RBModal} from 'react-bootstrap';

import './styles.scss';

type ModalFooterProps = {
    onConfirm?: (() => void) | null;
    confirmBtnText?: string;
    cancelBtnText?: string;
    onHide?: () => void;
    className?: string;
    confirmDisabled?: boolean;
    cancelDisabled?: boolean;
    confirmAction?: boolean;
}

const ModalFooter = ({onConfirm, onHide, cancelBtnText, confirmBtnText, className = '', confirmDisabled, cancelDisabled, confirmAction}: ModalFooterProps) : JSX.Element => (
    <RBModal.Footer className={confirmAction ? 'azd-modal__confirm-action' : `azd-modal__footer d-flex flex-column justify-content-center align-items-center ${className}`}>
        {onConfirm && (
            <button
                className={`plugin-btn btn ${confirmAction ? 'btn-danger' : 'btn-primary azd-modal__confirm-btn'}`}
                onClick={onConfirm}
                disabled={confirmDisabled}
            >
                {confirmBtnText || 'Confirm'}
            </button>
        )}
        {onHide && (
            <button
                className={`plugin-btn btn btn-link ${!confirmAction && 'azd-modal__cancel-btn'}`}
                onClick={onHide}
                disabled={cancelDisabled}
            >
                {cancelBtnText || 'Cancel'}
            </button>
        )}
    </RBModal.Footer>
);

export default ModalFooter;
