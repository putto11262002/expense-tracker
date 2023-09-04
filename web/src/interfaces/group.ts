import { IUser } from "./user"

export interface IGroup {
    id: string
    name: string
    members: IUser[]
    createdAt: Date
    updatedAt: Date
}