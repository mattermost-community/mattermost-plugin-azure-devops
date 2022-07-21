import {combineReducers} from 'redux';

// TODO: for reference of developers, remove when actual dev work is done
import services from 'services';

import openTaskModalReducer from './taskModal';
import testReducer from './testReducer';

const reducers = combineReducers({
    openTaskModalReducer,
    testReducer,
    [services.reducerPath]: services.reducer,
});

export default reducers;
