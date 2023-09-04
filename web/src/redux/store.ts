import { configureStore } from "@reduxjs/toolkit";
import { TypedUseSelectorHook, useDispatch, useSelector } from "react-redux";
import authReducer from "./authSlice";
import dashboardReducer from "./dashboardSlice";

const store = configureStore({
    reducer: {
        auth: authReducer,
        dashboard: dashboardReducer
    }
})

export type RootState = ReturnType<typeof store.getState>


export type AppDispatch = typeof store.dispatch

export const useAppDispatch: () => AppDispatch = useDispatch
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector

export default store