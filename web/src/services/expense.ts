import { CreateExpenseRequestType, ExpenseType, GetExpensesQueryType } from "../interfaces/expense";
import api from "./api";

export const createExpense = (payload: CreateExpenseRequestType) => {
    return api.post("/expense", payload).then(res => res.data as {id: string})
}

export const getExpenses = (query: GetExpensesQueryType) => {
    return api.get("/expense", {params: query}).then(res => res.data as ExpenseType[])
}