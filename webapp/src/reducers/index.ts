import {combineReducers} from 'redux';

import services from 'services';

import openLinkModalReducer from './linkModal';
import openTaskModalReducer from './taskModal';
import projectDetailsSlice from './projectDetails';
import userConnectionSlice from './userAcountDetails';

const reducers = combineReducers({
    openLinkModalReducer,
    openTaskModalReducer,
    projectDetailsSlice,
    userConnectionSlice,
    [services.reducerPath]: services.reducer,
});

export default reducers;
