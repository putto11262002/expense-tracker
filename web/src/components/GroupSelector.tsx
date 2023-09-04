import React, { useRef } from "react";
import {
  Command,
  CommandGroup,
  CommandItem,
  CommandList,
  CommandSeparator,
} from "./ui/command";
import { init } from "../redux/authSlice";
import { Popover, PopoverContent, PopoverTrigger } from "./ui/popover";
import { Button } from "./ui/button";
import { ChevronDoubleDown } from "react-bootstrap-icons";
import { CreateGroup, getMyGroups } from "../services/group";
import { useMutation, useQuery } from "react-query";
import { useAppDispatch, useAppSelector } from "../redux/store";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "./ui/dialog";
import { Label } from "./ui/label";
import { Input } from "./ui/input";
import { useForm } from "react-hook-form";
import { queryClient } from "../main";
import { Alert } from "./ui/alert";
import { setSelectGroup } from "../redux/dashboardSlice";

function GroupSelector() {
  const [open, setOpen] = React.useState(false);
  const [showNewTeamDialog, setShowNewTeamDialog] = React.useState(false);
  const { auth, dashboard } = useAppSelector((state) => state);
  const { isLoggedIn, user } = auth;
  const { selectedGroup } = dashboard;

  const formBtnRef = useRef<HTMLButtonElement>(null);

  const dispatch = useAppDispatch();

  const {
    register,
    formState: { errors },
    handleSubmit,
  } = useForm<{ name: string }>();

  const { data: groups, isLoading: isLoadingGroups } = useQuery(
    ["group"],
    () => getMyGroups(),
    {
      enabled: !!isLoggedIn,
    }
  );

  const {
    mutate: handleCreateGroup,
    isError: isErrorCreatingGroup,
    isLoading: isCreatingGroup,
  } = useMutation({
    mutationFn: (payload: { name: string }) => CreateGroup(payload),
    onSuccess: () => {
      queryClient.invalidateQueries("group");
      setShowNewTeamDialog(false);
    },
  });
  return (
    <Dialog
      open={showNewTeamDialog && isLoggedIn}
      onOpenChange={setShowNewTeamDialog}
    >
      <Popover open={open && isLoggedIn} onOpenChange={setOpen}>
        <PopoverTrigger asChild>
          <Button
            disabled={!init || !isLoggedIn || isLoadingGroups}
            className="justify-between w-[200px]"
            variant="outline"
          >
            { selectedGroup ? selectedGroup.name  : 'Select Group...' }
            <ChevronDoubleDown className="h-4 w-4" />
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-[200px] p-0">
          <Command>
            <CommandList>
              <CommandGroup>
                {groups?.map((group) => (
                  <CommandItem
                    onSelect={() =>{ 
                      dispatch(setSelectGroup({ group }))
                      setOpen(false)
                    }}
                    key={group.id}
                  >
                    {group.name}
                  </CommandItem>
                ))}
              </CommandGroup>
            </CommandList>
            <CommandSeparator />
            <CommandList>
              <CommandGroup>
                <CommandItem
                  onSelect={() => {
                    setOpen(false);
                    setShowNewTeamDialog(true);
                  }}
                >
                  Create group
                </CommandItem>
              </CommandGroup>
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>

      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create group</DialogTitle>
        </DialogHeader>
        <form
          onSubmit={handleSubmit((formData) => handleCreateGroup(formData))}
          className="flex flex-col gap-2"
        >
          {isErrorCreatingGroup && (
            <Alert variant="destructive">
              <p className="text-center">
                {" "}
                Something went wrong please try again later
              </p>
            </Alert>
          )}
          <div className="space-y-2">
            <Label>Group name</Label>
            <Input
              {...register("name", {
                required: { value: true, message: "Please enter group name" },
              })}
            />
            {errors["name"]?.message && (
              <p className="text-xs">{errors["name"].message.toString()}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label>Group name</Label>
            <Input disabled value={`${user?.firstName} ${user?.lastName}`} />
          </div>
          <button type="submit" hidden ref={formBtnRef}></button>
        </form>

        <DialogFooter>
          <Button
            disabled={isCreatingGroup}
            onClick={() => formBtnRef.current?.click()}
          >
            Create
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

export default GroupSelector;
