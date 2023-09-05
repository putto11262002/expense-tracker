import { z } from "zod"

export interface IUser {
    id: string
    email: string
    user: string
    firstName: string
    lastName: string
    createdAt: Date
    updatedAt: Date
}

export interface IRegisterRequest  {
    email: string
    firstName: string
    lastName: string
    username: string
    password: string
}



export const UserRegisterFormSchema = z.object({
    email: z.string({required_error: "Please enter your email"}).nonempty({message: "Please enter your email"}).email({message: "Invalid email"}),
    username:z.string({required_error: "Please enter your username"}).nonempty({message: "Please enter your username"}),
    firstName: z.string({required_error: "Please enter your first name"}).nonempty({message: "Please enter your first nmae"}),
    lastName: z.string({required_error: "Please enter your lastname"}).nonempty({message: "Please enter your last name"}),
    password: z.string({required_error: "Please enter your password"}),
    confirmPassword: z.string({required_error: "Please confirm your password"}).nonempty({message: "Please confirm your password"})
}).superRefine((data, ctx) => {
    if (data.confirmPassword !== data.password) {
        ctx.addIssue({
            code: "custom",
            message: "Passwords do not match",
            path: ["confirmPassword"]
        })
    }
})

export interface IUserRegisterForm extends z.infer<typeof UserRegisterFormSchema> {}