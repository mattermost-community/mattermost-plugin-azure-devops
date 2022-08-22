import {combineReducers} from 'redux';

import services from 'services';

import globalModalSlice from './globalModal';
import apiRequestCompletionSlice from './apiRequest';
import openLinkModalSlice from './linkModal';
import openSubscribeModalSlice from './subscribeModal';
import openTaskModalReducer from './taskModal';
import projectDetailsSlice from './projectDetails';
import userConnectedSlice from './userConnected';

const reducers = combineReducers({
    apiRequestCompletionSlice,
    globalModalSlice,
    openLinkModalSlice,
    openTaskModalReducer,
    openSubscribeModalSlice,
    projectDetailsSlice,
    userConnectedSlice,
    [services.reducerPath]: services.reducer,
});

export default reducers;
