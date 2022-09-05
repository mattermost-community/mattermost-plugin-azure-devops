import React, {useEffect} from 'react';
import {useDispatch} from 'react-redux';

import Modal from 'components/modal';
import Form from 'components/form';

import plugin_constants from 'plugin_constants';

import {toggleShowLinkModal} from 'reducers/linkModal';
import {getLinkModalState} from 'selectors';

import usePluginApi from 'hooks/usePluginApi';
import useForm from 'hooks/useForm';
import useApiRequestCompletionState from 'hooks/useApiRequestCompletionState';

const LinkModal = () => {
    const {linkProjectModal} = plugin_constants.form

    // Hooks
    const {
        formFields,
        errorState,
        setSpecificFieldValue,
        onChangeOfFormField,
        resetFormFields,
        isErrorInFormValidation,
    } = useForm(linkProjectModal);
    const {makeApiRequestWithCompletionStatus, state, getApiState} = usePluginApi();
    const dispatch = useDispatch();

    // State variables
    const {visibility, organization, project} = getLinkModalState(state);

    // Function to hide the modal and reset all the states.
    const resetModalState = (isActionDone?: boolean) => {
        dispatch(toggleShowLinkModal({isVisible: false, commandArgs: [], isActionDone}));
        resetFormFields();
    };

    // Handles on confirming link project
    const onConfirm = () => {
        if (!isErrorInFormValidation()) {
            // Make POST api request
            makeApiRequestWithCompletionStatus(
                plugin_constants.pluginApiServiceConfigs.createLink.apiServiceName,
                formFields as LinkPayload,
            );
        }
    };

    useApiRequestCompletionState({
        serviceName: plugin_constants.pluginApiServiceConfigs.createLink.apiServiceName,
        payload: formFields as LinkPayload,
        handleSuccess: () => resetModalState(true),
    });

    // Set modal field values
    useEffect(() => {
        setSpecificFieldValue({
            organization,
            project,
        });
    }, [visibility]);

    const {isLoading} = getApiState(plugin_constants.pluginApiServiceConfigs.createLink.apiServiceName, formFields as LinkPayload);

    return (
        <Modal
            show={visibility}
            title='Link New Project'
            onHide={resetModalState}
            onConfirm={onConfirm}
            confirmBtnText='Link new project'
            cancelDisabled={isLoading}
            confirmDisabled={isLoading}
            loading={isLoading}
        >
            <>
                {
                    Object.keys(linkProjectModal).map((field) => (
                        <Form
                            key={linkProjectModal[field as LinkProjectModalFields].label}
                            fieldConfig={linkProjectModal[field as LinkProjectModalFields]}
                            value={formFields[field as LinkProjectModalFields] ?? null}
                            onChange={(newValue) => onChangeOfFormField(field as LinkProjectModalFields, newValue)}
                            error={errorState[field as LinkProjectModalFields]}
                            isDisabled={isLoading}
                        />
                    ))
                }
            </>
        </Modal>
    );
};

export default LinkModal;
