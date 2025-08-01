import { Button, Nav } from '@douyinfe/semi-ui'
import {
  IconWrench,
  IconPercentage,
  IconHourglass,
  IconSidebar,
  IconMenu,
  IconArchive,
  IconLayers,
  IconServer,
  IconUser,
} from '@douyinfe/semi-icons'
import { useRouter } from 'next/router'
import { atom } from 'nanostores'
import { useStore } from '@nanostores/react'
import cn from 'clsx'
import { useEffect, useState } from 'react'
import { $user } from '@/store/userState'

export const $expandSidebar = atom(false)

export const SideBarComponent = ({ selected }: { selected: string }) => {
  const router = useRouter()
  const expand = useStore($expandSidebar)
  const [invisible, setInvisible] = useState(expand)
  const user = useStore($user);
  const routeMap = {
    workers: '/admin',
    status: '/nodes',
    settings: '/admin',
    task: '/task',
    oss: '/oss',
    kv: '/kv',
    sql: '/sql',
    mysql: '/mysql',
    user: '/user',
    users: '/users'
  } as any
  useEffect(() => {
    if (!expand) {
      setTimeout(() => setInvisible(!expand), 150)
    } else {
      setInvisible(!expand)
    }
  }, [expand])
  return (
    <>
      <Nav
        className={cn('z-10 fixed md:relative md:visible', { invisible })}
        style={{ height: '93vh' }}
        items={[
          { itemKey: 'workers', text: 'Workers', icon: <IconPercentage /> },
          { itemKey: 'status', text: 'Status', icon: <IconHourglass /> },
          { itemKey: 'task', text: 'Task', icon: <IconMenu /> },
          { itemKey: 'sql', text: 'SQL', icon: <IconServer /> },
          { itemKey: 'mysql', text: 'MySQL', icon: <IconServer /> },
          { itemKey: 'oss', text: 'OSS', icon: <IconArchive /> },
          { itemKey: 'kv', text: 'KV', icon: <IconLayers /> },
          { itemKey: 'user', text: 'User', icon: <IconUser /> },
          { itemKey: 'users', text: 'Users Management', icon: <IconUser />, hide: user?.role != 'admin' },
        ].filter(a => !a.hide)}
        onSelect={(data) => console.log('trigger onSelect: ', data)}
        onClick={(data) => {
          window.location.assign(routeMap[data.itemKey || ''])
        }}
        footer={{ collapseButton: true }}
        isCollapsed={!expand}
        onCollapseChange={() => $expandSidebar.set(!expand)}
        defaultSelectedKeys={[selected]}
      />
    </>
  )
}