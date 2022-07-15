import React from 'react';
import {Modal as RBModal} from 'react-bootstrap';

import './styles.scss';

type ModalProps = {
    show: boolean;
    onHide: () => void;
    showCloseIconInHeader?: boolean;
    children?: JSX.Element;
    title?: string | JSX.Element;
    subTitle?: string | JSX.Element;
    onConfirm?: () => void;
    confirmBtnText?: string;
    cancelBtnText?: string;
    className?: string;
}

const Modal = ({show, onHide, showCloseIconInHeader = true, children, title, subTitle, onConfirm, confirmBtnText, cancelBtnText, className = ''}: ModalProps) => {
    return (
        <RBModal
            show={show}
            onHide={onHide}
            centered={true}
            className={`modal ${className}`}
        >
            {(title || subTitle || showCloseIconInHeader) && (
                <div className='modal__header'>
                    <div className='modal__title d-flex align-items-center justify-content-between'>
                        {title && <p className='modal__title'>{title}</p>}
                        {showCloseIconInHeader && (
                            <button
                                className='style--none'
                                onClick={onHide}
                            ><i className='icon icon-close modal__close-icon'/></button>
                        )}
                    </div>
                    {subTitle && <p className='modal__subtitle'>{subTitle}</p>}
                </div>
            )}
            {children && (
                <RBModal.Body className='modal__body'>
                    {children}
                </RBModal.Body>
            )}
            <RBModal.Footer className='modal__footer d-flex flex-column justify-content-center align-items-center'>
                {onConfirm && (
                    <button
                        className='btn btn-primary modal__confirm-btn'
                        onClick={onConfirm}
                    >
                        {confirmBtnText || 'Confirm'}
                    </button>
                )}
                <button
                    className='btn btn-link modal__cancel-btn'
                    onClick={onHide}
                >
                    {cancelBtnText || 'Cancel'}
                </button>
            </RBModal.Footer>
        </RBModal>
    );
};

export default Modal;
