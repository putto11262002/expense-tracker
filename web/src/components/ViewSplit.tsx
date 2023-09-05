import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "./ui/dialog";

import { ExpenseType } from "../interfaces/expense";
import { formatAsCurrency } from "../utils/format";
import { Button } from "./ui/button";
import { useAppSelector } from "../redux/store";
import { cn } from "../lib/utils";

function ViewSplit({ expense }: { expense: ExpenseType }) {
  const [openDialog, setOpenDialog] = useState(false);
  const { dashboard, auth } = useAppSelector((state) => state);

  const { selectedGroup } = dashboard;
  const { user } = auth;
  const getGroupMember = (id: string) =>
    selectedGroup?.members.find((user) => user.id === id);
  return (
    <Dialog open={openDialog} onOpenChange={setOpenDialog}>
      <DialogTrigger asChild>
        <Button>View Split</Button>
      </DialogTrigger>

      <DialogContent>
        <DialogHeader>
          <DialogTitle>View split</DialogTitle>
        </DialogHeader>

        <div className="flex flex-col gap-3">
          {expense.splits.map((split, idx) => {
            const member = getGroupMember(split.userID);
            return (
              <div
                className="flex border rounded-md py-3 px-4 gap-2 items-center"
                key={idx}
              >
                <p
                  className={cn(
                    "py-1 px-2 rounded-md",
                    "bg-yellow-400/50",
                    split.settle && "bg-green-600/50"
                  )}
                >
                  {split.settle ? "Settled" : "Pending"}
                </p>
                <p className="grow">
                  {member?.id === user?.id
                    ? "You"
                    : `${member?.firstName} ${member?.lastName}`}
                </p>{" "}
                <p>{formatAsCurrency(split.value)}</p>
              </div>
            );
          })}
        </div>
      </DialogContent>
    </Dialog>
  );
}

export default ViewSplit;
