import {combineReducers} from 'redux';

import services from 'services';

import globalModalSlice from './globalModal';
import apiRequestCompletionSlice from './apiRequest';
import linkProjectModalSlice from './linkModal';
import subscriptionModalSlice from './subscribeModal';
import openTaskModalReducer from './taskModal';
import projectDetailsSlice from './projectDetails';
import websocketEventSlice from './websocketEvent';

const reducers = combineReducers({
    apiRequestCompletionSlice,
    globalModalSlice,
    linkProjectModalSlice,
    openTaskModalReducer,
    subscriptionModalSlice,
    projectDetailsSlice,
    websocketEventSlice,
    [services.reducerPath]: services.reducer,
});

export default reducers;
