import React, { useState, useEffect, useCallback } from 'react';
import { Breadcrumb, ButtonGroup, Button, Card, List, Modal, Form, Descriptions } from '@douyinfe/semi-ui';
import { IconHome } from '@douyinfe/semi-icons';
import { t } from '@/lib/i18n';
import { createAccessKey, getAccessKeys, deleteAccessKey } from '@/api/auth';
import { changePassword } from '@/api/users';
import type { AccessKey } from '@/types/user';
import { $user } from '@/store/userState';
import { useStore } from '@nanostores/react';

export const UserCom: React.FC = () => {
    const [accessKeys, setAccessKeys] = useState<AccessKey[]>([]);
    const [openDeleteDialog, setOpenDeleteDialog] = useState(false);
    const [openChangePasswordDialog, setOpenChangePasswordDialog] = useState(false);
    const [newPassword, setNewPassword] = useState('');
    const [isChangingPassword, setIsChangingPassword] = useState(false);
    const [accessKeyToDelete, setAccessKeyToDelete] = useState<string | null>(null);
    const [openCreateDialog, setOpenCreateDialog] = useState(false);
    const [newAccessKeyName, setNewAccessKeyName] = useState('');
    const [isCreating, setIsCreating] = useState(false);
    const [isDeleting, setIsDeleting] = useState(false);
    const user = useStore($user);

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

    const handleChangePasswordClick = () => {
        setOpenChangePasswordDialog(true);
    };

    const handleChangePasswordConfirm = async () => {
        if (!newPassword) {
            alert('旧密码和新密码不能为空');
            return;
        }
        if (!user) return;

        setIsChangingPassword(true);
        try {
            await changePassword(user.id, newPassword);
            setOpenChangePasswordDialog(false);
            setNewPassword('');
        } catch (error) {
            console.error('修改密码失败:', error);
            alert('修改密码失败，请检查旧密码是否正确');
        } finally {
            setIsChangingPassword(false);
        }
    };

    const handleChangePasswordCancel = () => {
        setOpenChangePasswordDialog(false);
        setNewPassword('');
    };

    return (
        <div className="m-4">
            <div className="flex justify-between">
                <Breadcrumb>
                    <Breadcrumb.Item
                        href="/admin"
                        icon={<IconHome size="small" />}
                    ></Breadcrumb.Item>
                    <Breadcrumb.Item href="/admin">User</Breadcrumb.Item>
                </Breadcrumb>
            </div>
            <p>Config</p>
            <Card style={{ width: '100%', marginBottom: '16px' }}>
                <div className="flex justify-between items-center">
                    <span>用户密码</span>
                    <Button onClick={handleChangePasswordClick} disabled={isChangingPassword}>修改密码</Button>
                </div>
            </Card>
            <div className='flex justify-between'>
                <p>Access Key</p>
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
            <Modal
                title="修改密码"
                visible={openChangePasswordDialog}
                onOk={handleChangePasswordConfirm}
                onCancel={handleChangePasswordCancel}
                maskClosable={false}
                closeOnEsc={true}
            >
                <div style={{ padding: '16px' }}>
                    <Form>
                        <Form.Input
                            field='newPassword'
                            label="新密码"
                            type="password"
                            initValue={newPassword}
                            onChange={(value) => setNewPassword(value)}
                        />
                    </Form>
                </div>
            </Modal>
        </div>
    );
};