import React, { useState, useEffect, useCallback } from 'react';
import { Breadcrumb, ButtonGroup, Button, Card, List, Modal, Form, Descriptions } from '@douyinfe/semi-ui';
import { IconHome } from '@douyinfe/semi-icons';
import { t } from '@/lib/i18n';
import { getLogs, interruptTask } from '@/api/workers';
import type { TaskLog } from '@/types/workers';
import { Pagination, Tag } from '@douyinfe/semi-ui';
import router from 'next/router';

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

export const LogsComponent: React.FC = () => {
    const [tasks, setTasks] = useState<TaskLog[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(25);
    const [total, setTotal] = useState(0);
    const { trace_id, worker_uid } = router.query
    console.log(router.query)
    const fetchTasks = useCallback(async (_page: number, _page_size: number) => {
        setLoading(true);
        try {
            const data = (await getLogs(worker_uid as string, trace_id as string, _page, _page_size)).data.data;
            setTasks(data.logs);
            setTotal(data.total);
        }
        catch (err) {
            setError('Failed to fetch tasks');
        }
        finally {
            setLoading(false);
        }

    }, [worker_uid, trace_id])

    useEffect(() => {
        fetchTasks(page, pageSize);
        // 设置每秒刷新一次
        const intervalId = setInterval(() => {
            fetchTasks(page, pageSize);
        }, 1000);

        // 组件卸载时清除定时器
        return () => clearInterval(intervalId);
    }, [page, pageSize, fetchTasks]);

    return (
        <div className="m-4 flex flex-col">
            <div className="flex flex-row justify-between">
                <div className="flex flex-col">
                    <Breadcrumb>
                        <Breadcrumb.Item href="/admin">
                            <IconHome />
                        </Breadcrumb.Item>
                        <Breadcrumb.Item href="/task">{t.task}</Breadcrumb.Item>
                    </Breadcrumb>
                </div>
                <div className="flex flex-col">
                    <ButtonGroup>
                        <Button onClick={() => interruptTask(trace_id as string, worker_uid as string)}>{t.cancel}</Button>
                    </ButtonGroup>
                </div>
            </div>
            <Pagination total={total} currentPage={page} onPageChange={setPage} pageSize={pageSize} style={{ marginBottom: 12 }}></Pagination>
            <List
                layout="vertical"
                dataSource={tasks}
                split={false}
                renderItem={(item) => (
                    <List.Item
                        style={{ padding: '5px 0px', fontSize: '16px' }}
                        main={
                            <div>
                                {/* 格式化时间 */}
                                <Tag size="small" color='light-blue'>{formatDate(new Date(item.time))}</Tag>
                                <span className="ml-2">{item.content}</span>
                            </div>
                        }
                    />
                )}
            />
        </div>
    );
};