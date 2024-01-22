import {PageContainer} from '@ant-design/pro-components';
import React from 'react';


const DashBoard: React.FC = () => {
  return (
    <PageContainer
    >
      这里是仪表盘，每个 card 展示数据统计信息，如：实时流量信息、流量访问统计、汇总、客户端连接个数、状态等等。
    </PageContainer>
  );
};

export default DashBoard;
