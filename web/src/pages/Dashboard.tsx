import AddMember from '../components/AddMember'
import CreateExpense from '../components/CreateExpense'
import ExpenseTable from '../components/ExpenseTable'


function Dashboard() {

 
  return (
   <>
   <div className='pb-4 flex gap-2'>
    <CreateExpense/>
    <AddMember/>

   
    
   </div>

   <ExpenseTable/>
   </>


  )
}

export default Dashboard