import {createSlice, PayloadAction} from '@reduxjs/toolkit';

const initialState: ApiRequestCompletionState = {
    serviceName: '',
};

export const apiRequestCompletionSlice = createSlice({
    name: 'globalModal',
    initialState,
    reducers: {
        setApiRequestCompletionState: (state: ApiRequestCompletionState, action: PayloadAction<ApiRequestCompletionState>) => {
            state.serviceName = action.payload.serviceName;
        },
        resetApiRequestCompletionState: (state: ApiRequestCompletionState) => {
            state.serviceName = '';
        },
    },
});

export const {setApiRequestCompletionState, resetApiRequestCompletionState} = apiRequestCompletionSlice.actions;

export default apiRequestCompletionSlice.reducer;
