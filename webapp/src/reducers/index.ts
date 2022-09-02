import {combineReducers} from 'redux';

import services from 'services';

import globalModalSlice from './globalModal';
import openLinkModalSlice from './linkModal';
import openSubscribeModalSlice from './subscribeModal';
import openTaskModalSlice from './taskModal';
import projectDetailsSlice from './projectDetails';
import userConnectedSlice from './userConnected';
import testReducer from './testReducer';

const reducers = combineReducers({
    globalModalSlice,
    openLinkModalSlice,
    openTaskModalSlice,
    openSubscribeModalSlice,
    projectDetailsSlice,
    userConnectedSlice,
    testReducer,
    [services.reducerPath]: services.reducer,
});

export default reducers;
