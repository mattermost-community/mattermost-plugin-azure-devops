import React, {useCallback, useEffect, useState} from 'react';
import {useDispatch} from 'react-redux';

import Input from 'components/inputField';
import Modal from 'components/modal';

import Constants from 'plugin_constants';
import usePluginApi from 'hooks/usePluginApi';
import {hideLinkModal} from 'reducers/linkModal';

const LinkModal = () => {
    const [state, setState] = useState({
        linkOrganization: '',
        linkProject: '',
    });
    const [linkOrganizationError, setLinkOrganizationError] = useState('');
    const [linkProjectError, setLinkProjectError] = useState('');
    const [linkPayload, setLinkPayload] = useState<LinkPayload | null>();
    const usePlugin = usePluginApi();
    const {visibility, organization, project} = usePlugin.state['plugins-mattermost-plugin-azure-devops'].openLinkModalReducer;
    const dispatch = useDispatch();

    useEffect(() => {
        if (organization && project) {
            setState({
                linkOrganization: organization,
                linkProject: project,
            });
        }
    }, [visibility]);

    // Function to hide the modal and reset all the states.
    const onHide = useCallback(() => {
        setState({
            linkOrganization: '',
            linkProject: '',
        });
        setLinkOrganizationError('');
        setLinkProjectError('');
        setLinkPayload(null);
        dispatch(hideLinkModal());
    }, []);

    const onOrganizationChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
        setLinkOrganizationError('');
        setState({...state, linkOrganization: (e.target as HTMLInputElement).value});
    }, [state]);

    const onProjectChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
        setLinkProjectError('');
        setState({...state, linkProject: (e.target as HTMLInputElement).value});
    }, [state]);

    const onConfirm = useCallback(() => {
        if (state.linkOrganization === '') {
            setLinkOrganizationError('Organization is required');
        }
        if (state.linkProject === '') {
            setLinkProjectError('Project is required');
        }

        if (!state.linkOrganization || !state.linkProject) {
            return;
        }

        // Create payload to send in the POST request.
        const payload = {
            organization: state.linkOrganization,
            project: state.linkProject,
        };

        // TODO: save the payload in a state variable to use it while reading the state
        // we can see later if there exists a better way for this
        setLinkPayload(payload);

        // Make POST api request
        usePlugin.makeApiRequest(Constants.pluginApiServiceConfigs.createLink.apiServiceName, payload);
    }, [state, linkOrganizationError, linkProjectError]);

    useEffect(() => {
        if (linkPayload) {
            const {isSuccess, isError} = usePlugin.getApiState(Constants.pluginApiServiceConfigs.createLink.apiServiceName, linkPayload);
            console.log(isSuccess, isError);
            if ((isSuccess && !isError) || (!isSuccess && isError)) {
                onHide();
            }
        }
    }, [usePlugin.state]);

    if (visibility) {
        return (
            <Modal
                show={visibility}
                title='Link new project'
                onHide={onHide}
                onConfirm={onConfirm}
                confirmBtnText='Link new project'
                loading={linkPayload ? usePlugin.getApiState(Constants.pluginApiServiceConfigs.createLink.apiServiceName, linkPayload).isLoading : false}
            >
                <>
                    <Input
                        type='text'
                        placeholder='Organization name'
                        value={state.linkOrganization}
                        onChange={onOrganizationChange}
                        error={linkOrganizationError}
                        required={true}
                    />
                    <Input
                        type='text'
                        placeholder='Project name'
                        value={state.linkProject}
                        onChange={onProjectChange}
                        required={true}
                        error={linkProjectError}
                    />
                </>
            </Modal>
        );
    }
    return null;
};

export default LinkModal;
