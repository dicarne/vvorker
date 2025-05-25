import { HeaderComponent } from '@/components/header'
import { Layout } from '@/components/layout'
import { ResourceList } from '@/components/resource_list'
import { SideBarComponent } from '@/components/sidebar'
import dynamic from 'next/dynamic'

export function OSSPage() {
  return (
    <Layout
      header={<HeaderComponent />}
      side={<SideBarComponent selected="oss" />}
            main={<ResourceList rtype='oss'/>}
    />
  )
}
export default dynamic(() => Promise.resolve(OSSPage), { ssr: false })
