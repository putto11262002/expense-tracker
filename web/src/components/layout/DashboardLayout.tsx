import { Outlet } from "react-router-dom";
import Header from "../Header";


function DashboardLayout() {
   
  return (
    <div className="flex flex-col">
   <div className="grow-0">
   <Header/>
   </div>
   <div className="py-2 pt-4 px-3 grow">
     
      <Outlet />
    </div>
    </div>
  );
}

export default DashboardLayout;
