import React, {useRef, useState} from 'react';
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
import IconButton from 'components/buttons/iconButton';
import SVGWrapper from 'components/svgWrapper';
import useOutsideClick from 'hooks/useClickOutside';

type HeaderProps = {
    projectDetails: ProjectDetails
    showAllSubscriptions: boolean
    setShowAllSubscriptions: (active: boolean) => void
    handlePagination: (reset: boolean) => void
    handleResetProjectDetails: () => void
    filter: SubscriptionFilters
    setFilter: (filter: SubscriptionFilters) => void
    setSubscriptionList: (subscriptionDetails: SubscriptionDetails[]) => void
}

const Header = ({projectDetails, showAllSubscriptions, handlePagination, setShowAllSubscriptions, handleResetProjectDetails, filter, setFilter, setSubscriptionList}: HeaderProps) => {
    const {projectName} = projectDetails;
    const {subscriptionFilters} = plugin_constants.common;
    const [showProjectConfirmationModal, setShowProjectConfirmationModal] = useState(false);

    const [showFilter, setShowFilter] = useState(false);

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

    const isFilterApplied = () => showAllSubscriptions || filter.createdBy !== subscriptionFilters.createdBy.me || filter.serviceType !== subscriptionFilters.serviceType.boards;

    const {isLoading: isUnlinkProjectLoading} = getApiState(plugin_constants.pluginApiServiceConfigs.unlinkProject.apiServiceName, projectDetails);

    const wrapperRef = useRef(null);
    useOutsideClick(wrapperRef, () => {
        setShowFilter(false);
    });

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
            <div className='position-relative rhs-header-divider'>
                <div className='d-flex align-item-center'>
                    <BackButton onClick={handleResetProjectDetails}/>
                    <p className='rhs-title'>{projectName}</p>
                    <IconButton
                        tooltipText='Filter'
                        extraClass={`margin-left-auto flex-basis-initial ${isFilterApplied() && 'filter-button'}`}
                        onClick={() => setShowFilter(!showFilter)}
                    >
                        <SVGWrapper
                            width={18}
                            height={12}
                            viewBox='0 0 18 12'
                        >
                            {plugin_constants.SVGIcons.filter}
                        </SVGWrapper>
                    </IconButton>
                    <PrimaryButton
                        text='Unlink'
                        iconName='fa fa-chain-broken'
                        extraClass='margin-left-5'
                        onClick={handleUnlinkProject}
                    />
                </div>
            </div>
            {
                showFilter && (
                    <div
                        ref={wrapperRef}
                        className='rhs-filter-popover'
                    >
                        <div className='d-flex align-item-center margin-bottom-15'>
                            <ToggleSwitch
                                active={showAllSubscriptions}
                                onChange={(active) => {
                                    handlePagination(true);
                                    setShowAllSubscriptions(active);
                                }}
                                label={'Show For All Channels'}
                                labelPositioning='right'
                            />
                        </div>
                        <div className='margin-bottom-15'>
                            <Dropdown
                                placeholder='Created By'
                                value={filter.createdBy}
                                onChange={(newValue) => {
                                    setFilter({...filter, createdBy: newValue});
                                    setSubscriptionList([]);
                                }}
                                options={plugin_constants.form.subscriptionFilterCreatedByOptions}
                                disabled={false}
                            />
                        </div>
                        <Dropdown
                            placeholder='Service Type'
                            value={filter.serviceType}
                            onChange={(newValue) => {
                                setFilter({...filter, serviceType: newValue});

                                // TODO: add logic of refreshing list while integration
                            }}
                            options={plugin_constants.form.subscriptionFilterServiceTypeOptions}
                            disabled={false}
                        />
                    </div>
                )
            }
        </>
    );
};

export default Header;
