import { HeaderComponent } from '@/components/header'
import { Layout } from '@/components/layout'
import { SideBarComponent } from '@/components/sidebar'
import { ResourceList } from '@/components/resource_list'
import dynamic from 'next/dynamic'

export function KVPage() {
  return (
    <Layout
      header={<HeaderComponent />}
      side={<SideBarComponent selected="kv" />}
      main={<ResourceList rtype='kv'/>}
    />
  )
}
export default dynamic(() => Promise.resolve(KVPage), { ssr: false })
