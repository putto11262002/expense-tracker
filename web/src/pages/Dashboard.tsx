import AddMember from '../components/AddMember'
import CreateExpense from '../components/CreateExpense'
import ExpenseTable from '../components/ExpenseTable'
import { useAppSelector } from '../redux/store'


function Dashboard() {

  const {selectedGroup} = useAppSelector(state => state.dashboard)

 
  return (
   <>
   <div className='pb-4 flex gap-2'>
    <CreateExpense/>
    <AddMember/>

   
    
   </div>

   {selectedGroup && <ExpenseTable/>}
   
   </>


  )
}

export default Dashboard