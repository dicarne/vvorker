import React, { useState, useEffect, useCallback } from 'react';
import { Breadcrumb, ButtonGroup, Button, Card, List, Modal, Form, Descriptions, Toast } from '@douyinfe/semi-ui';
import { IconHome } from '@douyinfe/semi-icons';
import { t } from '@/lib/i18n';

interface User {
    ID: number;
    user_name: string;
    email: string;
    role: string;
    createdAt: string;
    status: number;
}

import { getUsers, createUser, deleteUser, banUser, unbanUser, changePassword } from '@/api/users';
import { useStore } from '@nanostores/react';
import { $user } from '@/store/userState';

export const UsersManagementCom: React.FC = () => {
    const [users, setUsers] = useState<User[]>([]);
    const [openDeleteDialog, setOpenDeleteDialog] = useState(false);
    const [userToDelete, setUserToDelete] = useState<number | null>(null);
    const [openCreateDialog, setOpenCreateDialog] = useState(false);
    const [openChangePasswordDialog, setOpenChangePasswordDialog] = useState(false);
    const [userToChangePassword, setUserToChangePassword] = useState<number | null>(null);
    const [newPassword, setNewPassword] = useState('');
    const [newUserForm, setNewUserForm] = useState({ username: '', email: '', password: '' });
    const [isLoading, setIsLoading] = useState(false);
    const user = useStore($user);

    const loadUsers = useCallback(async () => {
        setIsLoading(true);
        try {
            const response = await getUsers();
            console.log(response)
            setUsers((response.data.data.users as User[]).filter(u => u.user_name !== user?.userName));
        } catch (error) {
            Toast.error('Failed to load users');
            console.error('Failed to load users:', error);
        } finally {
            setIsLoading(false);
        }
    }, []);

    useEffect(() => {
        loadUsers();
    }, [loadUsers]);

    const handleCreateUser = () => {
        setOpenCreateDialog(true);
    };

    const handleCreateConfirm = async () => {
        if (!newUserForm.username || !newUserForm.email || !newUserForm.password) {
            Toast.warning('Please fill in all required fields');
            return;
        }
        setIsLoading(true);
        try {
            await createUser(newUserForm);
            Toast.success('User created successfully');
            loadUsers();
            setOpenCreateDialog(false);
            setNewUserForm({ username: '', email: '', password: '' });
        } catch (error) {
            Toast.error('Failed to create user');
            console.error('Failed to create user:', error);
        } finally {
            setIsLoading(false);
        }
    };

    const handleDeleteClick = (userId: number) => {
        setUserToDelete(userId);
        setOpenDeleteDialog(true);
    };

    const handleDeleteConfirm = async () => {
        if (!userToDelete) return;
        setIsLoading(true);
        try {
            await deleteUser(userToDelete);
            Toast.success('User deleted successfully');
            loadUsers();
        } catch (error) {
            Toast.error('Failed to delete user');
            console.error('Failed to delete user:', error);
        } finally {
            setIsLoading(false);
            setOpenDeleteDialog(false);
            setUserToDelete(null);
        }
    };

    const handleBanUser = async (userId: number, isBanned: boolean) => {
        setIsLoading(true);
        try {
            if (isBanned) {
                await unbanUser(userId);
            } else {
                await banUser(userId);
            }
            Toast.success(`User ${isBanned ? 'unbanned' : 'banned'} successfully`);
            loadUsers();
        } catch (error) {
            Toast.error(`Failed to ${isBanned ? 'unban' : 'ban'} user`);
            console.error(`Failed to ${isBanned ? 'unban' : 'ban'} user:`, error);
        } finally {
            setIsLoading(false);
        }
    };

    const handleChangePasswordClick = (userId: number) => {
        setUserToChangePassword(userId);
        setNewPassword('');
        setOpenChangePasswordDialog(true);
    };

    const handleChangePasswordConfirm = async () => {
        if (!userToChangePassword || !newPassword) {
            Toast.warning('Please enter a new password');
            return;
        }
        setIsLoading(true);
        try {
            await changePassword(userToChangePassword, newPassword);
            Toast.success('Password changed successfully');
            setOpenChangePasswordDialog(false);
            setUserToChangePassword(null);
            setNewPassword('');
        } catch (error) {
            Toast.error('Failed to change password');
            console.error('Failed to change password:', error);
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="m-4">
            <div className="flex justify-between">
                <Breadcrumb>
                    <Breadcrumb.Item
                        href="/admin"
                        icon={<IconHome size="small" />}
                    />
                    <Breadcrumb.Item>Users Management</Breadcrumb.Item>
                </Breadcrumb>
                <ButtonGroup>
                    <Button onClick={handleCreateUser} disabled={isLoading}>Create User</Button>
                </ButtonGroup>
            </div>
            <List loading={isLoading}>
                {users.map(user => (
                    <List.Item key={user.user_name}>
                        <Card
                            title={user.user_name}
                            style={{ width: '100%' }}
                            headerExtraContent={
                                <ButtonGroup>
                                    <Button
                                        type={user.status === 1 ? 'secondary' : 'warning'}
                                        onClick={() => handleChangePasswordClick(user.ID)}
                                        disabled={isLoading || user.role === 'admin'}
                                    >
                                        Change Password
                                    </Button>
                                    <Button
                                        type={user.status === 1 ? 'secondary' : 'warning'}
                                        onClick={() => handleBanUser(user.ID, user.status === 1)}
                                        disabled={isLoading || user.role === 'admin'}
                                    >
                                        {user.status === 1 ? 'Unban' : 'Ban'}
                                    </Button>
                                    <Button
                                        type="danger"
                                        onClick={() => handleDeleteClick(user.ID)}
                                        disabled={isLoading || user.role === 'admin'}
                                    >
                                        Delete
                                    </Button>
                                </ButtonGroup>
                            }
                        >
                            <Descriptions data={[
                                { key: 'Email', value: user.email },
                                { key: 'Role', value: user.role === 'admin' ? 'Admin' : 'User' },
                                { key: 'Status', value: user.status === 1 ? 'Banned' : 'Active' },
                                { key: 'Created At', value: new Date(user.createdAt).toLocaleString() }
                            ]} />
                        </Card>
                    </List.Item>
                ))}
            </List>

            <Modal
                title="Create New User"
                visible={openCreateDialog}
                onOk={handleCreateConfirm}
                onCancel={() => setOpenCreateDialog(false)}
                maskClosable={false}
                closeOnEsc={true}
            >
                <Form>
                    <Form.Input
                        field="username"
                        label="Username"
                        required
                        initValue={newUserForm.username}
                        onChange={val => setNewUserForm(prev => ({ ...prev, username: val }))}
                    />
                    <Form.Input
                        field="email"
                        label="Email"
                        required
                        initValue={newUserForm.email}
                        onChange={val => setNewUserForm(prev => ({ ...prev, email: val }))}
                    />
                    <Form.Input
                        field="password"
                        label="Password"
                        type="password"
                        required
                        initValue={newUserForm.password}
                        onChange={val => setNewUserForm(prev => ({ ...prev, password: val }))}
                    />
                </Form>
            </Modal>

            <Modal
                title="Delete User"
                visible={openDeleteDialog}
                onOk={handleDeleteConfirm}
                onCancel={() => setOpenDeleteDialog(false)}
                maskClosable={false}
                closeOnEsc={true}
            >
                <p>Are you sure you want to delete this user? This action cannot be undone.</p>
            </Modal>

            <Modal
                title="Change Password"
                visible={openChangePasswordDialog}
                onOk={handleChangePasswordConfirm}
                onCancel={() => {
                    setOpenChangePasswordDialog(false);
                    setUserToChangePassword(null);
                    setNewPassword('');
                }}
                maskClosable={false}
                closeOnEsc={true}
            >
                <Form>
                    <Form.Input
                        field="newPassword"
                        label="New Password"
                        type="password"
                        required
                        initValue={newPassword}
                        onChange={val => setNewPassword(val)}
                        placeholder="Enter new password"
                    />
                </Form>
            </Modal>
        </div>
    );
};