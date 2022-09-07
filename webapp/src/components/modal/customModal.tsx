import React from 'react';
import {Modal as RBModal} from 'react-bootstrap';

type ModalProps = {
    show: boolean;
    onHide: () => void;
    children: JSX.Element;
    className?: string;
}

const CustomModal = ({show, onHide, children, className = ''}: ModalProps) => (
    <RBModal
        show={show}
        onHide={onHide}
        className={`modal ${className}`}
    >
        {children}
    </RBModal>
);

export default CustomModal;
