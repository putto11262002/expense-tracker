import React from "react";
import { useQuery } from "react-query";
import { getMe } from "../services/auth";
import { useDispatch } from "react-redux";
import { init } from "../redux/authSlice";

function InitAuth({ children }: { children: React.ReactNode }) {
  const dispatch = useDispatch();
  useQuery(["user", "me"], () => getMe(), {
    onSuccess(user) {
      dispatch(init({ user: user }));
    },
    onError: () => {
      console.log("error")
      dispatch(init({user: undefined}))
    },
    retry: false
  });

  return <>{children}</>;
}

export default InitAuth;
