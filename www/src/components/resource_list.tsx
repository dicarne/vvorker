import React, { useState, useEffect, useCallback } from 'react';
import { Dialog, Button, Card, List, ListItem, ListItemText, IconButton, TextField } from '@mui/material';
import DeleteIcon from '@mui/icons-material/Delete';
import { getResourceList, deleteResource, createResource } from '@/api/resources';
import { ResourceData } from '@/types/resources';

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
    }, []);

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
        <div>
            <Button
                variant="contained"
                color="primary"
                onClick={handleCreateResource}
                disabled={isCreating} // 创建时禁用按钮
            >
                创建资源
            </Button>
            <List>
                {resources.map(resource => (
                    <ListItem key={resource.uid}>
                        <Card sx={{ width: '100%', padding: '16px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                            <ListItemText primary={resource.name} />
                            <ListItemText primary={"id: " + resource.uid} style={{ color: "#909090" }} />
                            <IconButton
                                onClick={() => handleDeleteClick(resource.uid)}
                                disabled={isDeleting} // 删除时禁用按钮
                            >
                                <DeleteIcon />
                            </IconButton>
                        </Card>
                    </ListItem>
                ))}
            </List>
            {/* 创建资源对话框 */}
            <Dialog
                open={openCreateDialog}
                // 加载时阻止关闭
                onClose={isCreating ? () => { } : handleCreateCancel}
            >
                <div style={{ padding: '16px' }}>
                    <h2>创建{rtype}</h2>
                    <TextField
                        label="资源名称"
                        value={newResourceName}
                        onChange={(e) => setNewResourceName(e.target.value)}
                        fullWidth
                        margin="normal"
                        disabled={isCreating} // 加载时禁用输入框
                    />
                    <Button
                        onClick={handleCreateCancel}
                        disabled={isCreating} // 加载时禁用按钮
                    >
                        取消
                    </Button>
                    <Button
                        onClick={handleCreateConfirm}
                        color="primary"
                        disabled={isCreating} // 加载时禁用按钮
                    >
                        {isCreating ? '创建中...' : '创建'}
                    </Button>
                </div>
            </Dialog>
            <Dialog
                open={openDeleteDialog}
                // 加载时阻止关闭
                onClose={isDeleting ? () => { } : handleDeleteCancel}
            >
                <div style={{ padding: '16px' }}>
                    <p>确定要删除这个资源吗？</p>
                    <Button
                        onClick={handleDeleteCancel}
                        disabled={isDeleting} // 加载时禁用按钮
                    >
                        取消
                    </Button>
                    <Button
                        onClick={handleDeleteConfirm}
                        color="error"
                        disabled={isDeleting} // 加载时禁用按钮
                    >
                        {isDeleting ? '删除中...' : '删除'}
                    </Button>
                </div>
            </Dialog>
        </div>
    );
};
