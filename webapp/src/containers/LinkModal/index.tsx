import React, {useEffect, useState} from 'react';
import {useDispatch} from 'react-redux';

import Input from 'components/inputField';
import Modal from 'components/modal';

import usePluginApi from 'hooks/usePluginApi';
import {hideLinkModal, toggleIsLinked} from 'reducers/linkModal';
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
        dispatch(hideLinkModal());
    };

    // Set organization name
    const onOrganizationChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setProjectDetails({...projectDetails, organization: (e.target as HTMLInputElement).value});
    };

    // Set project name
    const onProjectChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setProjectDetails({...projectDetails, project: (e.target as HTMLInputElement).value});
    };

    // Handles on confirming link project
    const onConfirm = () => {
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
            return;
        }

        // Make POST api request
        linkTask(projectDetails);
    };

    // Make POST api request to link a project
    const linkTask = async (payload: LinkPayload) => {
        const createTaskRequest = await usePlugin.makeApiRequest(plugin_constants.pluginApiServiceConfigs.createLink.apiServiceName, payload);
        if (createTaskRequest) {
            // TODO: remove later
            // eslint-disable-next-line
            console.log('test', usePlugin.getApiState(plugin_constants.pluginApiServiceConfigs.createLink.apiServiceName, payload));
            dispatch(toggleIsLinked(true));
            resetModalState();
        }
    };

    useEffect(() => {
        setProjectDetails({
            organization: getLinkModalState(usePlugin.state).organization,
            project: getLinkModalState(usePlugin.state).project,
        });
    }, [getLinkModalState(usePlugin.state)]);

    return (
        <Modal
            show={getLinkModalState(usePlugin.state).visibility}
            title='Link new project'
            onHide={resetModalState}
            onConfirm={onConfirm}
            confirmBtnText='Link new project'
            cancelDisabled={usePlugin.getApiState(plugin_constants.pluginApiServiceConfigs.createLink.apiServiceName, projectDetails).isLoading}
            confirmDisabled={usePlugin.getApiState(plugin_constants.pluginApiServiceConfigs.createLink.apiServiceName, projectDetails).isLoading}
            loading={usePlugin.getApiState(plugin_constants.pluginApiServiceConfigs.createLink.apiServiceName, projectDetails).isLoading}
        >
            <>
                <Input
                    type='text'
                    placeholder='Organization name'
                    value={projectDetails.project}
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
