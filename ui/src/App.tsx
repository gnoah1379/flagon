import { ConfigProvider } from 'antd';
import { BrowserRouter as Router } from 'react-router-dom';
import AppRoutes from './routes';
import '@ant-design/v5-patch-for-react-19';

function App() {
  return (
    <ConfigProvider
      theme={{
        token: {
          colorPrimary: '#1677ff',
        },
      }}
    >
      <Router>
        <AppRoutes />
      </Router>
    </ConfigProvider>
  );
}

export default App;
