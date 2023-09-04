import { Input } from "../components/ui/input";
import { Label } from "../components/ui/label";
import { Button } from "../components/ui/button";
import { useForm } from "react-hook-form";
import { ILoginRequest } from "../interfaces/login";
import { useMutation } from "react-query";
import { login } from "../services/auth";
import { ArrowClockwise } from "react-bootstrap-icons";
import { Alert } from "../components/ui/alert";
import { useNavigate } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "../redux/store";
import { useEffect } from "react";
import { init } from "../redux/authSlice";

function LoginPage() {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<ILoginRequest>();

  const navigate = useNavigate();

  const {isLoggedIn} = useAppSelector(state => state.auth)
  const dispatch = useAppDispatch()

  useEffect(() => {
    if (isLoggedIn){
        navigate("/")
    }
  }, [isLoggedIn, navigate])


  const {mutate: handleLogin, isLoading, error} = useMutation({
    mutationFn: (paylaod: ILoginRequest) => login(paylaod),
   
    onSuccess: (data) => {
        // setting client user state
        dispatch(init({user: data.user}))

        // redirecting to main page
        navigate("/")
        
    }
  })

 

  return (
    <div className="h-screen w-screen flex justify-center items-center">
      <form
        onSubmit={handleSubmit((formData) => handleLogin(formData))}
        className="w-full max-w-[300px] flex flex-col justify-center gap-4"
      >
        <h2 className="text-center text-xl font-bold">
          Welcome to <br /> Expense tracker
        </h2>

       {typeof error === "string" && <Alert variant="destructive">
            <p className="text-center text-sm">{error}</p>

        </Alert>}
        <div className="flex flex-col gap-2">
          <Label>Email</Label>
          <Input
            {...register("key", {
              required: {
                value: true,
                message: "Please enter your email address",
              },
            })}
          />
          <p className="text-xs">{errors["key"]?.message?.toString()}</p>
        </div>

        <div className="flex flex-col gap-2">
          <Label>Password</Label>
          <Input
            {...register("secret", {
              required: { value: true, message: "Please enter your password" },
            })}
            type="password"
          />
          <p className="text-xs">{errors["secret"]?.message?.toString()}</p>
        </div>

        <div className="text-center pt-3">
          <Button type="submit" disabled={isLoading} className="margin-x-auto px-5">{isLoading && <ArrowClockwise className="animate-spin mr-2 h-4 w-4"/>} Login</Button>
        </div>
      </form>
    </div>
  );
}

export default LoginPage;
