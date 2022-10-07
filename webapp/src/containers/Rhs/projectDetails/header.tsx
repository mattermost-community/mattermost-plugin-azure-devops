import React, {useState} from 'react';
import {useDispatch} from 'react-redux';

import plugin_constants from 'plugin_constants';

import ConfirmationModal from 'components/modal/confirmationModal';
import BackButton from 'components/buttons/backButton';
import PrimaryButton from 'components/buttons/primaryButton';
import ToggleSwitch from 'components/toggleSwitch';
import Dropdown from 'components/dropdown';

import {toggleIsLinkedProjectListChanged} from 'reducers/linkModal';

import usePluginApi from 'hooks/usePluginApi';
import useApiRequestCompletionState from 'hooks/useApiRequestCompletionState';

type HeaderProps = {
    projectDetails: ProjectDetails
    showAllSubscriptions: boolean
    setShowAllSubscriptions: (active: boolean) => void
    handlePagination: (reset: boolean) => void
    handleResetProjectDetails: () => void
    filter: string
    setFilter: (filter: string) => void
    setSubscriptionList: (subscriptionDetails: SubscriptionDetails[]) => void
}

const Header = ({projectDetails, showAllSubscriptions, handlePagination, setShowAllSubscriptions, handleResetProjectDetails, filter, setFilter, setSubscriptionList}: HeaderProps) => {
    const {projectName} = projectDetails;
    const [showProjectConfirmationModal, setShowProjectConfirmationModal] = useState(false);

    // const [showFilter, setShowFilter] = useState(false); // TODO: uncomment when need to toggle filter

    const dispatch = useDispatch();
    const {makeApiRequestWithCompletionStatus, getApiState, state} = usePluginApi();

    // Opens a confirmation modal to confirm unlinking a project
    const handleUnlinkProject = () => {
        setShowProjectConfirmationModal(true);
    };

    // Handles unlinking a project and fetching the modified project list
    const handleConfirmUnlinkProject = () => {
        makeApiRequestWithCompletionStatus(plugin_constants.pluginApiServiceConfigs.unlinkProject.apiServiceName, projectDetails);
    };

    useApiRequestCompletionState({
        serviceName: plugin_constants.pluginApiServiceConfigs.unlinkProject.apiServiceName,
        payload: projectDetails,
        handleSuccess: () => {
            dispatch(toggleIsLinkedProjectListChanged(true));
            handleResetProjectDetails();
            setShowProjectConfirmationModal(false);
        },
    });

    const {isLoading: isUnlinkProjectLoading} = getApiState(plugin_constants.pluginApiServiceConfigs.unlinkProject.apiServiceName, projectDetails);

    return (
        <>
            <ConfirmationModal
                isOpen={showProjectConfirmationModal}
                onHide={() => setShowProjectConfirmationModal(false)}
                onConfirm={handleConfirmUnlinkProject}
                isLoading={isUnlinkProjectLoading}
                confirmBtnText='Unlink'
                description={`Are you sure you want to unlink ${projectName}?`}
                title='Confirm Project Unlink'
            />
            <div className='rhs-header-divider'>
                <div className='d-flex align-item-center margin-bottom-15'>
                    <BackButton onClick={handleResetProjectDetails}/>
                    <p className='rhs-title'>{projectName}</p>
                    <PrimaryButton
                        text='Unlink'
                        iconName='fa fa-chain-broken'
                        extraClass='rhs-project-details-unlink-button'
                        onClick={handleUnlinkProject}
                    />
                </div>
                <div className='d-flex align-item-center'>
                    <ToggleSwitch
                        active={showAllSubscriptions}
                        onChange={(active) => {
                            handlePagination(true);
                            setShowAllSubscriptions(active);
                        }}
                        label={'Show For All Channels'}
                        labelPositioning='right'
                    />
                    {/* TODO: uncomment when need to toggle filter */}
                    {/* <IconButton
                        tooltipText='Filter'
                        extraClass='margin-left-auto flex-basis-initial'
                        onClick={() => setShowFilter(!showFilter)}
                    >
                        <SVGWrapper
                            width={18}
                            height={12}
                            viewBox='0 0 18 12'
                        >
                            {plugin_constants.SVGIcons.filter}
                        </SVGWrapper>
                    </IconButton> */}
                </div>
                <div className='filter-dropdown-container filter-dropdown-container__show'>
                    <Dropdown
                        placeholder='Created By'
                        value={filter}
                        onChange={(newValue) => {
                            setFilter(newValue);
                            setSubscriptionList([]);
                        }}
                        options={plugin_constants.form.subscriptionFilterOptions}
                        disabled={false}
                    />
                </div>
            </div>
        </>
    );
};

export default Header;
