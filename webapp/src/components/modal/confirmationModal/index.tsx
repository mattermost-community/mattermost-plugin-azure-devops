import React from 'react';

import Modal from 'components/modal';

type ConfirmationModalProps = {
    isOpen: boolean
    title: string
    description: string
    confirmBtnText: string
    onHide: () => void
    onConfirm?: () => void
    isLoading?: boolean
}

const ConfirmationModal = ({isOpen, title, confirmBtnText, description, onHide, onConfirm}: ConfirmationModalProps) => (
    <Modal
        show={isOpen}
        title={title}
        onHide={onHide}
        onConfirm={onConfirm}
        confirmAction={true}
        confirmBtnText={confirmBtnText}
    >
        <p>{description}</p>
    </Modal>
);

export default ConfirmationModal;
