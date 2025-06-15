import React, { useEffect } from 'react';
import { UsersManagementCom } from '@/components/usersManagement';
import { $user } from '@/store/userState';
import { useStore } from '@nanostores/react';
import { useRouter } from 'next/router';
import { HeaderComponent } from '@/components/header'
import { Layout } from '@/components/layout'
import { SideBarComponent } from '@/components/sidebar'
import { WorkerEditComponent } from '@/components/worker_edit'
import dynamic from 'next/dynamic'

const UsersPage: React.FC = () => {
    const router = useRouter();
    const user = useStore($user);
    const [isAdmin, setIsAdmin] = React.useState<boolean>(false);

    // 检查用户是否为管理员，如果不是则重定向到首页
    useEffect(() => {
        if (user && user?.userName && user?.role !== 'admin') {
            router.replace('/admin');
        }
        if (user && user?.userName && user?.role === 'admin') {
            setIsAdmin(true);
        }
    }, [user, router]);

    return <Layout
        header={<HeaderComponent />}
        side={isAdmin ? <SideBarComponent selected="workers" /> : null}
        main={< UsersManagementCom />}
    />
};

export default dynamic(() => Promise.resolve(UsersPage), { ssr: false })
