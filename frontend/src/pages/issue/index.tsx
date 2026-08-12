import React, { useEffect, useState } from 'react';
import {
  Row,
  Col,
  Typography,
  Tag,
  Button,
  Space,
  Card,
  Avatar,
  Input,
  List,
  Comment,
  Tabs,
  Descriptions,
  Modal,
  Form,
  Select,
  message,
} from 'antd';
import {
  UserOutlined,
  FlagOutlined,
  ClockCircleOutlined,
  SendOutlined,
  EditOutlined,
  ArrowLeftOutlined,
} from '@ant-design/icons';
import { useParams, useNavigate } from 'react-router-dom';
import { issueApi, projectApi } from '../../api';
import { Issue, IssueComment, IssueChangelog, Status, IssueType, Priority } from '../../types';
import dayjs from 'dayjs';

const { Title, Text } = Typography;
const { TextArea } = Input;
const { Option } = Select;

const IssueDetailPage: React.FC = () => {
  const { issueKey } = useParams<{ issueKey: string }>();
  const navigate = useNavigate();
  const [issue, setIssue] = useState<Issue | null>(null);
  const [comments, setComments] = useState<IssueComment[]>([]);
  const [changelogs, setChangelogs] = useState<IssueChangelog[]>([]);
  const [statuses, setStatuses] = useState<Status[]>([]);
  const [issueTypes, setIssueTypes] = useState<IssueType[]>([]);
  const [priorities, setPriorities] = useState<Priority[]>([]);
  const [commentText, setCommentText] = useState('');
  const [editModalVisible, setEditModalVisible] = useState(false);
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);

  const fetchIssue = async () => {
    if (!issueKey) return;
    try {
      const data = await issueApi.getIssueByKey(issueKey);
      setIssue(data);

      const [c, logs] = await Promise.all([
        issueApi.getComments(data.id),
        issueApi.getChangelogs(data.id),
        projectApi.listStatuses(data.project_id),
        projectApi.listIssueTypes(data.project_id),
        projectApi.listPriorities(data.project_id),
      ]);
      setComments(c);
      setChangelogs(logs);
    } catch (error) {
      console.error('Failed to fetch issue:', error);
    }
  };

  const fetchProjectConfig = async (projectId: number) => {
    try {
      const [s, it, p] = await Promise.all([
        projectApi.listStatuses(projectId),
        projectApi.listIssueTypes(projectId),
        projectApi.listPriorities(projectId),
      ]);
      setStatuses(s);
      setIssueTypes(it);
      setPriorities(p);
    } catch (error) {
      console.error('Failed to fetch project config:', error);
    }
  };

  useEffect(() => {
    fetchIssue();
  }, [issueKey]);

  useEffect(() => {
    if (issue) {
      fetchProjectConfig(issue.project_id);
    }
  }, [issue]);

  const handleAddComment = async () => {
    if (!commentText.trim() || !issue) return;
    try {
      await issueApi.addComment(issue.id, commentText);
      setCommentText('');
      fetchIssue();
      message.success('评论发表成功');
    } catch (error) {
      console.error('Failed to add comment:', error);
    }
  };

  const handleEditIssue = async (values: any) => {
    if (!issue) return;
    setLoading(true);
    try {
      await issueApi.updateIssue(issue.id, values);
      message.success('更新成功');
      setEditModalVisible(false);
      fetchIssue();
    } catch (error) {
      console.error('Failed to update issue:', error);
    } finally {
      setLoading(false);
    }
  };

  const getStatusName = (statusId: number) => {
    const status = statuses.find((s) => s.id === statusId);
    return status?.name || '未知';
  };

  const getStatusColor = (statusId: number) => {
    const status = statuses.find((s) => s.id === statusId);
    return status?.color || '#666';
  };

  const getIssueTypeName = (typeId: number) => {
    const type = issueTypes.find((t) => t.id === typeId);
    return type?.name || '未知';
  };

  const getIssueTypeColor = (typeId: number) => {
    const type = issueTypes.find((t) => t.id === typeId);
    return type?.color || '#666';
  };

  const getPriorityName = (priorityId?: number) => {
    if (!priorityId) return '未设置';
    const priority = priorities.find((p) => p.id === priorityId);
    return priority?.name || '未知';
  };

  const getPriorityColor = (priorityId?: number) => {
    if (!priorityId) return '#999';
    const priority = priorities.find((p) => p.id === priorityId);
    return priority?.color || '#666';
  };

  const tabItems = [
    {
      key: 'comments',
      label: `评论 (${comments.length})`,
      children: (
        <div>
          <div style={{ marginBottom: 16 }}>
            <TextArea
              rows={3}
              placeholder="添加评论..."
              value={commentText}
              onChange={(e) => setCommentText(e.target.value)}
              style={{ marginBottom: 8 }}
            />
            <div style={{ textAlign: 'right' }}>
              <Button
                type="primary"
                icon={<SendOutlined />}
                onClick={handleAddComment}
                disabled={!commentText.trim()}
              >
                发表评论
              </Button>
            </div>
          </div>
          <List
            dataSource={comments}
            renderItem={(item) => (
              <Comment
                author={<Text strong>用户 {item.author_id}</Text>}
                avatar={<Avatar icon={<UserOutlined />} />}
                content={item.content}
                datetime={dayjs(item.created_at).format('YYYY-MM-DD HH:mm')}
              />
            )}
          />
        </div>
      ),
    },
    {
      key: 'activity',
      label: `活动 (${changelogs.length})`,
      children: (
        <List
          dataSource={changelogs}
          renderItem={(item) => (
            <List.Item>
              <List.Item.Meta
                avatar={<Avatar icon={<UserOutlined />} />}
                title={
                  <Space>
                    <Text strong>用户 {item.author_id}</Text>
                    <Text type="secondary">修改了 {item.field}</Text>
                  </Space>
                }
                description={
                  <div>
                    <Text delete type="secondary">{item.old_value || '空'}</Text>
                    <Text style={{ margin: '0 8px' }}>→</Text>
                    <Text>{item.new_value || '空'}</Text>
                    <div style={{ marginTop: 4 }}>
                      <Text type="secondary" style={{ fontSize: 12 }}>
                        <ClockCircleOutlined /> {dayjs(item.created_at).format('YYYY-MM-DD HH:mm')}
                      </Text>
                    </div>
                  </div>
                }
              />
            </List.Item>
          )}
        />
      ),
    },
  ];

  if (!issue) {
    return <div style={{ padding: 50, textAlign: 'center' }}>加载中...</div>;
  }

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <Button
          type="text"
          icon={<ArrowLeftOutlined />}
          onClick={() => navigate(-1)}
          style={{ paddingLeft: 0 }}
        >
          返回
        </Button>
      </div>

      <Row gutter={24}>
        <Col span={17}>
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
            <div>
              <Tag color={getIssueTypeColor(issue.issue_type_id)} style={{ marginBottom: 8 }}>
                {getIssueTypeName(issue.issue_type_id)}
              </Tag>
              <Title level={3} style={{ margin: 0 }}>{issue.summary}</Title>
              <Text type="secondary">{issue.issue_key}</Text>
            </div>
            <Button
              type="primary"
              icon={<EditOutlined />}
              onClick={() => {
                form.setFieldsValue({
                  summary: issue.summary,
                  description: issue.description,
                  issue_type_id: issue.issue_type_id,
                  status_id: issue.status_id,
                  priority_id: issue.priority_id,
                  assignee_id: issue.assignee_id,
                });
                setEditModalVisible(true);
              }}
            >
              编辑
            </Button>
          </div>

          <Card style={{ marginBottom: 16 }}>
            <div style={{ marginBottom: 12 }}>
              <Text strong style={{ fontSize: 16 }}>描述</Text>
            </div>
            <div style={{ whiteSpace: 'pre-wrap', color: '#333' }}>
              {issue.description || '暂无描述'}
            </div>
          </Card>

          <Card>
            <Tabs items={tabItems} defaultActiveKey="comments" />
          </Card>
        </Col>

        <Col span={7}>
          <Card style={{ marginBottom: 16 }}>
            <Descriptions column={1} size="small">
              <Descriptions.Item label="状态">
                <Tag color={getStatusColor(issue.status_id)}>
                  {getStatusName(issue.status_id)}
                </Tag>
              </Descriptions.Item>
              <Descriptions.Item label="优先级">
                <span style={{ color: getPriorityColor(issue.priority_id) }}>
                  <FlagOutlined /> {getPriorityName(issue.priority_id)}
                </span>
              </Descriptions.Item>
              <Descriptions.Item label="报告人">
                <Space>
                  <Avatar size="small" icon={<UserOutlined />} />
                  <span>用户 {issue.reporter_id}</span>
                </Space>
              </Descriptions.Item>
              <Descriptions.Item label="经办人">
                {issue.assignee_id ? (
                  <Space>
                    <Avatar size="small" icon={<UserOutlined />} />
                    <span>用户 {issue.assignee_id}</span>
                  </Space>
                ) : (
                  <Text type="secondary">未分配</Text>
                )}
              </Descriptions.Item>
              {issue.due_date && (
                <Descriptions.Item label="截止日期">
                  <ClockCircleOutlined /> {dayjs(issue.due_date).format('YYYY-MM-DD')}
                </Descriptions.Item>
              )}
              {issue.story_points !== undefined && issue.story_points !== null && (
                <Descriptions.Item label="故事点">
                  {issue.story_points}
                </Descriptions.Item>
              )}
            </Descriptions>
          </Card>

          <Card>
            <Descriptions column={1} size="small">
              <Descriptions.Item label="创建时间">
                {dayjs(issue.created_at).format('YYYY-MM-DD HH:mm')}
              </Descriptions.Item>
              <Descriptions.Item label="更新时间">
                {dayjs(issue.updated_at).format('YYYY-MM-DD HH:mm')}
              </Descriptions.Item>
            </Descriptions>
          </Card>
        </Col>
      </Row>

      <Modal
        title="编辑任务"
        open={editModalVisible}
        onCancel={() => setEditModalVisible(false)}
        footer={null}
        width={600}
      >
        <Form form={form} layout="vertical" onFinish={handleEditIssue}>
          <Form.Item
            name="summary"
            label="标题"
            rules={[{ required: true, message: '请输入标题' }]}
          >
            <Input />
          </Form.Item>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="status_id" label="状态">
                <Select>
                  {statuses.map((s) => (
                    <Option key={s.id} value={s.id}>
                      <Tag color={s.color}>{s.name}</Tag>
                    </Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="priority_id" label="优先级">
                <Select allowClear>
                  {priorities.map((p) => (
                    <Option key={p.id} value={p.id}>
                      <span style={{ color: p.color }}>●</span> {p.name}
                    </Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="issue_type_id" label="任务类型">
                <Select>
                  {issueTypes.map((t) => (
                    <Option key={t.id} value={t.id}>
                      <Tag color={t.color}>{t.name}</Tag>
                    </Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
          </Row>

          <Form.Item name="description" label="描述">
            <TextArea rows={4} />
          </Form.Item>

          <Form.Item style={{ marginBottom: 0 }}>
            <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
              <Button onClick={() => setEditModalVisible(false)}>取消</Button>
              <Button type="primary" htmlType="submit" loading={loading}>
                保存
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default IssueDetailPage;
