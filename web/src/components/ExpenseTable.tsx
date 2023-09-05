import { useEffect } from "react";
import { useMutation, useQuery, useQueryClient } from "react-query";
import { getExpenses } from "../services/expense";
import { useAppSelector } from "../redux/store";
import { formatAsCurrency } from "../utils/format";
import { Button } from "./ui/button";
import { ExpenseType } from "../interfaces/expense";
import api from "../services/api";
function ExpenseTable() {
  const { dashboard, auth } = useAppSelector((state) => state);
  const queryClient = useQueryClient()

  const { selectedGroup } = dashboard;
  const { user, isLoggedIn } = auth;

  const { data, isLoading, isError, refetch } = useQuery(
    ["expense"],
    () =>
      getExpenses({
        ...(selectedGroup ? { groupID: selectedGroup.id } : {}),
        ...(user ? { userID: user.id } : {}),
      }),
    {}
  );
  const {mutate: handleSettleDept} = useMutation({
    mutationFn: (expenseId: string) => api.post("/expense/dept/settle", {expenseID: expenseId}),
  onSuccess: () => {
    queryClient.invalidateQueries({queryKey: ["expense"]})
  }
  })

  useEffect(() => {
    if (isLoggedIn) refetch();
  }, [selectedGroup, isLoggedIn]);

  if (isError) {
    return <div>Error</div>;
  }

  if (isLoading) {
    return <div>Loading</div>;
  }



  const calDepth = (expense: ExpenseType) => {
    let depth = 0
    for (const split of expense.splits){
      if (split.userID === user?.id){
        depth += split.settle ? 0 : split.value
        return depth
      }
    }
    return depth
  }

  const calLending = (expense: ExpenseType) => {
    let lending = 0
    for (const split of expense.splits){
      if (split.userID !== user?.id){
        lending += split.settle ? 0 :  split.value
      }
    }
    return lending
  }
  return (
    <div className="flex flex-col gap-2">
      {data?.map((expense) => {
        const paidBy =
          expense.paidBy === user?.id
            ? user
            : selectedGroup?.members?.find(
                (user) => user.id === expense.paidBy
              );

        
        const paidByMe = paidBy?.id === user?.id
        const lending = paidByMe ? calLending(expense) : 0
        const dept = paidByMe ? 0 : calDepth(expense)

          
        return (
          <div
            className="showdow-sm border rounded-md py-4 px-3 flex gap-5 "
            key={expense.ID}
          >
            <div className="flex flex-col">
              {expense.paidBy === user?.id ? (
                <p className="font-bold"> {"You paid"}</p>
              ) : (
                <p>
                  <span className="font-bold">{`${paidBy?.firstName} ${paidBy?.lastName}`}</span>{" "}
                  paid
                </p>
              )}
              <p className="text-sm">{formatAsCurrency(expense.amount)}</p>
            </div>

            <div className="flex flex-col">
              <p className="font-bold">Category:</p>
              <p className="text-sm">{expense.category}</p>
            </div>

            <div>
              <p className="font-bold">Description</p>
              <p className="text-sm ">{expense.description}</p>
            </div>

            {
              paidByMe ? <div>
                <p className="font-bold">Pending dept</p>
                <p className="text-sm">{formatAsCurrency(lending)}</p>
              </div> : <div>
                <p>You owe <span className="font-bold">{paidBy?.firstName} {paidBy?.lastName}</span></p>
                <p className="text-sm">{formatAsCurrency(dept)}</p>
              </div>
            }

            <div className="grow"></div>
            <div className="flex gap-3">
            <Button>View Splits</Button>
          {(!paidByMe && dept >0 ) &&  <Button onClick={() => handleSettleDept(expense.ID)}>Settle</Button>}
            </div>
          </div>
        );
      })}
    </div>
  );
}

export default ExpenseTable;
