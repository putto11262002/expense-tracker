import { useEffect } from "react";
import { useQuery } from "react-query";
import { getExpenses } from "../services/expense";
import dayjs from "dayjs";
import { useAppSelector } from "../redux/store";
function ExpenseTable() {
    const {dashboard, auth} = useAppSelector(state => state)

    const {selectedGroup} = dashboard
    const {user, isLoggedIn} = auth

  const { data, isLoading, isError , refetch} = useQuery(
    ["expense"],
    () => getExpenses({...(selectedGroup ? {groupID: selectedGroup.id} : {}), ...(user ? {userID: user.id} : {})}),
    {
        enabled: false
    }
  );

  useEffect(() => {
    if (isLoggedIn) refetch()
  }, [selectedGroup, isLoggedIn])

  if (isError) {
    return <div>Error</div>;
  }

  if (isLoading) {
    return <div>Loading</div>;
  }

  return (
    <div className="flex flex-col gap-2">
      
      {data?.map((expense) => (
        <div
          className="showdow-sm border rounded-md py-4 px-3 flex "
          key={expense.ID}
        >
          <div>
            <p className="font-bold">{expense.description}</p>
            <p className="text-sm text-slate-400">
              {dayjs(expense.date).format("YYYY MMM DD")}
            </p>
          </div>
          <div></div>
        </div>
      ))}
    </div>
  );
}

export default ExpenseTable;
