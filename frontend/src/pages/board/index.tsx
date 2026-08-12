import React, { useEffect, useState } from 'react';
import {
  Row,
  Col,
  Card,
  Button,
  Modal,
  Form,
  Input,
  Select,
  Tag,
  Dropdown,
  Menu,
  Typography,
  Space,
  message,
} from 'antd';
import {
  PlusOutlined,
  MoreOutlined,
  UserOutlined,
  FlagOutlined,
  ClockCircleOutlined,
} from '@ant-design/icons';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { issueApi, projectApi, boardApi } from '../../api';
import { Issue, Status, IssueType, Priority, Project } from '../../types';

const { Title, Text } = Typography;
const { TextArea } = Input;
const { Option } = Select;

const BoardPage: React.FC = () => {
  const [searchParams] = useSearchParams();
  const projectId = parseInt(searchParams.get('project_id') || '0');
  const [project, setProject] = useState<Project | null>(null);
  const [statuses, setStatuses] = useState<Status[]>([]);
  const [issueTypes, setIssueTypes] = useState<IssueType[]>([]);
  const [priorities, setPriorities] = useState<Priority[]>([]);
  const [boardData, setBoardData] = useState<Record<number, Issue[]>>({});
  const [loading, setLoading] = useState(false);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [selectedStatusId, setSelectedStatusId] = useState<number | null>(null);
  const [form] = Form.useForm();
  const navigate = useNavigate();

  const fetchProjectData = async () => {
    if (!projectId) return;
    try {
      const [proj, statusList, issueTypeList, priorityList] = await Promise.all([
        projectApi.getProject(projectId),
        projectApi.listStatuses(projectId),
        projectApi.listIssueTypes(projectId),
        projectApi.listPriorities(projectId),
      ]);
      setProject(proj);
      setStatuses(statusList);
      setIssueTypes(issueTypeList);
      setPriorities(priorityList);
    } catch (error) {
      console.error('Failed to fetch project data:', error);
    }
  };

  const fetchBoardData = async () => {
    if (!projectId) return;
    setLoading(true);
    try {
      const result = await issueApi.listIssues({ project_id: projectId, page_size: 100 });
      const grouped: Record<number, Issue[]> = {};
      statuses.forEach((s) => {
        grouped[s.id] = [];
      });
      result.list.forEach((issue) => {
        if (!grouped[issue.status_id]) {
          grouped[issue.status_id] = [];
        }
        grouped[issue.status_id].push(issue);
      });
      setBoardData(grouped);
    } catch (error) {
      console.error('Failed to fetch board data:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchProjectData();
  }, [projectId]);

  useEffect(() => {
    if (statuses.length > 0) {
      fetchBoardData();
    }
  }, [statuses]);

  const handleCreateIssue = async (values: any) => {
    try {
      await issueApi.createIssue({
        ...values,
        project_id: projectId,
        status_id: selectedStatusId || statuses[0]?.id,
      });
      message.success('任务创建成功');
      setIsModalVisible(false);
      form.resetFields();
      fetchBoardData();
    } catch (error) {
      console.error('Failed to create issue:', error);
    }
  };

  const openCreateModal = (statusId: number) => {
    setSelectedStatusId(statusId);
    setIsModalVisible(true);
  };

  const handleViewIssue = (issueKey: string) => {
    navigate(`/issues/${issueKey}`);
  };

  const getIssueTypeColor = (typeId: number) => {
    const type = issueTypes.find((t) => t.id === typeId);
    return type?.color || '#666';
  };

  const getIssueTypeName = (typeId: number) => {
    const type = issueTypes.find((t) => t.id === typeId);
    return type?.name || '未知';
  };

  const getPriorityColor = (priorityId?: number) => {
    if (!priorityId) return '#999';
    const priority = priorities.find((p) => p.id === priorityId);
    return priority?.color || '#666';
  };

  const getPriorityName = (priorityId?: number) => {
    if (!priorityId) return '未设置';
    const priority = priorities.find((p) => p.id === priorityId);
    return priority?.name || '未知';
  };

  if (!projectId) {
    return (
      <div style={{ textAlign: 'center', padding: '100px 0' }}>
        <Title level={3}>请选择一个项目</Title>
        <Text type="secondary">从项目列表选择一个项目查看看板</Text>
      </div>
    );
  }

  return (
    <div>
      <div style={{ marginBottom: 24, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div>
          <Title level={3} style={{ margin: 0 }}>{project?.name || '项目看板'}</Title>
          <Text type="secondary">{project?.key}</Text>
        </div>
        <Space>
          <Button onClick={() => navigate('/projects')}>返回项目列表</Button>
        </Space>
      </div>

      <div
        style={{
          display: 'flex',
          gap: 16,
          overflowX: 'auto',
          paddingBottom: 16,
        }}
      >
        {statuses.map((status) => (
          <div
            key={status.id}
            style={{
              minWidth: 300,
              width: 300,
              flexShrink: 0,
              background: '#f5f5f5',
              borderRadius: 8,
              padding: 12,
            }}
          >
            <div
              style={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                marginBottom: 12,
              }}
            >
              <Space>
                <Tag color={status.color} style={{ margin: 0 }}>
                  {status.name}
                </Tag>
                <Text type="secondary" style={{ fontSize: 12 }}>
                  {boardData[status.id]?.length || 0}
                </Text>
              </Space>
              <Button
                type="text"
                icon={<PlusOutlined />}
                size="small"
                onClick={() => openCreateModal(status.id)}
              />
            </div>

            <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
              {boardData[status.id]?.map((issue) => (
                <Card
                  key={issue.id}
                  size="small"
                  hoverable
                  onClick={() => handleViewIssue(issue.issue_key)}
                  style={{ cursor: 'pointer' }}
                  bodyStyle={{ padding: 12 }}
                >
                  <div style={{ marginBottom: 8 }}>
                    <Text type="secondary" style={{ fontSize: 12 }}>
                      {issue.issue_key}
                    </Text>
                  </div>
                  <div
                    style={{
                      fontWeight: 500,
                      marginBottom: 8,
                      wordBreak: 'break-word',
                    }}
                  >
                    {issue.summary}
                  </div>
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <Space size={8}>
                      <Tag
                        color={getIssueTypeColor(issue.issue_type_id)}
                        style={{ margin: 0, fontSize: 11 }}
                      >
                        {getIssueTypeName(issue.issue_type_id)}
                      </Tag>
                      {issue.priority_id && (
                        <FlagOutlined style={{ color: getPriorityColor(issue.priority_id) }} />
                      )}
                    </Space>
                    <Space size={8}>
                      {issue.assignee_id && (
                        <span style={{ fontSize: 12, color: '#888' }}>
                          <UserOutlined />
                        </span>
                      )}
                    </Space>
                  </div>
                </Card>
              ))}
            </div>
          </div>
        ))}
      </div>

      <Modal
        title="创建任务"
        open={isModalVisible}
        onCancel={() => {
          setIsModalVisible(false);
          form.resetFields();
        }}
        footer={null}
        width={600}
      >
        <Form form={form} layout="vertical" onFinish={handleCreateIssue}>
          <Form.Item
            name="summary"
            label="标题"
            rules={[
              { required: true, message: '请输入任务标题' },
              { max: 255, message: '最多255个字符' },
            ]}
          >
            <Input placeholder="请输入任务标题" />
          </Form.Item>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="issue_type_id"
                label="任务类型"
                rules={[{ required: true, message: '请选择任务类型' }]}
              >
                <Select placeholder="请选择">
                  {issueTypes.map((type) => (
                    <Option key={type.id} value={type.id}>
                      <Tag color={type.color}>{type.name}</Tag>
                    </Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="priority_id" label="优先级">
                <Select placeholder="请选择" allowClear>
                  {priorities.map((p) => (
                    <Option key={p.id} value={p.id}>
                      <span style={{ color: p.color }}>●</span> {p.name}
                    </Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
          </Row>

          <Form.Item name="description" label="描述">
            <TextArea rows={4} placeholder="请输入任务描述" />
          </Form.Item>

          <Form.Item style={{ marginBottom: 0 }}>
            <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
              <Button
                onClick={() => {
                  setIsModalVisible(false);
                  form.resetFields();
                }}
              >
                取消
              </Button>
              <Button type="primary" htmlType="submit">
                创建
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default BoardPage;
