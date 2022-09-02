import React, {useCallback, useEffect, useState} from 'react';
import {useDispatch} from 'react-redux';

import Input from 'components/inputField';
import Modal from 'components/modal';

import usePluginApi from 'hooks/usePluginApi';
import {toggleShowLinkModal, toggleIsLinked} from 'reducers/linkModal';
import {getLinkModalState} from 'selectors';
import plugin_constants from 'plugin_constants';

const LinkModal = () => {
    // State variables
    const [projectDetails, setProjectDetails] = useState<LinkPayload>({
        organization: '',
        project: '',
    });

    const [errorState, setErrorState] = useState<LinkPayload>({
        organization: '',
        project: '',
    });

    // Hooks
    const usePlugin = usePluginApi();
    const dispatch = useDispatch();

    // Function to hide the modal and reset all the states.
    const resetModalState = () => {
        setProjectDetails({
            organization: '',
            project: '',
        });
        setErrorState({
            organization: '',
            project: '',
        });
        dispatch(toggleShowLinkModal({isVisible: false, commandArgs: []}));
    };

    // Set organization name
    const onOrganizationChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setErrorState({...errorState, organization: ''});
        setProjectDetails({...projectDetails, organization: (e.target as HTMLInputElement).value});
    };

    // Set project name
    const onProjectChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setErrorState({...errorState, project: ''});
        setProjectDetails({...projectDetails, project: (e.target as HTMLInputElement).value});
    };

    // Handles on confirming link project
    const onConfirm = useCallback(() => {
        const errorStateChanges: LinkPayload = {
            organization: '',
            project: '',
        };

        if (projectDetails.organization === '') {
            errorStateChanges.organization = 'Organization is required';
        }

        if (projectDetails.project === '') {
            errorStateChanges.project = 'Project is required';
        }

        if (errorStateChanges.organization || errorStateChanges.project) {
            setErrorState(errorStateChanges);
            return;
        }

        // Make POST api request
        linkTask(projectDetails);
    }, [errorState]);

    // Make POST API request to link a project
    const linkTask = async (payload: LinkPayload) => {
        const createTaskRequest = await usePlugin.makeApiRequest(plugin_constants.pluginApiServiceConfigs.createLink.apiServiceName, payload);
        if (createTaskRequest) {
            dispatch(toggleIsLinked(true));
            resetModalState();
        }
    };

    // Set modal field values
    useEffect(() => {
        setProjectDetails({
            organization: getLinkModalState(usePlugin.state).organization,
            project: getLinkModalState(usePlugin.state).project,
        });
    }, [getLinkModalState(usePlugin.state).visibility]);

    const {isLoading} = usePlugin.getApiState(plugin_constants.pluginApiServiceConfigs.createLink.apiServiceName, projectDetails);

    return (
        <Modal
            show={getLinkModalState(usePlugin.state).visibility}
            title='Link new project'
            onHide={resetModalState}
            onConfirm={onConfirm}
            confirmBtnText='Link new project'
            cancelDisabled={isLoading}
            confirmDisabled={isLoading}
            loading={isLoading}
        >
            <>
                <Input
                    type='text'
                    placeholder='Organization name'
                    value={projectDetails.organization}
                    onChange={onOrganizationChange}
                    error={errorState.organization}
                    required={true}
                />
                <Input
                    type='text'
                    placeholder='Project name'
                    value={projectDetails.project}
                    onChange={onProjectChange}
                    required={true}
                    error={errorState.project}
                />
            </>
        </Modal>
    );
};

export default LinkModal;
