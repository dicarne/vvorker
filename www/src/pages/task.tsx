import { HeaderComponent } from '@/components/header'
import { Layout } from '@/components/layout'
import { SideBarComponent } from '@/components/sidebar'
import { TaskComponent } from '@/components/task'
import dynamic from 'next/dynamic'

export function TaskPage() {
  return (
    <Layout
      header={<HeaderComponent />}
      side={<SideBarComponent selected="task" />}
      main={<TaskComponent />}
    />
  )
}
export default dynamic(() => Promise.resolve(TaskPage), { ssr: false })
