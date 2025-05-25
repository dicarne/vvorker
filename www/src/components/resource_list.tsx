import React, { useState, useEffect, useCallback } from 'react';
import { ListItem, TextField, ListItemText } from '@mui/material';
import { getResourceList, deleteResource, createResource } from '@/api/resources';
import { ResourceData } from '@/types/resources';
import { Breadcrumb, ButtonGroup, Button, Card, List, Modal, Form, Descriptions } from '@douyinfe/semi-ui';
import { IconHome } from '@douyinfe/semi-icons';
import { t } from '@/lib/i18n';

const fetchResources = async (type: string) => {
    return await getResourceList(0, 10000, type)
};



// 定义组件的 Props 类型
type ResourceListProps = {
    rtype: string;
};

export const ResourceList: React.FC<ResourceListProps> = ({ rtype }) => {
    const [resources, setResources] = useState<ResourceData[]>([]);
    const [openDeleteDialog, setOpenDeleteDialog] = useState(false);
    const [resourceToDelete, setResourceToDelete] = useState<string | null>(null);
    // 新增状态控制创建资源对话框的显示
    const [openCreateDialog, setOpenCreateDialog] = useState(false);
    // 新增状态存储用户输入的资源名称
    const [newResourceName, setNewResourceName] = useState('');

    const loadResources = useCallback(async () => {
        const data = await fetchResources(rtype);
        setResources(data.data);
    }, [rtype]);

    // 页面加载时获取资源列表
    useEffect(() => {
        loadResources();
    }, [loadResources]);

    // 处理创建资源按钮点击事件
    const handleCreateResource = () => {
        setOpenCreateDialog(true);
    };

    // 新增加载状态
    const [isCreating, setIsCreating] = useState(false);
    const [isDeleting, setIsDeleting] = useState(false);

    // 处理创建资源对话框的确认按钮点击事件
    const handleCreateConfirm = async () => {
        const trimmedName = newResourceName.trim();
        if (!trimmedName) {
            // 提示用户资源名称不能为空
            alert('资源名称不能为空，请输入有效的名称。');
            return;
        }
        // 开始创建，设置加载状态
        setIsCreating(true);
        try {
            const newResource = await createResource(trimmedName, rtype);
            setResources([...resources, newResource as ResourceData]);
        } catch (error) {
            console.error('创建资源失败:', error);
        } finally {
            // 结束创建，取消加载状态
            setIsCreating(false);
            setOpenCreateDialog(false);
            setNewResourceName('');
        }
    };

    // 处理创建资源对话框的取消按钮点击事件
    const handleCreateCancel = () => {
        setOpenCreateDialog(false);
        setNewResourceName('');
    };

    // 处理删除按钮点击事件
    const handleDeleteClick = (id: string) => {
        setResourceToDelete(id);
        setOpenDeleteDialog(true);
    };

    // 处理删除确认弹窗的确认按钮点击事件
    const handleDeleteConfirm = async () => {
        if (resourceToDelete !== null) {
            // 开始删除，设置加载状态
            setIsDeleting(true);
            try {
                await deleteResourcefn(resourceToDelete, rtype);
                setResources(resources.filter(resource => resource.uid !== resourceToDelete));
            } catch (error) {
                console.error('删除资源失败:', error);
            } finally {
                // 结束删除，取消加载状态
                setIsDeleting(false);
                setOpenDeleteDialog(false);
                setResourceToDelete(null);
            }
        }
    };

    // 处理删除确认弹窗的取消按钮点击事件
    const handleDeleteCancel = () => {
        setOpenDeleteDialog(false);
        setResourceToDelete(null);
    };

    const deleteResourcefn = async (id: string, type: string) => {
        await deleteResource(id, type);
        await loadResources();
    };

    return (
        <div className="m-4">
            <div className="flex justify-between">
                <Breadcrumb>
                    {/* <Breadcrumb.Item key={}
                        href="/admin"
                        icon={<IconHome size="small" />}
                    ></Breadcrumb.Item>
                    <Breadcrumb.Item href="/admin">{rtype.toUpperCase()}</Breadcrumb.Item> */}
                    <p>hi</p>
                </Breadcrumb>
                <ButtonGroup>
                    <Button onClick={handleCreateResource} disabled={isCreating}>{t.create}</Button>
                </ButtonGroup>
            </div>
            <List>
                {resources.map(resource => (
                    <ListItem key={resource.uid}>
                        <Card title={resource.name} style={{ width: '100%' }} headerExtraContent={
                            <ButtonGroup>
                                <Button onClick={() => handleDeleteClick(resource.uid)} disabled={isDeleting}>{t.delete}</Button>
                            </ButtonGroup>
                        }>

                            <Descriptions data={[
                                {

                                    key: t.id,
                                    value: resource.uid

                                }
                            ]} />
                        </Card>
                    </ListItem>
                ))}
            </List>
            {/* 创建资源对话框 */}
            <Modal
                title={t.create + " " + rtype.toUpperCase()}
                visible={openCreateDialog}
                onOk={handleCreateConfirm}
                onCancel={handleCreateCancel}
                maskClosable={false}
                closeOnEsc={true}
            >
                <div style={{ padding: '16px' }}>
                    {/* <Input
                        label={t.resourceName}
                        value={newResourceName}
                        onChange={(e) => setNewResourceName(e.target.value)}
                        fullWidth
                        margin="normal"
                        disabled={isCreating} // 加载时禁用输入框
                    /> */}
                    <Form onValueChange={values => setNewResourceName(values.name)}>
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
