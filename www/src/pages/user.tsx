import { HeaderComponent } from '@/components/header'
import { Layout } from '@/components/layout'
import { ResourceList } from '@/components/resource_list'
import { SideBarComponent } from '@/components/sidebar'
import { UserCom } from '@/components/user'
import dynamic from 'next/dynamic'

export function UserPage() {
    return (
        <Layout
            header={<HeaderComponent />}
            side={<SideBarComponent selected="user" />}
            main={ <UserCom />}
        />
    )
}
export default dynamic(() => Promise.resolve(UserPage), { ssr: false })
