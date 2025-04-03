import { Typography, Card, Row, Col, Statistic } from 'antd';
import { UserOutlined, ShoppingCartOutlined, DollarOutlined } from '@ant-design/icons';

const { Title } = Typography;

const Home = () => {
  return (
    <div>
      <Title level={2}>Dashboard</Title>
      
      <Row gutter={16}>
        <Col span={8}>
          <Card>
            <Statistic
              title="Users"
              value={1128}
              prefix={<UserOutlined />}
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic
              title="Orders"
              value={93}
              prefix={<ShoppingCartOutlined />}
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic
              title="Revenue"
              value={11280}
              prefix={<DollarOutlined />}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={16} style={{ marginTop: 16 }}>
        <Col span={12}>
          <Card title="Recent Activity">
            <p>No recent activity</p>
          </Card>
        </Col>
        <Col span={12}>
          <Card title="Quick Actions">
            <p>No quick actions available</p>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default Home; 