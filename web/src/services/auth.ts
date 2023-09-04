import axios from "axios";
import { ILoginRequest, ILoginResponse } from "../interfaces/login";
import api from "./api";

export const login = (payload: ILoginRequest) => {
  return api
    .post("/auth/login", payload)
    .then((res) => res.data as ILoginResponse)
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
