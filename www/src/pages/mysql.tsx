import { HeaderComponent } from '@/components/header'
import { Layout } from '@/components/layout'
import { ResourceList } from '@/components/resource_list'
import { SideBarComponent } from '@/components/sidebar'
import dynamic from 'next/dynamic'

export function MySQLPage() {
  return (
    <Layout
      header={<HeaderComponent />}
      side={<SideBarComponent selected="mysql" />}
            main={<ResourceList rtype='mysql'/>}
    />
  )
}
export default dynamic(() => Promise.resolve(MySQLPage), { ssr: false })
