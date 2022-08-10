import {combineReducers} from 'redux';

import services from 'services';

import openLinkModalReducer from './linkModal';
import openSubscribeModalReducer from './subscribeModal';
import openTaskModalReducer from './taskModal';
import projectDetailsSlice from './projectDetails';
import userConnectionSlice from './userAcountDetails';

const reducers = combineReducers({
    openLinkModalReducer,
    openTaskModalReducer,
    projectDetailsSlice,
    userConnectionSlice,
    openSubscribeModalReducer,
    [services.reducerPath]: services.reducer,
});

export default reducers;
