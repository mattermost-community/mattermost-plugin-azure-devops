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
    const [error, setError] = useState({
        linkOrganizationError: '',
        linkProjectError: '',
    });
    const [linkPayload, setLinkPayload] = useState<LinkPayload | null>();
    const usePlugin = usePluginApi();
    const {visibility, organization, project} = usePlugin.state['plugins-mattermost-plugin-azure-devops'].openLinkModalReducer;
    const [loading, setLoading] = useState(false);
    const [APIError, setAPIError] = useState('');
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
        setError({
            linkOrganizationError: '',
            linkProjectError: '',
        });
        setLinkPayload(null);
        setLoading(false);
        setAPIError('');
        dispatch(hideLinkModal());
    }, []);

    const onOrganizationChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
        setError({...error, linkOrganizationError: ''});
        setState({...state, linkOrganization: (e.target as HTMLInputElement).value});
    }, [state, error]);

    const onProjectChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
        setError({...error, linkProjectError: ''});
        setState({...state, linkProject: (e.target as HTMLInputElement).value});
    }, [state, error]);

    const onConfirm = useCallback(() => {
        if (state.linkOrganization === '') {
            setError((value) => ({...value, linkOrganizationError: 'Organization is required'}));
        }
        if (state.linkProject === '') {
            setError((value) => ({...value, linkProjectError: 'Project is required'}));
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
    }, [state]);

    useEffect(() => {
        if (linkPayload) {
            const {isLoading, isSuccess, isError} = usePlugin.getApiState(Constants.pluginApiServiceConfigs.createLink.apiServiceName, linkPayload);
            setLoading(isLoading);
            if (isSuccess) {
                onHide();
            }
            if (isError) {
                setAPIError('Organization or project name is wrong');
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
                cancelDisabled={loading}
                confirmDisabled={loading}
                loading={loading}
                error={APIError}
            >
                <>
                    <Input
                        type='text'
                        placeholder='Organization name'
                        value={state.linkOrganization}
                        onChange={onOrganizationChange}
                        error={error.linkOrganizationError}
                        required={true}
                    />
                    <Input
                        type='text'
                        placeholder='Project name'
                        value={state.linkProject}
                        onChange={onProjectChange}
                        required={true}
                        error={error.linkProjectError}
                    />
                </>
            </Modal>
        );
    }
    return null;
};

export default LinkModal;
