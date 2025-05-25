import { HeaderComponent } from '@/components/header'
import { Layout } from '@/components/layout'
import { ResourceList } from '@/components/resource_list'
import { SideBarComponent } from '@/components/sidebar'
import dynamic from 'next/dynamic'

export function PGSQLPage() {
  return (
    <Layout
      header={<HeaderComponent />}
      side={<SideBarComponent selected="sql" />}
            main={<ResourceList rtype='pgsql'/>}
    />
  )
}
export default dynamic(() => Promise.resolve(PGSQLPage), { ssr: false })
