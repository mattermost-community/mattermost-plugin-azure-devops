import React, {useEffect} from 'react';
import {useDispatch} from 'react-redux';

import SubscriptionCard from 'components/card/subscription';
import IconButton from 'components/buttons/iconButton';
import BackButton from 'components/buttons/backButton';

import {resetProjectDetails} from 'reducers/projectDetails';

// TODO: dummy data, remove later
const data: SubscriptionDetails[] = [
    {
        id: 'abc',
        name: 'Listen for all new tasks created',
        eventType: 'create',
    },
    {
        id: 'abc1',
        name: 'Listen for any task updated',
        eventType: 'update',
    },
    {
        id: 'abc2',
        name: 'Listen for all any task deleted',
        eventType: 'delete',
    },
];

type ProjectDetailsProps = {
    title: string
}

const ProjectDetails = ({title}: ProjectDetailsProps) => {
    const dispatch = useDispatch();

    const handleResetProjectDetails = () => {
        dispatch(resetProjectDetails());
    };

    // Reset the state when the component is unmounted
    useEffect(() => {
        return () => {
            handleResetProjectDetails();
        };
    }, []);

    return (
        <>
            <BackButton onClick={handleResetProjectDetails}/>
            <div className='d-flex'>
                <p className='rhs-title'>{title}</p>
                <IconButton
                    tooltipText='Unlink project'
                    iconClassName='fa fa-chain-broken'
                    extraClass='project-details-unlink-button unlink-button'
                />
            </div>
            <div className='bottom-divider'>
                <p className='font-size-14 font-bold margin-0 show-selected'>{'Subscriptions'}</p>
            </div>
            {
                data.map((item) => (
                    <SubscriptionCard
                        subscriptionDetails={{...item}}
                        key={item.id}
                    />
                ),
                )
            }
        </>
    );
};

export default ProjectDetails;
