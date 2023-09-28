import axios from "axios";
import { ILoginRequest, ILoginResponse } from "../interfaces/login";
import api from "./api";
import { IRegisterRequest, IUser } from "../interfaces/user";

export const login = (payload: ILoginRequest) => {
  return api
    .post("/auth/login", payload)
    .then((res) => {
        const data = res.data as ILoginResponse
        // save token to user storage
        localStorage.setItem("token", data.token)
        // save token to api headers 
        api.defaults.headers.common.Authorization = `Bearer ${data.token}`
        return data
    })
    .catch((err) => {
        let errMessage: string
        if(axios.isAxiosError<{error: string}>(err) && err.response){
          
            errMessage = err.response.data.error
           
        }else {
            console.error(err)
            errMessage = "Something went wrong please try again"
        }

        return Promise.reject(errMessage)
      
    });
};


export const getMe = () => {
    return api.get("/user/me").then(res => res.data as IUser)
}


export const registerUser = (payload: IRegisterRequest) => {
    return api.post("/auth/register", payload).then(res => res.data as IUser)
}