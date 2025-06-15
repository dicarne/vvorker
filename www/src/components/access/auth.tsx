import React, { useState, useEffect, useRef } from 'react';
import { Tabs, Button, Modal, Form, Input, List, Tag, ButtonGroup, Toast } from '@douyinfe/semi-ui';
import {
    createInternalWhiteList,
    listInternalWhiteLists,
    deleteInternalWhiteList,
    createAccessToken,
    listAccessTokens,
    deleteAccessToken
} from '@/api/workers';
import {
    InternalWhiteListCreateRequest,
    InternalWhiteListDeleteRequest,
    AccessTokenCreateRequest,
    AccessTokenDeleteRequest,
    InternalServerWhiteList,
    ExternalServerToken
} from '@/types/access';
import { t } from '@/lib/i18n'

interface AuthTabProps {
    workerUid: string;
}

const AuthTab: React.FC<AuthTabProps> = ({ workerUid }) => {
    const [isInternalModalVisible, setIsInternalModalVisible] = useState(false);
    const [isTokenModalVisible, setIsTokenModalVisible] = useState(false);
    const [isDeleteConfirmVisible, setIsDeleteConfirmVisible] = useState(false);
    const [internalWhiteLists, setInternalWhiteLists] = useState<InternalServerWhiteList[]>([]);
    const [accessTokens, setAccessTokens] = useState<ExternalServerToken[]>([]);
    const [deleteType, setDeleteType] = useState<'internal' | 'token'>('internal');
    const [deleteId, setDeleteId] = useState<number | string | null>(null);
    const internalFormRef = useRef<Form>(null);
    const tokenFormRef = useRef<Form>(null);
    // 新增状态，控制显示 token 的弹窗
    const [isShowTokenModalVisible, setIsShowTokenModalVisible] = useState(false);
    // 新增状态，存储新生成的 token
    const [newToken, setNewToken] = useState<string>('');

    // 获取内部白名单列表
    const fetchInternalWhiteLists = async () => {
        try {
            const response = await listInternalWhiteLists({ worker_uid: workerUid, page: 1, page_size: 100 });
            setInternalWhiteLists(response.data.internal_white_lists);
        } catch (error) {
            console.error('Failed to fetch internal white lists', error);
        }
    };

    // 获取访问令牌列表
    const fetchAccessTokens = async () => {
        try {
            const response = await listAccessTokens({ worker_uid: workerUid, page: 1, page_size: 100 });
            setAccessTokens(response.data.access_tokens);
        } catch (error) {
            console.error('Failed to fetch access tokens', error);
        }
    };

    // 显示新增内部白名单弹窗
    const showInternalModal = () => {
        setIsInternalModalVisible(true);
    };

    // 显示新增访问令牌弹窗
    const showTokenModal = () => {
        setIsTokenModalVisible(true);
    };

    // 处理新增内部白名单提交
    const handleInternalOk = async () => {
        try {
            const v = internalFormRef.current?.formApi.getValues();
            const request: InternalWhiteListCreateRequest = {
                worker_uid: workerUid,
                ...v
            };
            try {
                await createInternalWhiteList(request);
            } catch (error) {
                const throttleOpts = {
                    content: String(error),
                    duration: 10,
                    stack: true,
                };
                Toast.info(throttleOpts)
                throw error
            }


            setIsInternalModalVisible(false);
            fetchInternalWhiteLists();
        } catch (error) {
            console.error('Failed to add internal white list', error);
        }
    };

    // 处理新增访问令牌提交
    const handleTokenOk = async () => {
        try {
            const v = tokenFormRef.current?.formApi.getValues();
            const request: AccessTokenCreateRequest = {
                worker_uid: workerUid,
                ...v
            };

            try {
                const response = await createAccessToken(request);
                // 存储新生成的 token
                setNewToken(response.data.access_token);
                // 显示 token 弹窗
                setIsShowTokenModalVisible(true);
            } catch (error) {
                const throttleOpts = {
                    content: String(error),
                    duration: 10,
                    stack: true,
                };
                Toast.info(throttleOpts)
                throw error
            }

            setIsTokenModalVisible(false);
            fetchAccessTokens();
        } catch (error) {
            console.error('Failed to add access token', error);
        }
    };

    // 处理取消新增
    const handleCancel = () => {
        setIsInternalModalVisible(false);
        setIsTokenModalVisible(false);
    };

    // 显示删除确认弹窗
    const showDeleteConfirm = (type: 'internal' | 'token', id: number | string) => {
        setDeleteType(type);
        setDeleteId(id);
        setIsDeleteConfirmVisible(true);
    };

    // 处理删除确认
    const handleDeleteConfirm = async () => {
        try {
            if (deleteType === 'internal' && deleteId !== null) {
                const request: InternalWhiteListDeleteRequest = {
                    worker_uid: workerUid,
                    id: deleteId as number
                };
                try {
                    await deleteInternalWhiteList(request);
                } catch (error) {
                    const throttleOpts = {
                        content: String(error),
                        duration: 10,
                        stack: true,
                    };
                    Toast.info(throttleOpts)
                    throw error
                }

                fetchInternalWhiteLists();
            } else if (deleteType === 'token' && deleteId !== null) {
                const request: AccessTokenDeleteRequest = {
                    worker_uid: workerUid,
                    id: deleteId as number
                };
                try {
                    await deleteAccessToken(request);
                } catch (error) {
                    const throttleOpts = {
                        content: String(error),
                        duration: 10,
                        stack: true,
                    };
                    Toast.info(throttleOpts)
                    throw error
                }
                fetchAccessTokens();
            }
            setIsDeleteConfirmVisible(false);
        } catch (error) {
            console.error('Failed to delete item', error);
        }
    };

    // 处理取消删除
    const handleDeleteCancel = () => {
        setIsDeleteConfirmVisible(false);
    };

    useEffect(() => {
        fetchInternalWhiteLists();
        fetchAccessTokens();
    }, [workerUid]);

    return (
        <>
            <div style={{ marginBottom: 24 }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
                    <h3>{t.internalAccess}</h3>
                    <Button type="primary" onClick={showInternalModal}>
                        {t.add}
                    </Button>
                </div>
                <List
                    layout="horizontal"
                    dataSource={internalWhiteLists}
                    split
                    renderItem={(item) => (
                        <List.Item
                            style={{ padding: '10px', width: "100%", border: '1px solid var(--semi-color-border)', }}
                            main={
                                <div>
                                    <div style={{ color: 'var(--semi-color-text-0)', fontWeight: 500 }}>{item.WorkerName}</div>
                                    <div>{item.description}</div>
                                </div>
                            }
                            extra={
                                <ButtonGroup theme="borderless">
                                    <Button onClick={() => showDeleteConfirm('internal', item.id)}>删除</Button>
                                </ButtonGroup>
                            }
                        />
                    )}
                />
            </div>
            <div>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
                    <h3>{t.accessKey}</h3>
                    <Button type="primary" onClick={showTokenModal}>
                        {t.add}
                    </Button>
                </div>
                <List
                    layout="horizontal"
                    dataSource={accessTokens}
                    split
                    renderItem={(item) => (
                        <List.Item
                            style={{ padding: '10px', width: "100%", border: '1px solid var(--semi-color-border)', }}
                            main={
                                <div>
                                    <div style={{ color: 'var(--semi-color-text-0)', fontWeight: 500 }}>{item.token}</div>
                                    <div>{item.description}</div>
                                </div>
                            }
                            extra={
                                <ButtonGroup theme="borderless">
                                    <Button onClick={() => showDeleteConfirm('token', item.id)}>删除</Button>
                                </ButtonGroup>
                            }
                        />
                    )}
                />
            </div>
            <Modal title={t.addInternalAccess} visible={isInternalModalVisible} onOk={handleInternalOk} onCancel={handleCancel}>
                <Form ref={internalFormRef}>
                    <Form.Input
                        field="name"
                        label="Name"
                        initValue=""
                        trigger='blur'
                    />
                    <Form.Input
                        field="description"
                        label="Description"
                        initValue=""
                        trigger='blur'
                    />
                </Form>
            </Modal>
            <Modal title={t.addAccessKey} visible={isTokenModalVisible} onOk={handleTokenOk} onCancel={handleCancel}>
                <Form ref={tokenFormRef}>
                    <Form.Input
                        field="description"
                        label="Description"
                        initValue=""
                        trigger='blur'
                    />
                </Form>
            </Modal>
            <Modal
                title={t.delete}
                visible={isDeleteConfirmVisible}
                onOk={handleDeleteConfirm}
                onCancel={handleDeleteCancel}
                okType="danger"
            >
                <p>{t.deleteConfirm}</p>
            </Modal>
            <Modal
                title={t.important}
                visible={isShowTokenModalVisible}
                onOk={() => setIsShowTokenModalVisible(false)}
                onCancel={() => setIsShowTokenModalVisible(false)}
            >
                <p style={{
                    textAlign: "center"
                }}>{t.tokenOnce}</p>
                <pre style={{
                    background: "#ededed",
                    padding: 10,
                    borderRadius: 5,
                    marginTop: 10,
                    textAlign: "center"

                }}>{newToken}</pre>
            </Modal>
        </>
    );
};

export default AuthTab;
