import React from 'react'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from './ui/dialog'
import { Button } from './ui/button'

function CreateExpense() {
    const [openDialog, setOpenDialog] = React.useState(false)
  return (
    <Dialog open={openDialog} onOpenChange={setOpenDialog}>
            <DialogTrigger asChild>
            <Button>Add Expense</Button>
            </DialogTrigger>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>
                        Create expense
                    </DialogTitle>
                </DialogHeader>
            </DialogContent>
    </Dialog>
  )
}

export default CreateExpense