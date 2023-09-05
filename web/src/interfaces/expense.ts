/**
 * 
 * {
    "description": "Flat toilet rolls",
    "category": "Necessities",
    "groupID": "085d448e-3781-43e8-aa84-11817e3fb51f",
    "splitMode": "percentage",
    "paidBy": "f138effa-ef76-4386-ba7e-52ba60946a5c",
    "Amount": 100.0,
    "date": "2023-09-04T15:30:00Z",
    "splits": [
        {
            "userID": "f138effa-ef76-4386-ba7e-52ba60946a5c",
            "value": 10.0
        }
    ]
}
 */

import { z } from "zod";
import { formatAsCurrency } from "../utils/format";

export type CreateExpenseRequestType  = {
  description: string;
  category: string;
  groupID: string;
  paidBy: string;
  amount: number;
  date: Date;
  splits: { userID: string; value: number }[];
}

export const createExpenseFormSchema = z
  .object({
    description: z.string({ required_error: "Please enter a description" }).max(20, "Description can only be 20 characters long"),
    category: z.string({ required_error: "Please select a category" }),
    amount: z.string().transform((val, ctx) => {
      const parsed = parseFloat(val);
      if (isNaN(parsed) || !isFinite(parsed)) {
        ctx.addIssue({
          code: "custom",
          message: "Invalid amount",
          fatal: true,
        });
        return z.NEVER;
      } else {
        if (parsed < 0) {
          ctx.addIssue({
            code: "custom",
            message: "Amount cannot be negative",
            fatal: true,
          });
        }
        return parsed;
      }
    }),
    date: z.date(),
    splits: z.array(
      z.object({
        userID: z.string(),
        value: z
          .string()
          .optional()
          .transform((val, ctx) => {
            if (val === undefined || val === "") return 0;
            const parsed = parseFloat(val);
            if (isNaN(parsed)) {
              ctx.addIssue({
                code: "custom",
                message: "Invalid amount",
                fatal: true
                
              });
              return z.NEVER;
            } else {
                if (parsed < 0) {
                    ctx.addIssue({
                      code: "custom",
                      message: "Amount cannot be negative",
                      fatal: true,
                    });
                  }
              return parsed;
            }
          }),
      })
    ),
    splitSum: z.number().optional(),
    left: z.number().optional(),
  })
  .superRefine((data, ctx) => {
    data.splitSum =
      Math.round(
        data.splits?.map(({ value }) => value).reduce((x, y) => x + y) * 100
      ) / 100;
    data.left = Math.round((data.amount - data.splitSum) * 100) / 100;
    if (data.left !== 0) {
      ctx.addIssue({
        code: "custom",
        message: `${formatAsCurrency(data.left)} Left`,
        path: ["left"],
      });
    }
    return data;
  });

export type CreateExpenseFormDataType = z.infer<typeof createExpenseFormSchema>;




/**
 *  {
        "groupID": "085d448e-3781-43e8-aa84-11817e3fb51f",
        "id": "0636ec85-82c4-42bf-b982-ed2028584bcf",
        "description": "Flat toilet rolls",
        "category": "Necessities",
        "date": "2023-09-04T15:30:00Z",
        "paidBy": "f138effa-ef76-4386-ba7e-52ba60946a5c",
        "amount": 10000,
        "splits": [
            {
                "expenseID": "0636ec85-82c4-42bf-b982-ed2028584bcf",
                "value": 10000,
                "userID": "f138effa-ef76-4386-ba7e-52ba60946a5c"
            }
        ],
        "createdAt": "2023-09-03T23:17:52.711Z",
        "updatedAt": "2023-09-03T23:17:52.711Z"
    },
 */



export type GetExpensesQueryType = {
  userID?: string
  groupID?: string
  from?: string
  to?: string
}

export type  SplitType = {
  expenseID: string
  value: number
  userID: string
  settle: boolean
}

export type ExpenseType = {
  ID: string
  groupID: string
  description: string
  category: string
  date: Date
  paidBy: string
  amount: number
  splits: SplitType[]
  createdAt: Date
  updatedAt: Date
}