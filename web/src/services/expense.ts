import { CreateExpenseRequestType } from "../interfaces/expense";
import api from "./api";

export const createExpense = (payload: CreateExpenseRequestType) => {
    return api.post("/expense", payload).then(res => res.data as {id: string})
}