import React from 'react';
import {Modal as RBModal} from 'react-bootstrap';

import ModalHeader from './subComponents/modalHeader';
import ModalLoader from './subComponents/modalLoader';
import ModalBody from './subComponents/modalBody';
import ModalFooter from './subComponents/modalFooter';
import ModalSubTitleAndError from './subComponents/modalSubtitleAndError';

type ModalProps = {
    show: boolean;
    onHide?: () => void;
    showCloseIconInHeader?: boolean;
    children?: JSX.Element;
    title?: string | JSX.Element;
    subTitle?: string | JSX.Element;
    onConfirm?: () => void;
    confirmBtnText?: string;
    cancelBtnText?: string;
    confirmAction?: boolean;
    className?: string;
    loading?: boolean;
    error?: string | JSX.Element;
    confirmDisabled?: boolean;
    cancelDisabled?: boolean;
}

const Modal = ({show, onHide, showCloseIconInHeader = true, children, title, subTitle, onConfirm, confirmAction, confirmBtnText, cancelBtnText, className = '', loading = false, error, confirmDisabled = false, cancelDisabled = false}: ModalProps) => (
    <RBModal
        show={show}
        onHide={onHide}
        centered={true}
        className={`modal ${className}`}
    >
        <ModalHeader
            title={title}
            showCloseIconInHeader={showCloseIconInHeader}
            onHide={onHide}
        />
        <ModalLoader loading={loading}/>
        <ModalBody>
            <>
                <ModalSubTitleAndError
                    subTitle={subTitle}
                />
                {children}
                <ModalSubTitleAndError
                    error={error}
                />
            </>
        </ModalBody>
        <ModalFooter
            onHide={onHide}
            onConfirm={onConfirm}
            cancelBtnText={cancelBtnText}
            confirmBtnText={confirmBtnText}
            confirmAction={confirmAction}
            confirmDisabled={confirmDisabled}
            cancelDisabled={cancelDisabled}
        />
    </RBModal>
);

export default Modal;
