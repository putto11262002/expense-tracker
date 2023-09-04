
import { Button } from "./ui/button";

import { useAppSelector } from "../redux/store";

import { Link } from "react-router-dom";
import GroupSelector from "./GroupSelector";



function Header() {

  const { isLoggedIn, user } = useAppSelector((state) => state.auth);


  return (
    <div className="py-3 px-5 border-b flex justify-between items-center">
      <GroupSelector/>
   

      {isLoggedIn ? (
        <p>
          Logged in as{" "}
          <span className="font-bold">
            {user?.firstName} {user?.lastName}
          </span>
        </p>
      ) : (
        <Link to="/login">
          <Button>Login</Button>
        </Link>
      )}
    </div>
  );
}

export default Header;
