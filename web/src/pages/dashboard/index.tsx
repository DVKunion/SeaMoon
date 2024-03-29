import {PageContainer} from '@ant-design/pro-components';
import React, {useState} from 'react';
import {StatisticCard} from '@ant-design/pro-components';
import RcResizeObserver from 'rc-resize-observer';
import IconFont from "@/components/IconFont";

// const {Statistic} = StatisticCard;

const Index: React.FC = () => {
  const [responsive, setResponsive] = useState(false);

  return (
    <PageContainer
      title={"数据统计"}
    >
      <RcResizeObserver
        key="resize-observer"
        onResize={(offset) => {
          setResponsive(offset.width < 596);
        }}
      >
        <StatisticCard.Group title={"系统信息"} direction={responsive ? 'column' : 'row'}>
          <StatisticCard
            statistic={{
              title: '隧道数',
              value: 3,
              icon: <IconFont type={"icon-ico_tunnel-A"} style={{fontSize: "300%"}}/>,
            }}
          />
          <StatisticCard
            statistic={{
              title: '代理服务数',
              value: 10,
              icon: <IconFont type={"icon-PROXY-A"} style={{fontSize: "300%"}}/>,
            }}
          />
          <StatisticCard
            statistic={{
              title: '云账户数',
              value: 3,
              icon: <IconFont type={"icon-cloud1"} style={{fontSize: "300%"}}/>,
            }}
          />
          <StatisticCard
            statistic={{
              title: '花费总计',
              value: 900,
              icon: <IconFont type={"icon-feiyong"} style={{fontSize: "300%"}}/>,
            }}
          />
        </StatisticCard.Group>
        <br/>
      </RcResizeObserver>
    </PageContainer>
  );
};

export default Index;
