import { PayloadAction, createSlice } from "@reduxjs/toolkit"
import { IGroup } from "../interfaces/group"

export interface IDashboardState {
    selectedGroup: IGroup | undefined
}

const dashboardInitialState: IDashboardState = {
    selectedGroup: undefined
}


export const dashboardSlice = createSlice({
    name: "dashboard",
    initialState: dashboardInitialState,
    reducers: {
        setSelectGroup: (state, action: PayloadAction<{group: IGroup}>) => {
            state.selectedGroup = action.payload.group
        }
    }
})


export const {setSelectGroup} = dashboardSlice.actions

const dashboardReducer = dashboardSlice.reducer

export default dashboardReducer