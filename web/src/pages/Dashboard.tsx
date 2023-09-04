import CreateExpense from '../components/CreateExpense'
import ExpenseTable from '../components/ExpenseTable'


function Dashboard() {

 
  return (
   <>
   <div className='pb-4'>
    <CreateExpense/>

   
    
   </div>

   <ExpenseTable/>
   </>


  )
}

export default Dashboard