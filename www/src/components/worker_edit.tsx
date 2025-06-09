import {
  Breadcrumb,
  Button,
  ButtonGroup,
  Divider,
  Input,
  Notification,
  Select,
  TabPane,
  Tabs,
  Toast,
  Typography,
  Pagination,
  List,
  Tag
} from '@douyinfe/semi-ui'
import { DEFAUTL_WORKER_ITEM, WorkerItem } from '@/types/workers'
import * as api from '@/api/workers'
import { useRouter } from 'next/router'
import { useMutation, useQuery } from '@tanstack/react-query'
import { $code, $vorkerSettings } from '@/store/workers'
import { useStore } from '@nanostores/react'
import { IconArticle, IconHome } from '@douyinfe/semi-icons'
import { getNodes } from '@/api/nodes'
import dynamic from 'next/dynamic'
import { i18n } from '@/lib/i18n'
import { TemplateEditor } from './template_editor'
import type { TaskLog, WorkerLog } from '@/types/workers';
import { useEffect, useState, useCallback, useRef } from 'react'

const MonacoEditor = dynamic(
  import('./editor').then((m) => m.MonacoEditor),
  { ssr: false }
)

// 定义时间格式化函数
const formatDate = (date: Date) => {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  const seconds = String(date.getSeconds()).padStart(2, '0');
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
};

