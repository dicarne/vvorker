import { HeaderComponent } from '@/components/header'
import { Layout } from '@/components/layout'
import { ResourceList } from '@/components/resource_list'
import { SideBarComponent } from '@/components/sidebar'
import { UserCom } from '@/components/user'
import { $user } from '@/store/userState'
import { useStore } from '@nanostores/react'
import dynamic from 'next/dynamic'

export function UserPage() {
    const user = useStore($user);
    return (
        <Layout
            header={<HeaderComponent />}
            side={<SideBarComponent selected="user" />}
            main={user?.id ? <UserCom /> : null}
        />
    )
}
export default dynamic(() => Promise.resolve(UserPage), { ssr: false })
