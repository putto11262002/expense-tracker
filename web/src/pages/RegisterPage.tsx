
import { useState } from "react";
import { ArrowClockwise } from "react-bootstrap-icons";
import { Link, useNavigate } from "react-router-dom";
import { Alert } from "../components/ui/alert";
import { Input } from "../components/ui/input";
import { useForm } from "react-hook-form";
import { useMutation } from "react-query";
import {
  IUserRegisterForm,
  UserRegisterFormSchema,
} from "../interfaces/user";
import { zodResolver } from "@hookform/resolvers/zod";
import { registerUser } from "../services/auth";
import { extractErrorMessage } from "../utils/extractError";
import { Label } from "../components/ui/label";
import { Button } from "../components/ui/button";

function RegisterPage() {
  const navigate = useNavigate()
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<IUserRegisterForm>({
    resolver: zodResolver(UserRegisterFormSchema),
  });

  const [error, setError] = useState<string | undefined>(undefined);

  const { mutate: handleRegister, isLoading } = useMutation({
    mutationFn: (formData: IUserRegisterForm) => {
      return registerUser({
        ...formData,
      });
    },
    onSuccess: () => {
        setError(undefined)
        navigate("/login")
    },
    onError: async (err) => {
        const errMessage = await extractErrorMessage(err)
        console.log(errMessage)
        setError(errMessage)
    }
  });



  
  return (
    <div className="h-screen w-screen flex justify-center items-center">
      <form
        onSubmit={handleSubmit((formData) => handleRegister(formData))}
        className="w-full max-w-[400px] grid grid-cols-2 gap-4"
      >
        <h2 className="text-center text-xl font-bold col-span-2">
          Welcome to <br /> Expense tracker
        </h2>

        {typeof error === "string" && (
          <Alert variant="destructive" className="col-span-2">
            <p className="text-center text-sm">{error}</p>
          </Alert>
        )}
        <div className="col-span-2 flex flex-col">
          <Label className="mb-3">Username</Label>
          <Input {...register("username", {})} />
          <p className="text-xs mt-1">{errors["username"]?.message?.toString()}</p>
        </div>

        <div className="col-span-2 flex flex-col">
          <Label className="mb-3">Email</Label>
          <Input {...register("email")} type="text" />
          <p className="text-xs mt-1">{errors["email"]?.message?.toString()}</p>
        </div>

        <div className="col-span-1 flex flex-col ">
          <Label className="mb-3">First Name</Label>
          <Input {...register("firstName")} type="text" />
          <p className="text-xs mt-1">{errors["firstName"]?.message?.toString()}</p>
        </div>

        <div className="col-span-1 flex flex-col ">
          <Label className="mb-3">Last Name</Label>
          <Input {...register("lastName")} type="text" />
          <p className="text-xs mt-1">{errors["lastName"]?.message?.toString()}</p>
        </div>

        <div className="col-span-2 flex flex-col ">
          <Label className="mb-3">Password</Label>
          <Input {...register("password")} type="password" />
          <p className="text-xs mt-1">{errors["password"]?.message?.toString()}</p>
        </div>

        <div className="col-span-2 flex flex-col ">
          <Label className="mb-3">Confirm Password</Label>
          <Input {...register("confirmPassword")} type="password" />
          <p className="text-xs mt-1">{errors["confirmPassword"]?.message?.toString()}</p>
        </div>

        <div className="col-span-2 text-center pt-3">
          <Button
            type="submit"
            disabled={isLoading}
            className="margin-x-auto px-5"
          >
            {isLoading && (
              <ArrowClockwise className="animate-spin mr-2 h-4 w-4" />
            )}{" "}
            Register
          </Button>
        </div>

        <p className="col-span-2 text-center">
          Already have an account?{" "}
          <Link className="underline" to={"/login"}>
            Login
          </Link>
        </p>
      </form>
    </div>
  );
}

export default RegisterPage;
