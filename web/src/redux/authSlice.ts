import { PayloadAction, createSlice } from "@reduxjs/toolkit";
import { IUser } from "../interfaces/user";

export interface IAuthState {
    user: IUser | undefined
    init: boolean
    isLoggedIn: boolean
}

const authInitialState: IAuthState  =  {
    user: undefined,
    init: false,
    isLoggedIn: false
}

export const authSlice = createSlice({
    name: "auth",
    initialState: authInitialState,
    reducers: {
        init: (state, action: PayloadAction<{user: IUser | undefined}>) => {
            state.user = action.payload.user
            state.isLoggedIn = action.payload.user !== undefined
            state.init = true
        },
        updateUser: (state, action: PayloadAction<{user: IUser}>) => {
            state.user = action.payload.user
        },
        logout: (state) => {
            state.user = undefined
            state.isLoggedIn = false
        }
    }
})

export const {init, updateUser, logout} = authSlice.actions


const authReducer =  authSlice.reducer

export default authReducer