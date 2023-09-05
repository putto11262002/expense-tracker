import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "./ui/dialog";

import { useAppSelector } from "../redux/store";
import { Button } from "./ui/button";
import { useMutation, useQuery, useQueryClient } from "react-query";
import api from "../services/api";
import { Input } from "./ui/input";
import { useForm } from "react-hook-form";
import { IUser } from "../interfaces/user";


function AddMember() {
  const { selectedGroup } = useAppSelector((state) => state.dashboard);
  const [openDialog, setOpenDialog] = useState(false);
  const { register, handleSubmit, watch } = useForm();
  const [userToAdd, setUserToAdd] = useState<undefined | IUser>(undefined)
  const email = watch("email");
  const queryClient = useQueryClient()

  const {  refetch } = useQuery(
    ["user"],
    () => {
      return api
        .get("/user", { params: { email, notInGroup: selectedGroup?.id } })
        .then((res) => (res.data || []) as IUser[]);
    },
    {
      enabled: false,
      onSuccess(data) {
          setUserToAdd(data[0])
      },
    }
  );


  const {mutate: handleAddUser} = useMutation({
    mutationFn: () => {
        return api.post(`/group/user/add`, {groupID: selectedGroup?.id, userID: userToAdd?.id})
    },
    onSuccess: () => {
        queryClient.invalidateQueries("group")
        setUserToAdd(undefined)
    },
    onError: (err) => {
        console.error(err)
    }
  })

  return (
    <Dialog open={openDialog} onOpenChange={setOpenDialog}>
      <DialogTrigger asChild>
        <Button disabled={!selectedGroup}>Add Member</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add Member</DialogTitle>
        </DialogHeader>
        <div className="space-y-4">
          <form onSubmit={handleSubmit(() => refetch())} className="flex gap-2">
            <Input {...register("email")} placeholder="email..." />
            <Button>Search</Button>
          </form>

          {userToAdd  && (
            <div className="flex items-center border rounded-md  py-2 px-4">
              <p className="grow font-bold">{`${userToAdd.firstName} ${userToAdd.lastName}`}</p>
              <Button onClick={() => handleAddUser()} className="" variant="secondary">
                Add
              </Button>
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}

export default AddMember;