export const WorkerEditComponent = () => {
  const router = useRouter()
  const { UID } = router.query
  const [editItem, setEditItem] = useState(DEFAUTL_WORKER_ITEM)
  const [templateContent, setTemplateContent] = useState('')
  const [logs, setLogs] = useState<WorkerLog[]>([]);
  const [loadingLogs, setLoadingLogs] = useState(false);
  const [errorLogs, setErrorLogs] = useState<string | null>(null);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(50);
  const [totalLogs, setTotalLogs] = useState(0);

  const appConf = useStore($vorkerSettings)
  const code = useStore($code)
  const { Paragraph, Text, Numeral, Title } = Typography
  const { data: resp } = useQuery(['getNodes'], () => getNodes())

  const { data: worker } = useQuery(['getWorker', UID], () => {
    return UID ? api.getWorker(UID as string) : null
  })

  const updateWorker = useMutation(async () => {
    await api.updateWorker(UID as string, editItem)
    Toast.info(i18n('workerSaveSuccess'))
  })

  const runWorker = useMutation(async (UID: string) => {
    let resp = await api.runWorker(UID)
    let raw_resp = JSON.stringify(resp)
    let run_resp = Buffer.from(resp?.data?.run_resp, 'base64').toString('utf8')
    let opts = {
      title: 'worker run result',
      content: (
        <>
          <Paragraph spacing="extended">
            <code className="overflow-scroll w-full">
              {run_resp.length > 100
                ? run_resp.slice(0, 100) + '......'
                : run_resp.length == 0
                  ? 'data is undefined, raw resp: ' + raw_resp
                  : run_resp}
            </code>
          </Paragraph>
          <div className="flex flex-row justify-end">
            <Text>copy to see full content</Text>
            <Paragraph
              copyable={{ content: run_resp }}
              spacing="extended"
              className="justify-end"
            />
          </div>
        </>
      ),
      duration: 10,
    }
    Notification.info({ ...opts, position: 'bottomRight' })
  })

  useEffect(() => {
    worker && setEditItem(worker)
  }, [UID, worker])

  useEffect(() => {
    if (worker) {
      setEditItem(worker)
      $code.set(Buffer.from(worker.Code, 'base64').toString('utf8'))
      if (worker.Template) setTemplateContent(worker.Template)
      else { setTemplateContent(DEFAUTL_WORKER_ITEM.Template) }
    }
  }, [worker])

  useEffect(() => {
    if (code && editItem)
      setEditItem((item) => ({
        ...item,
        Template: templateContent,
        Code: Buffer.from(code).toString('base64'),
      }))
  }, [code, editItem, templateContent])

  useEffect(() => {
    worker?.Code
  })

  const [activeTab, setActiveTab] = useState('code'); // 新增：记录当前激活的标签页
  const intervalRef = useRef<number | null>(null); // 新增：用于存储定时器 ID

  const fetchLogs = useCallback(async (_page: number, _page_size: number) => {
    if (!UID) return;
    setLoadingLogs(true);
    try {
      const data = (await api.getWorkerLogs(UID as string, _page, _page_size)).data.data;
      setLogs(data.logs);
      setTotalLogs(data.total);
    }
    catch (err) {
      setErrorLogs('Failed to fetch logs');
    }
    finally {
      setLoadingLogs(false);
    }
  }, [UID]);

  // 新增：当标签页变化时，处理定时器
  useEffect(() => {
    if (activeTab === 'logs') {
      // 每秒刷新日志
      intervalRef.current = window.setInterval(() => {
        fetchLogs(page, pageSize);
      }, 1000);
    } else {
      // 清除定时器
      if (intervalRef.current) {
        window.clearInterval(intervalRef.current);
      }
    }
    // 组件卸载时清除定时器
    return () => {
      if (intervalRef.current) {
        window.clearInterval(intervalRef.current);
      }
    };
  }, [activeTab, page, pageSize, fetchLogs]);

  useEffect(() => {
    fetchLogs(page, pageSize);
  }, [page, pageSize, fetchLogs]);

  const workerURL = `${appConf?.Scheme}://${editItem.Name}${appConf?.WorkerURLSuffix}`

  return (
    <div className="m-4 flex flex-col">
      <div className="flex flex-row justify-between">
        <div className="flex flex-col">
          <Breadcrumb compact={false}>
            <Breadcrumb.Item
              href="/admin"
              icon={<IconHome size="small" />}
            ></Breadcrumb.Item>
            <Breadcrumb.Item href="/admin">Workers</Breadcrumb.Item>
            <Breadcrumb.Item href={`/worker?UID=${editItem.UID}`}>
              {editItem.Name}
            </Breadcrumb.Item>
          </Breadcrumb>
        </div>
        <div className="flex flex-col">
          <ButtonGroup>
            <Button
              onClick={() => {
                window.location.assign('/admin')
              }}
            >
              Back
            </Button>
            <Button onClick={() => updateWorker.mutate()}>Save</Button>
          </ButtonGroup>
        </div>
      </div>
      <div className="flex flex-row gap-1">
        <div className="columns-1 md:columns-2">
          <div></div>
          <Title heading={5}>ID</Title>
          <Paragraph copyable={{ content: editItem.UID }} spacing="extended">
            <code>{editItem.UID}</code>
          </Paragraph>
          <Title heading={5}>URL</Title>
          <Paragraph copyable={{ content: workerURL }} spacing="extended">
            <code>{workerURL}</code>
          </Paragraph>
        </div>
      </div>

      <Divider margin={4}></Divider>
      <Tabs
        onChange={setActiveTab} // 新增：更新当前激活的标签页
        tabBarExtraContent={
          <Button
            theme="borderless"
            onClick={() => runWorker.mutate(editItem.UID)}
          >
            Run
          </Button>
        }
      >
        <TabPane
          itemKey="code"
          style={{ overflow: 'initial' }}
          tab={<span>Code</span>}
        >
          {worker ? (
            <div className="flex flex-col my-1">
              <div>
                <MonacoEditor uid={worker.UID} />
              </div>
            </div>
          ) : null}
        </TabPane>
        <TabPane itemKey="config" tab={<span>Config</span>}>
          <div className="flex flex-col">
            <div className="flex flex-row m-2">
              <p className="self-center">Entry: </p>
              <div className="grid grid-cols-1 lg:grid-cols-2">
                <Input
                  addonBefore={
                    <p className="invisible w-0 sm:visible sm:w-auto">
                      {appConf?.Scheme}://
                    </p>
                  }
                  addonAfter={
                    <p className="invisible w-0 sm:visible w-25">
                      {appConf?.WorkerURLSuffix}
                    </p>
                  }
                  value={editItem.Name}
                  defaultValue={worker?.Name}
                  onChange={(value) => {
                    if (worker) {
                      setEditItem((item) => ({ ...item, Name: value }))
                    }
                  }}
                />
              </div>
            </div>
            <div className="flex flex-row m-2">
              <p className="self-center">Node: </p>
              <Select
                placeholder="请选择节点"
                style={{ width: 180 }}
                optionList={resp?.data.nodes.map((node) => {
                  return {
                    label: node.Name,
                    value: node.Name,
                  }
                })}
                value={editItem.NodeName}
                onChange={(value) => {
                  if (worker) {
                    setEditItem((item) => ({
                      ...item,
                      NodeName: value as string,
                    }))
                  }
                }}
              ></Select>
            </div>
          </div>
        </TabPane>
        <TabPane
          itemKey="template"
          style={{ overflow: 'initial' }}
          tab={<span>Template</span>}
        >
          <TemplateEditor content={templateContent} setContent={setTemplateContent} />
        </TabPane>
        <TabPane
          itemKey="logs"
          style={{ overflow: 'initial' }}
          tab={<span>Logs</span>}
        >
          <Pagination total={totalLogs} currentPage={page} onPageChange={setPage} pageSize={pageSize} style={{ marginBottom: 12 }} />
          <List
            layout="vertical"
            dataSource={logs}
            split={false}
            renderItem={(item) => (
              <List.Item
                style={{ padding: '5px 0px', fontSize: '16px' }}
                main={
                  <div>
                    <Tag size="small" color='light-blue'>{formatDate(new Date(item.time))}</Tag>
                    <span className="ml-2">{item.output}</span>
                  </div>
                }
              />
            )}
          />
          <Pagination total={totalLogs} currentPage={page} onPageChange={setPage} pageSize={pageSize} style={{ marginBottom: 12 }} />

        </TabPane>
      </Tabs>
    </div>
  )
}
