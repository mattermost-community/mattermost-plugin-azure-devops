import {createSlice, PayloadAction} from '@reduxjs/toolkit';

const initialState: ApiRequestCompletionState = {
    requestes: [],
};

export const apiRequestCompletionSlice = createSlice({
    name: 'globalModal',
    initialState,
    reducers: {
        setApiRequestCompletionState: (state: ApiRequestCompletionState, action: PayloadAction<ApiServiceName>) => {
            state.requestes = [...state.requestes, action.payload];
        },
        resetApiRequestCompletionState: (state: ApiRequestCompletionState, action: PayloadAction<ApiServiceName>) => {
            state.requestes = state.requestes.filter(((request) => request !== action.payload));
        },
    },
});

export const {setApiRequestCompletionState, resetApiRequestCompletionState} = apiRequestCompletionSlice.actions;

export default apiRequestCompletionSlice.reducer;
