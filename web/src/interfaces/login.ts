import { IUser } from "./user"

export interface ILoginRequest {
    key: string
    secret: string
}

export interface ILoginResponse {
    user: IUser
    token: string
}