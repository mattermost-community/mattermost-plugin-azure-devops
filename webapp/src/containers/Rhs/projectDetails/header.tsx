import React, {ChangeEvent, useCallback, useEffect, useRef, useState} from 'react';
import {useDispatch} from 'react-redux';

import pluginConstants from 'pluginConstants';

import ConfirmationModal from 'components/modal/confirmationModal';
import BackButton from 'components/buttons/backButton';
import Button from 'components/buttons/button';
import ToggleSwitch from 'components/toggleSwitch';
import Dropdown from 'components/dropdown';

import {toggleIsLinkedProjectListChanged} from 'reducers/linkModal';

import usePluginApi from 'hooks/usePluginApi';
import useApiRequestCompletionState from 'hooks/useApiRequestCompletionState';
import IconButton from 'components/buttons/iconButton';
import SVGWrapper from 'components/svgWrapper';
import useOutsideClick from 'hooks/useClickOutside';
import {subscriptionFilterEventTypeReposOptions} from 'pluginConstants/form';
import utils from 'utils';

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
    const {defaultSubscriptionFilters, subscriptionFilters, filterLabelValuePairAll} = pluginConstants.common;
    const {subscriptionFilterCreatedByOptions, subscriptionFilterServiceTypeOptions, subscriptionFilterEventTypeBoardsOptions, subscriptionModal} = pluginConstants.form;
    const [showProjectConfirmationModal, setShowProjectConfirmationModal] = useState(false);
    const [deleteSubscriptions, setDeleteSubscriptions] = useState(false);
    const [showDeleteSubscriptionsCheckbox, setShowDeleteSubscriptionsCheckbox] = useState(true);
    const [confirmationModalDescription, setConfirmationModalDescription] = useState(`Are you sure you want to unlink ${projectName}?`);
    const [unlinkConfirmationModalError, setUnlinkConfirmationModalError] = useState<ConfirmationModalErrorPanelProps | null>(null);

    const [showFilter, setShowFilter] = useState(false);

    const dispatch = useDispatch();
    const {makeApiRequestWithCompletionStatus, getApiState, state} = usePluginApi();

    const getEventTypeOptions = useCallback((serviceType: string) => {
        switch (serviceType) {
        case subscriptionFilters.serviceType.boards:
            return subscriptionFilterEventTypeBoardsOptions();
        case subscriptionFilters.serviceType.repos:
            return subscriptionFilterEventTypeReposOptions();
        default:
            return [filterLabelValuePairAll];
        }
    }, [filter.serviceType]);

    // Opens a confirmation modal to confirm unlinking a project
    const handleUnlinkProject = () => {
        setShowProjectConfirmationModal(true);
    };

    // Handles unlinking a project and fetching the modified project list
    const handleConfirmUnlinkProject = () => {
        makeApiRequestWithCompletionStatus(pluginConstants.pluginApiServiceConfigs.unlinkProject.apiServiceName, {...projectDetails, deleteSubscriptions} as APIRequestPayload);
    };

    // Update the modal when project unlinking fails
    const handleActionsAfterUnlinkingProjectFailed = (err: ApiErrorResponse) => {
        const errorMessage = utils.getErrorMessage(true, 'ConfirmationModal', err);
        if (errorMessage === pluginConstants.messages.error.adminAccessError) {
            setConfirmationModalDescription(pluginConstants.messages.error.adminAccessErrorForUnlinking);
        }

        setUnlinkConfirmationModalError({
            title: errorMessage,
            onSecondaryBtnClick: () => {
                setShowProjectConfirmationModal(false);
                setUnlinkConfirmationModalError(null);
                setShowDeleteSubscriptionsCheckbox(true);
            },
        });

        setDeleteSubscriptions(false);
        setShowDeleteSubscriptionsCheckbox(false);
    };

    useApiRequestCompletionState({
        serviceName: pluginConstants.pluginApiServiceConfigs.unlinkProject.apiServiceName,
        payload: {...projectDetails, deleteSubscriptions} as APIRequestPayload,
        handleSuccess: () => {
            dispatch(toggleIsLinkedProjectListChanged(true));
            handleResetProjectDetails();
            setShowProjectConfirmationModal(false);
        },
        handleError: handleActionsAfterUnlinkingProjectFailed,
    });

    const isFilterApplied = useCallback(() => showAllSubscriptions || filter.createdBy !== defaultSubscriptionFilters.createdBy || filter.serviceType !== defaultSubscriptionFilters.serviceType || filter.eventType !== defaultSubscriptionFilters.eventType, [filter, showAllSubscriptions]);

    const {isLoading: isUnlinkProjectLoading} = getApiState(pluginConstants.pluginApiServiceConfigs.unlinkProject.apiServiceName, {...projectDetails, deleteSubscriptions} as APIRequestPayload);

    // Detects and closes the filter popover whenever it is opened and the user clicks outside of it
    const wrapperRef = useRef(null);
    useOutsideClick(wrapperRef, () => {
        setShowFilter(false);
    });

    // Whenever "serviceType" changes make "all" option as default in "eventType"
    useEffect(() => {
        if (filter.eventType !== subscriptionFilters.eventType.all) {
            setFilter({...filter, eventType: subscriptionFilters.eventType.all});
        }
    }, [filter.serviceType]);

    const handleCheckboxChange = (e: ChangeEvent<HTMLInputElement>) => {
        setDeleteSubscriptions(e.target.checked);
    };

    const deleteSubscriptionsCheckbox = (
        <div className='d-flex align-item-center'>
            <input
                type='checkbox'
                id='deleteSubscriptions'
                className='margin-0'
                onChange={handleCheckboxChange}
            />
            <label className='margin-left-5 margin-bottom-0 font-weight-normal'>{pluginConstants.common.deleteAllSubscriptionsMessage}</label>
        </div>
    );

    return (
        <>
            <ConfirmationModal
                isOpen={showProjectConfirmationModal}
                onHide={() => setShowProjectConfirmationModal(false)}
                onConfirm={handleConfirmUnlinkProject}
                isLoading={isUnlinkProjectLoading}
                confirmBtnText='Unlink'
                description={confirmationModalDescription}
                title='Confirm Project Unlink'
                showErrorPanel={unlinkConfirmationModalError}
            >
                {showDeleteSubscriptionsCheckbox ? deleteSubscriptionsCheckbox : <></>}
            </ConfirmationModal>
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
                            {pluginConstants.SVGIcons.filter}
                        </SVGWrapper>
                    </IconButton>
                    <Button
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
                                options={subscriptionFilterCreatedByOptions}
                                disabled={false}
                            />
                        </div>
                        <div className='margin-bottom-15'>
                            <Dropdown
                                placeholder='Service Type'
                                value={filter.serviceType}
                                onChange={(newValue) => {
                                    setFilter({...filter, serviceType: newValue});
                                    setSubscriptionList([]);
                                }}
                                options={subscriptionFilterServiceTypeOptions}
                                disabled={false}
                            />
                        </div>
                        <div className='margin-bottom-15'>
                            <Dropdown
                                placeholder='Event Type'
                                value={filter.eventType}
                                onChange={(newValue) => {
                                    setFilter({...filter, eventType: newValue});
                                    setSubscriptionList([]);
                                }}
                                options={getEventTypeOptions(filter.serviceType)}
                                disabled={filter.serviceType === subscriptionFilters.serviceType.all}
                            />
                        </div>
                        <div className='text-align-right'>
                            <Button
                                text='Reset'
                                onClick={() => {
                                    setFilter(defaultSubscriptionFilters);
                                    setShowAllSubscriptions(false);
                                }}
                                extraClass='margin-right-8'
                                isSecondaryButton={true}
                                isDisabled={!isFilterApplied()}
                            />
                            <Button
                                text='Close'
                                onClick={() => setShowFilter(false)}
                            />
                        </div>
                    </div>
                )
            }
        </>
    );
};

export default Header;
