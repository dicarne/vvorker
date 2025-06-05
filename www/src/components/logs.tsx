import React, { useState, useEffect, useCallback } from 'react';
import { Breadcrumb, ButtonGroup, Button, Card, List, Modal, Form, Descriptions } from '@douyinfe/semi-ui';
import { IconHome } from '@douyinfe/semi-icons';
import { t } from '@/lib/i18n';
import { getLogs } from '@/api/workers';
import type { TaskLog } from '@/types/workers';
import { Pagination } from '@douyinfe/semi-ui';
import router from 'next/router';

export const LogsComponent: React.FC = () => {
    const [tasks, setTasks] = useState<TaskLog[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [total, setTotal] = useState(0);
    const { trace_id, worker_uid } = router.query
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
                        main={
                            <div>
                                {item.content}
                            </div>
                        }
                    />
                )}
            />

        </div>
    );
};