import { Form, Input, Button, Checkbox, Card, Typography } from 'antd';
import { GoogleOutlined } from '@ant-design/icons';
import { useNavigate, Link } from 'react-router-dom';
import styles from './Login.module.css';

const { Title, Text } = Typography;

const Login = () => {
  const navigate = useNavigate();

  const onFinish = (values: any) => {
    console.log('Success:', values);
    // TODO: Implement login logic
    navigate('/');
  };

  return (
    <div className={styles.container}>
      <div className={styles.logoContainer}>
        <img src="/logo.svg" alt="InsideBox" className={styles.logo} />
        <Title level={4} style={{ marginBottom: 0 }}>InsideBox</Title>
      </div>
      
      <Card bordered={false} className={styles.card}>
        <Text type="secondary" style={{ display: 'block', marginBottom: 8 }}>
          Please enter your details
        </Text>
        <Title level={2} style={{ marginBottom: 24 }}>Welcome back</Title>

        <Form
          name="login"
          onFinish={onFinish}
          autoComplete="off"
          layout="vertical"
          requiredMark={false}
        >
          <Form.Item
            name="email"
            label="Email adress"
            rules={[{ required: true, message: 'Please input your email!' }]}
          >
            <Input size="large" placeholder="Enter your email" />
          </Form.Item>

          <Form.Item
            name="password"
            label="Password"
            rules={[{ required: true, message: 'Please input your password!' }]}
          >
            <Input.Password size="large" placeholder="Enter your password" />
          </Form.Item>

          <div className={styles.rememberForgot}>
            <Form.Item name="remember" valuePropName="checked" noStyle>
              <Checkbox>Remember for 30 days</Checkbox>
            </Form.Item>
            <Link to="/forgot-password" className={styles.forgotLink}>
              Forgot password
            </Link>
          </div>

          <Form.Item>
            <Button type="primary" htmlType="submit" size="large" block>
              Sign in
            </Button>
          </Form.Item>

          <Button 
            icon={<GoogleOutlined />} 
            size="large" 
            block
            className={styles.googleButton}
          >
            Sign in with Google
          </Button>

          <div className={styles.signupContainer}>
            <Text type="secondary">Don't have an account?</Text>
            <Link to="/register" className={styles.signupLink}>
              Sign up
            </Link>
          </div>
        </Form>
      </Card>
    </div>
  );
};

export default Login; 