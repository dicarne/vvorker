import React, { useState, useEffect, useCallback } from 'react';
import { Breadcrumb, ButtonGroup, Button, Card, List, Modal, Form, Descriptions } from '@douyinfe/semi-ui';
import { IconHome } from '@douyinfe/semi-icons';
import { t } from '@/lib/i18n';
import { createAccessKey, getAccessKeys, deleteAccessKey } from '@/api/auth';
import type { AccessKey } from '@/types/user';

export const UserCom: React.FC = () => {
    const [accessKeys, setAccessKeys] = useState<AccessKey[]>([]);
    const [openDeleteDialog, setOpenDeleteDialog] = useState(false);
    const [accessKeyToDelete, setAccessKeyToDelete] = useState<string | null>(null);
    const [openCreateDialog, setOpenCreateDialog] = useState(false);
    const [newAccessKeyName, setNewAccessKeyName] = useState('');
    const [isCreating, setIsCreating] = useState(false);
    const [isDeleting, setIsDeleting] = useState(false);

    const loadAccessKeys = useCallback(async () => {
        try {
            const data = await getAccessKeys();
            setAccessKeys(data.data.data);
        } catch (error) {
            console.error('获取 Access Key 列表失败:', error);
        }
    }, []);

    useEffect(() => {
        loadAccessKeys();
    }, [loadAccessKeys]);

    const handleCreateAccessKey = () => {
        setOpenCreateDialog(true);
    };

    const handleCreateConfirm = async () => {
        const trimmedName = newAccessKeyName.trim();
        if (!trimmedName) {
            alert('Access Key 名称不能为空，请输入有效的名称。');
            return;
        }
        setIsCreating(true);
        try {
            const newKey = await createAccessKey(trimmedName);
            setAccessKeys([...accessKeys, newKey.data.data as AccessKey]);
        } catch (error) {
            console.error('创建 Access Key 失败:', error);
        } finally {
            setIsCreating(false);
            setOpenCreateDialog(false);
            setNewAccessKeyName('');
        }
    };

    const handleCreateCancel = () => {
        setOpenCreateDialog(false);
        setNewAccessKeyName('');
    };

    const handleDeleteClick = (id: string) => {
        setAccessKeyToDelete(id);
        setOpenDeleteDialog(true);
    };

    const handleDeleteConfirm = async () => {
        if (accessKeyToDelete !== null) {
            setIsDeleting(true);
            try {
                await deleteAccessKey(accessKeyToDelete);
                setAccessKeys(accessKeys.filter(key => key.key !== accessKeyToDelete));
            } catch (error) {
                console.error('删除 Access Key 失败:', error);
            } finally {
                setIsDeleting(false);
                setOpenDeleteDialog(false);
                setAccessKeyToDelete(null);
            }
        }
    };

    const handleDeleteCancel = () => {
        setOpenDeleteDialog(false);
        setAccessKeyToDelete(null);
    };

    return (
        <div className="m-4">
            <div className="flex justify-between">
                <Breadcrumb>
                    <Breadcrumb.Item
                        href="/admin"
                        icon={<IconHome size="small" />}
                    ></Breadcrumb.Item>
                    <Breadcrumb.Item href="/admin">Access Key</Breadcrumb.Item>
                </Breadcrumb>
                <ButtonGroup>
                    <Button onClick={handleCreateAccessKey} disabled={isCreating}>{t.create}</Button>
                </ButtonGroup>
            </div>
            <List>
                {accessKeys.map(key => (
                    <List.Item key={key.key}>
                        <Card title={key.name} style={{ width: '100%' }} headerExtraContent={
                            <ButtonGroup>
                                <Button onClick={() => handleDeleteClick(key.key)} disabled={isDeleting}>{t.delete}</Button>
                            </ButtonGroup>
                        }>
                            <Descriptions data={[
                                {
                                    key: t.id,
                                    value: key.key
                                }
                            ]} />
                        </Card>
                    </List.Item>
                ))}
            </List>
            <Modal
                title={t.create + " Access Key"}
                visible={openCreateDialog}
                onOk={handleCreateConfirm}
                onCancel={handleCreateCancel}
                maskClosable={false}
                closeOnEsc={true}
            >
                <div style={{ padding: '16px' }}>
                    <Form onValueChange={values => setNewAccessKeyName(values.name)}>
                        <Form.Input field='name' label={t.resourceName} />
                    </Form>
                </div>
            </Modal>
            <Modal
                title="基本对话框"
                visible={openDeleteDialog}
                onOk={handleDeleteConfirm}
                onCancel={handleDeleteCancel}
                maskClosable={false}
                closeOnEsc={true}
            >
                <div style={{ padding: '16px' }}>
                    <p>{t.warnDeleteResource}</p>
                </div>
            </Modal>
        </div>
    );
};