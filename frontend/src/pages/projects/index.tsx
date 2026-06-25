import React, { useEffect, useState } from 'react';
import {
  Table,
  Button,
  Space,
  Modal,
  Form,
  Input,
  Select,
  Tag,
  message,
  Card,
  Row,
  Col,
  Typography,
} from 'antd';
import {
  PlusOutlined,
  ProjectOutlined,
  TeamOutlined,
  CalendarOutlined,
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { projectApi } from '../../api';
import { Project } from '../../types';

const { Title } = Typography;
const { TextArea } = Input;

const ProjectsPage: React.FC = () => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [form] = Form.useForm();
  const navigate = useNavigate();

  const fetchProjects = async () => {
    setLoading(true);
    try {
      const result = await projectApi.listProjects({ page, page_size: pageSize });
      setProjects(result.list);
      setTotal(result.total);
    } catch (error) {
      console.error('Failed to fetch projects:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchProjects();
  }, [page, pageSize]);

  const handleCreateProject = async (values: any) => {
    try {
      await projectApi.createProject(values);
      message.success('项目创建成功');
      setIsModalVisible(false);
      form.resetFields();
      fetchProjects();
    } catch (error) {
      console.error('Failed to create project:', error);
    }
  };

  const handleViewProject = (project: Project) => {
    navigate(`/board?project_id=${project.id}`);
  };

  const columns = [
    {
      title: '项目',
      dataIndex: 'name',
      key: 'name',
      render: (text: string, record: Project) => (
        <Space>
          <ProjectOutlined style={{ color: '#1890ff', fontSize: 18 }} />
          <div>
            <div style={{ fontWeight: 500, cursor: 'pointer' }} onClick={() => handleViewProject(record)}>
              {text}
            </div>
            <div style={{ color: '#888', fontSize: 12 }}>{record.key}</div>
          </div>
        </Space>
      ),
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const colorMap: Record<string, string> = {
          active: 'green',
          archived: 'default',
          deleted: 'red',
        };
        return <Tag color={colorMap[status] || 'default'}>{status}</Tag>;
      },
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => new Date(date).toLocaleDateString(),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Project) => (
        <Space>
          <Button type="link" onClick={() => handleViewProject(record)}>
            查看
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <div style={{ marginBottom: 24, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Title level={3} style={{ margin: 0 }}>项目列表</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setIsModalVisible(true)}>
          创建项目
        </Button>
      </div>

      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col span={8}>
          <Card>
            <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
              <ProjectOutlined style={{ fontSize: 28, color: '#1890ff' }} />
              <div>
                <div style={{ fontSize: 24, fontWeight: 'bold' }}>{total}</div>
                <div style={{ color: '#888' }}>项目总数</div>
              </div>
            </div>
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
              <TeamOutlined style={{ fontSize: 28, color: '#52c41a' }} />
              <div>
                <div style={{ fontSize: 24, fontWeight: 'bold' }}>-</div>
                <div style={{ color: '#888' }}>团队成员</div>
              </div>
            </div>
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
              <CalendarOutlined style={{ fontSize: 28, color: '#faad14' }} />
              <div>
                <div style={{ fontSize: 24, fontWeight: 'bold' }}>-</div>
                <div style={{ color: '#888' }}>活跃冲刺</div>
              </div>
            </div>
          </Card>
        </Col>
      </Row>

      <Table
        columns={columns}
        dataSource={projects}
        rowKey="id"
        loading={loading}
        pagination={{
          current: page,
          pageSize,
          total,
          onChange: (p, ps) => {
            setPage(p);
            setPageSize(ps);
          },
          showSizeChanger: true,
          showTotal: (t) => `共 ${t} 个项目`,
        }}
      />

      <Modal
        title="创建新项目"
        open={isModalVisible}
        onCancel={() => {
          setIsModalVisible(false);
          form.resetFields();
        }}
        footer={null}
        width={500}
      >
        <Form form={form} layout="vertical" onFinish={handleCreateProject}>
          <Form.Item
            name="key"
            label="项目键"
            rules={[
              { required: true, message: '请输入项目键' },
              { min: 2, message: '至少2个字符' },
              { max: 20, message: '最多20个字符' },
              { pattern: /^[A-Z0-9_]+$/, message: '只能包含大写字母、数字和下划线' },
            ]}
          >
            <Input placeholder="例如：PROJ" style={{ textTransform: 'uppercase' }} />
          </Form.Item>
          <Form.Item
            name="name"
            label="项目名称"
            rules={[
              { required: true, message: '请输入项目名称' },
              { max: 100, message: '最多100个字符' },
            ]}
          >
            <Input placeholder="请输入项目名称" />
          </Form.Item>
          <Form.Item name="description" label="项目描述">
            <TextArea rows={3} placeholder="请输入项目描述" />
          </Form.Item>
          <Form.Item style={{ marginBottom: 0 }}>
            <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
              <Button onClick={() => setIsModalVisible(false)}>取消</Button>
              <Button type="primary" htmlType="submit">创建</Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default ProjectsPage;
