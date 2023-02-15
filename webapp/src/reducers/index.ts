import {combineReducers} from 'redux';

import {azureDevOpsPluginApi} from 'services';

import globalModalSlice from './globalModal';
import apiRequestCompletionSlice from './apiRequest';
import linkProjectModalSlice from './linkModal';
import subscriptionModalSlice from './subscribeModal';
import createTaskModalSlice from './taskModal';
import projectDetailsSlice from './projectDetails';
import websocketEventSlice from './websocketEvent';

const reducers = combineReducers({
    apiRequestCompletionSlice,
    globalModalSlice,
    linkProjectModalSlice,
    createTaskModalSlice,
    subscriptionModalSlice,
    projectDetailsSlice,
    websocketEventSlice,
    [azureDevOpsPluginApi.reducerPath]: azureDevOpsPluginApi.reducer,
});

export default reducers;
