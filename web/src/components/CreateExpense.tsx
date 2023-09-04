import React, { useEffect, useRef, useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "./ui/dialog";
import { Button } from "./ui/button";
import { Label } from "./ui/label";
import { Input } from "./ui/input";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "./ui/select";
import { Calendar } from "./ui/calendar";
import {
  Popover,
  PopoverTrigger,
  PopoverContent,
} from "@radix-ui/react-popover";
import { format } from "date-fns";
import { cn } from "../lib/utils";
import { Controller, useForm } from "react-hook-form";
import {
  CreateExpenseFormDataType,
  CreateExpenseRequestType,
  createExpenseFormSchema,
} from "../interfaces/expense";
import { CalendarIcon } from "lucide-react";
import { useAppSelector } from "../redux/store";

import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useMutation } from "react-query";
import { createExpense } from "../services/expense";
import { queryClient } from "../main";

function formatAsCurrency(number: number, currencyCode = "USD") {
  return number.toLocaleString("en-US", {
    style: "currency",
    currency: currencyCode,
  });
}

function CreateExpense() {
  const [openDialog, setOpenDialog] = React.useState(false);
  const {
    watch,
    register,
    control,
    handleSubmit,
    formState: { errors },
    setValue,
    reset,
    getValues,
  } = useForm<CreateExpenseFormDataType>({
    reValidateMode: "onChange",
    resolver: zodResolver(createExpenseFormSchema),
  });

  const formBtnRef = useRef<HTMLButtonElement>(null);
  const date = watch("date");
  const { auth, dashboard } = useAppSelector((state) => state);
  const {isLoggedIn, user} = auth
  const {selectedGroup} = dashboard
  const [splitMode, setSplitMode] = useState("amount");
  const amount = watch("amount");
  const [left, setLeft] = useState(0);

  const {mutate: handleCreateExpense, isLoading: isCreatingExpense} = useMutation({
    mutationFn: (formData: CreateExpenseFormDataType) => {
      if (!selectedGroup || !user) throw new Error("no group or user selected")
      const payload: CreateExpenseRequestType = {
        description: formData.description,
        category: formData.category,
        groupID: selectedGroup?.id,
        paidBy: user?.id,
        amount: formData.amount,
        date: formData.date,
        splits: formData.splits

      }

      return createExpense(payload)
    },
    onSuccess: () => {
      queryClient.invalidateQueries(["expense"])
      setOpenDialog(false)
    }
  })



  useEffect(() => {
    if (openDialog) {
      reset({});
      setLeft(0);
    }
  }, [openDialog, reset]);

  const computeLeft = () => {
    const splitSum =
      Math.round(
        (getValues("splits")
          ?.filter(({ value }) => !isNaN(parseInt(String(value))))
          .map(({ value }) => parseFloat(String(value)))
          .reduce((x, y) => x + y) || 0) * 100
      ) / 100;

      const tempLeft = Math.round(((getValues("amount") || 0) - splitSum) * 100) / 100
    setLeft(tempLeft);
  };



  return (
    <Dialog open={openDialog} onOpenChange={setOpenDialog}>
      <DialogTrigger asChild>
        <Button disabled={!isLoggedIn || !selectedGroup}> Add Expense</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create expense</DialogTitle>
        </DialogHeader>
        <form
          onSubmit={handleSubmit((formData) => handleCreateExpense(formData))}
          className="grid grid-cols-2 gap-3"
        >
          <div className="space-y-2 col-span-1">
            <Label>Amount ($)</Label>
            <Controller
              control={control}
              name="amount"
              render={({ field: { value, onChange } }) => (
                <Input
                  value={value || ""}
                  onChange={(e) => {
                    onChange(e);

                    computeLeft();
                  }}
                />
              )}
            />
            {errors["amount"] && (
              <p className="text-xs">{errors["amount"]?.message?.toString()}</p>
            )}
          </div>

          <div className="space-y-2 col-span-1">
            <Label>Category</Label>
            <Controller
              control={control}
              name="category"
              render={({ field: { value, onChange } }) => (
                <Select value={value} onValueChange={onChange}>
                  <SelectTrigger className="">
                    <SelectValue className="" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="apple">Apple</SelectItem>
                    <SelectItem value="banana">Banana</SelectItem>
                    <SelectItem value="blueberry">Blueberry</SelectItem>
                    <SelectItem value="grapes">Grapes</SelectItem>
                    <SelectItem value="pineapple">Pineapple</SelectItem>
                  </SelectContent>
                </Select>
              )}
            />
            {errors["category"] && (
              <p className="text-xs">
                {errors["category"]?.message?.toString()}
              </p>
            )}
          </div>

          <div className="space-y-2 col-span-2">
            <Label>Description</Label>
            <Input {...register("description", {})} />
            {errors["description"] && (
              <p className="text-xs">
                {errors["description"]?.message?.toString()}
              </p>
            )}
          </div>
          <div className="space-y-2 col-span-1">
            <Label>Date</Label>
            <Controller
              control={control}
              name="date"
              render={({ field: { onChange, value } }) => (
                <>
                  <div>
                    <Popover>
                      <PopoverTrigger asChild>
                        <Button
                          variant={"outline"}
                          className={cn(
                            "w-full justify-start text-left font-normal",
                            !date && "text-muted-foreground"
                          )}
                        >
                          <CalendarIcon className="mr-2 h-4 w-4" />
                          {date ? (
                            format(date, "PPP")
                          ) : (
                            <span>Pick a date</span>
                          )}
                        </Button>
                      </PopoverTrigger>
                      <PopoverContent className="w-auto p-0" align="start">
                        <Calendar
                          mode="single"
                          className="bg-white shadow-md mb-2 rounded-sm"
                          selected={value}
                          onSelect={onChange}
                          initialFocus
                        />
                      </PopoverContent>
                    </Popover>
                  </div>
                </>
              )}
            />
            {errors["date"] && (
              <p className="text-xs">{errors["date"]?.message?.toString()}</p>
            )}
          </div>

          <div className="space-y-2 col-span-1">
            <Label>Paid by</Label>
            <Input disabled value={`${user?.firstName} ${user?.lastName}`} />
          </div>

          {amount > 0 && (
            <div className="col-span-2 space-y-2">
              <Label>Split by</Label>
              <Tabs defaultValue="amount" className="w-full">
                <div className="flex justify-between items-center">
                  <TabsList>
                    <TabsTrigger
                      onClick={() => setSplitMode("amount")}
                      value="amount"
                    >
                      Amount
                    </TabsTrigger>
                    <TabsTrigger
                      disabled
                      onClick={() => setSplitMode("equally")}
                      value="equally"
                    >
                      Equally
                    </TabsTrigger>
                  </TabsList>
                  <p className={cn("text-red-600", left === 0  && "text-green-600")}>{formatAsCurrency(left)}</p>
                </div>
                <TabsContent
                  className="max-h-[150px] overflow-y-scroll"
                  value="amount"
                >
                  <div className="flex flex-col px-1 gap-2">
                    {selectedGroup?.members?.map((member, idx) => {
                      setValue(`splits.${idx}.userID`, member.id);

                      return (
                        <div
                          key={member.id}
                          className="flex justify-between items-center border rounded-md py-3 px-4 "
                        >
                          <p className="">{`${member.firstName} ${member.lastName}`}</p>
                          <div>
                            <Controller
                              control={control}
                              name={`splits.${idx}.value`}
                              
                              render={({
                                field: { value, onChange },
                                fieldState: { error },
                              }) => (
                                <div className="space-y-2">
                                  <Input
                                    value={value || ""}
                                    onChange={(e) => {
                                      onChange(e);
                                      computeLeft();
                                    }}
                                  />
                                  {error && (
                                    <p className="text-xs">
                                      {error?.message?.toString()}
                                    </p>
                                  )}
                                </div>
                              )}
                            />
                          </div>
                        </div>
                      );
                    })}
                  </div>
                </TabsContent>
                <TabsContent value="equally">Equally</TabsContent>
              </Tabs>
            </div>
          )}
          <button type="submit" ref={formBtnRef} hidden></button>
        </form>
        <DialogFooter>
          <Button disabled={isCreatingExpense} onClick={() => formBtnRef?.current?.click()}>Create</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

export default CreateExpense;
