import React, {useEffect, useState} from 'react';
import {useDispatch} from 'react-redux';

import Input from 'components/inputField';
import Modal from 'components/modal';

import usePluginApi from 'hooks/usePluginApi';
import {hideLinkModal, toggleIsLinked} from 'reducers/linkModal';
import {getLinkModalState, getUserConnectionState} from 'selectors';
import plugin_constants from 'plugin_constants';
import LinearLoader from 'components/loader/linear';

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
        setErrorState({...errorState, organization: ''});
        setProjectDetails({...projectDetails, organization: (e.target as HTMLInputElement).value});
    };

    // Set project name
    const onProjectChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setErrorState({...errorState, project: ''});
        setProjectDetails({...projectDetails, project: (e.target as HTMLInputElement).value});
    };

    // Handles on confirming link project
    const onConfirm = () => {
        const newErrorState: LinkPayload = {
            organization: '',
            project: '',
        };

        if (projectDetails.organization === '') {
            newErrorState.organization = 'Organization is required';
        }

        if (projectDetails.project === '') {
            newErrorState.project = 'Project is required';
        }

        if (newErrorState.organization || newErrorState.project) {
            setErrorState(newErrorState);
            return;
        }

        // Make POST api request
        linkTask(projectDetails);
    };

    // Make POST api request to link a project
    const linkTask = async (payload: LinkPayload) => {
        const createTaskRequest = await usePlugin.makeApiRequest(plugin_constants.pluginApiServiceConfigs.createLink.apiServiceName, payload);
        if (createTaskRequest) {
            dispatch(toggleIsLinked(true));
            resetModalState();
        }
    };

    useEffect(() => {
        if (getLinkModalState(usePlugin.state).visibility) {
            setProjectDetails({
                organization: getLinkModalState(usePlugin.state).organization,
                project: getLinkModalState(usePlugin.state).project,
            });
        }
    }, [getLinkModalState(usePlugin.state).visibility]);

    return (
        <Modal
            show={getLinkModalState(usePlugin.state).visibility}
            title='Link new project'
            onHide={resetModalState}
            onConfirm={onConfirm}
            showFooter={
                !usePlugin.getUserAccountConnectionState().isLoading &&
                usePlugin.getUserAccountConnectionState().isSuccess &&
                usePlugin.getUserAccountConnectionState().data?.MattermostUserID
            }
            confirmBtnText='Link new project'
            cancelDisabled={usePlugin.getApiState(plugin_constants.pluginApiServiceConfigs.createLink.apiServiceName, projectDetails).isLoading}
            confirmDisabled={usePlugin.getApiState(plugin_constants.pluginApiServiceConfigs.createLink.apiServiceName, projectDetails).isLoading}
            loading={usePlugin.getApiState(plugin_constants.pluginApiServiceConfigs.createLink.apiServiceName, projectDetails).isLoading}
        >
            <>
                {
                    usePlugin.getUserAccountConnectionState().isLoading && (<LinearLoader/>)
                }
                {
                    !usePlugin.getUserAccountConnectionState().isLoading &&
                    usePlugin.getUserAccountConnectionState().isError &&
                    (<div className='not-linked'>{'You do not have any Azure Devops account connected. Kindly link the account first'}</div>)
                }
                {
                    !usePlugin.getUserAccountConnectionState().isLoading &&
                    usePlugin.getUserAccountConnectionState().isSuccess &&
                    usePlugin.getUserAccountConnectionState().data?.MattermostUserID && (
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
                    )
                }
            </>
        </Modal>
    );
};

export default LinkModal;
