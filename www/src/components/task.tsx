import React, { useState, useEffect, useCallback } from 'react';
import { Breadcrumb, ButtonGroup, Button, Card, List, Modal, Form, Descriptions } from '@douyinfe/semi-ui';
import { IconHome } from '@douyinfe/semi-icons';
import { t } from '@/lib/i18n';
import { listTasks } from '@/api/workers';
import type { Task } from '@/types/workers';
import { Pagination } from '@douyinfe/semi-ui';
import router from 'next/router';

export const TaskComponent: React.FC = () => {
    const [tasks, setTasks] = useState<Task[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [total, setTotal] = useState(0);
    const fetchTasks = useCallback(async (_page: number, _page_size: number) => {
        setLoading(true);
        try {
            const data = (await listTasks(_page, _page_size)).data.data;
            setTasks(data.tasks);
            setTotal(data.total);
        }
        catch (err) {
            setError('Failed to fetch tasks');
        }
        finally {
            setLoading(false);
        }

    }, [])

    useEffect(() => {
        fetchTasks(page, pageSize);
    }, [page, pageSize, fetchTasks]);

    return (
        <div className="m-4">
            <Breadcrumb>
                <Breadcrumb.Item href="/">
                    <IconHome />
                </Breadcrumb.Item>
                <Breadcrumb.Item>{t.task}</Breadcrumb.Item>
            </Breadcrumb>
            <Pagination total={total} currentPage={page} onPageChange={setPage} pageSize={pageSize} style={{ marginBottom: 12 }}></Pagination>
            <List
                layout="vertical"
                dataSource={tasks}
                loading={loading}
                renderItem={(item) => (
                    <List.Item
                        header={item.worker_name}
                        extra={
                            <ButtonGroup theme="borderless">
                                <Button onClick={() => {
                                    router.push({
                                        pathname: '/worker',
                                        query: { UID: item.worker_uid },
                                    })
                                }}>{t.look}</Button>
                                <Button onClick={() => {
                                    router.push({
                                        pathname: '/logs',
                                        query: { trace_id: item.trace_id, worker_uid: item.worker_uid },
                                    })
                                }}>{t.log}</Button>
                            </ButtonGroup>
                        }
                        main={
                            <div>
                                <div style={{ color: 'var(--semi-color-text-0)', fontWeight: 500 }}>{t[item.status]}</div>
                                {item.trace_id}
                            </div>
                        }
                    />
                )}
            />

        </div>
    );
};