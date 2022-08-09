import {combineReducers} from 'redux';

import services from 'services';

import openTaskModalReducer from './taskModal';
import projectDetailsSlice from './projectDetails';
import testReducer from './testReducer';

const reducers = combineReducers({
    openTaskModalReducer,
    projectDetailsSlice,
    testReducer,
    [services.reducerPath]: services.reducer,
});

export default reducers;
