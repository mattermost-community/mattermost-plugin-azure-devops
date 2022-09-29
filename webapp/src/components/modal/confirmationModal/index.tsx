import React from 'react';

import Modal from 'components/modal';
import ResultPanel from 'components/resultPanel';

type ConfirmationModalProps = {
    isOpen: boolean
    title: string
    description: string
    confirmBtnText: string
    onHide: () => void
    onConfirm?: () => void
    isLoading?: boolean
    showErrorPanel?: ConfirmationModalErrorPanelProps | null
}

const ConfirmationModal = ({isOpen, title, confirmBtnText, description, onHide, onConfirm, isLoading, showErrorPanel}: ConfirmationModalProps) => (
    <Modal
        show={isOpen}
        title={title}
        onHide={onHide}
        onConfirm={onConfirm}
        confirmAction={true}
        confirmBtnText={confirmBtnText}
        loading={isLoading}
        showFooter={!showErrorPanel}
    >
        {
            showErrorPanel ? (
                <ResultPanel
                    iconClass={'fa-times-circle-o result-panel-icon--error'}
                    header={showErrorPanel.title}
                    secondaryBtnText='Close'
                    onSecondaryBtnClick={showErrorPanel.onSecondaryBtnClick}
                />
            ) : <p>{description}</p>
        }
    </Modal>
);

export default ConfirmationModal;
