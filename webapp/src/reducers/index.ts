import {combineReducers} from 'redux';

// TODO: for reference of developers, remove when actual dev work is done
import services from 'services';

import testReducer from './testReducer';

export default combineReducers({
    testReducer,
    [services.reducerPath]: services.reducer,
});
