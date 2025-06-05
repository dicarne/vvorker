import { HeaderComponent } from '@/components/header'
import { Layout } from '@/components/layout'
import { LogsComponent } from '@/components/logs'
import { SideBarComponent } from '@/components/sidebar'
import dynamic from 'next/dynamic'

export function WorkerPage() {
    return (
        <Layout
            header={<HeaderComponent />}
            side={<SideBarComponent selected="task" />}
            main={<LogsComponent />}
        />
    )
}
export default dynamic(() => Promise.resolve(WorkerPage), { ssr: false })
