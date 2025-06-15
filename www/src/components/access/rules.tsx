import React, { useState, useEffect, useRef } from 'react';
import { Tabs, Switch, List, Tag, Button, Modal, Form, Input, Select, Col, Row, ButtonGroup } from '@douyinfe/semi-ui';
import {
    getAccessControl,
    updateEnableAccessControl,
    listAccessRules,
    addAccessRule,
    deleteAccessRule
} from '@/api/workers';
import { AccessRule, EnableAccessControlRequest, DeleteAccessRuleRequest } from '@/types/access';
import { t } from '@/lib/i18n';

interface RulesTabPaneProps {
    workerUid: string;
}


const RulesTabPane: React.FC<RulesTabPaneProps> = ({ workerUid }) => {
    const [isModalVisible, setIsModalVisible] = useState(false);
    const [rules, setRules] = useState<AccessRule[]>([]);
    const [isEnabled, setIsEnabled] = useState(false);
    const formRef = useRef<Form>(null);

    // 获取规则控制状态
    const fetchAccessControl = async () => {
        try {
            const response = await getAccessControl({ worker_uid: workerUid });
            setIsEnabled(response.data.EnableAccessControl);
        } catch (error) {
            console.error('Failed to fetch access control status', error);
        }
    };

    // 获取规则列表
    const fetchRules = async () => {
        try {
            const response = await listAccessRules({
                worker_uid: workerUid,
                page: 1,
                page_size: 100
            });
            setRules(response.data.access_rules);
        } catch (error) {
            console.error('Failed to fetch access rules', error);
        }
    };

    // 切换规则启用状态
    const handleSwitchChange = async (checked: boolean) => {
        try {
            const request: EnableAccessControlRequest = {
                enable: checked,
                worker_uid: workerUid
            };
            await updateEnableAccessControl(request);
            setIsEnabled(checked);
        } catch (error) {
            console.error('Failed to update access control status', error);
        }
    };

    // 显示新增规则弹窗
    const showModal = () => {
        setIsModalVisible(true);
    };

    // 处理新增规则提交
    const handleOk = async () => {
        try {
            const v = formRef.current?.formApi.getValues();
            await addAccessRule({
                worker_uid: workerUid,
                rule_type: v.rule_type,
                path: v.path,
                description: v.description
            });
            setIsModalVisible(false);
            fetchRules();
        } catch (error) {
            console.error('Failed to add access rule', error);
        }
    };

    // 处理取消新增规则
    const handleCancel = () => {
        setIsModalVisible(false);
    };

    // 删除规则
    const handleDeleteRule = async (id: number) => {
        try {
            const request: DeleteAccessRuleRequest = { worker_uid: workerUid, rule_id: id };
            await deleteAccessRule(request);
            fetchRules();
        } catch (error) {
            console.error('Failed to delete access rule', error);
        }
    };

    useEffect(() => {
        fetchAccessControl();
        fetchRules();
    }, [workerUid]);

    return (
        <>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
                <div className='flex items-center px-2'>
                    <span className='pr-6'>{t.enableAccessControl} </span>
                    <Switch checked={isEnabled} onChange={handleSwitchChange} />
                </div>
                <Button type="primary" onClick={showModal}>
                    {t.addRule}
                </Button>
            </div>
            <List
                layout="horizontal"
                dataSource={rules}
                split
                renderItem={(item) => (
                    <List.Item
                        style={{ padding: '10px', width: "100%", border: '1px solid var(--semi-color-border)', }}
                        main={
                            <div>
                                <div style={{ color: 'var(--semi-color-text-0)', fontWeight: 500 }}>{item.path}</div>
                                <div>{item.rule_type}</div>
                            </div>
                        }
                        extra={
                            <ButtonGroup theme="borderless">
                                <Button onClick={() => handleDeleteRule(item.id!)}>删除</Button>
                                <Button>更多</Button>
                            </ButtonGroup>
                        }
                    />
                )}
            />
            <Modal title={t.addRule} visible={isModalVisible} onOk={handleOk} onCancel={handleCancel}>
                <Form ref={formRef}>
                    <Row>
                        <Form.Input
                            field="path"
                            label={t.prefix}
                            initValue={'/'}
                            trigger='blur'
                        />
                    </Row>
                    <Row>
                        <Form.Select
                            field="rule_type"
                            placeholder='请选择控制类型'
                            label="该路径前缀的控制类型"
                            initValue={'internal'}
                        >
                            <Form.Select.Option value="internal">内部认证</Form.Select.Option>
                            <Form.Select.Option value="token">外部TOKEN</Form.Select.Option>
                            <Form.Select.Option value="open">开放</Form.Select.Option>
                        </Form.Select>
                    </Row>
                    <Row>
                        <Form.Input field="description" label="描述" placeholder="" />
                    </Row>
                </Form>
            </Modal>
        </>
    );
};

export default RulesTabPane;
