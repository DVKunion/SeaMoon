import {PageContainer, ProCard, Statistic, StatisticCard} from '@ant-design/pro-components';
import React, {useEffect, useRef, useState} from 'react';
import RcResizeObserver from 'rc-resize-observer';
import {Badge, Progress, Space, Tag} from "antd";
import {Line} from '@ant-design/plots';
import styles from "./index.less";
import IconFont from "@/components/IconFont";

const {Divider} = StatisticCard;
const colors = [
  "#28a745", "#32ad46", "#3cb548", "#46bd49", "#51c44b",
  "#5bcc4c", "#66d34d", "#70db4f", "#7ae351", "#84ea52",
  "#8ff254", "#99fa55", "#a3fc57", "#adfe5b", "#b7ff5f",
  "#c1ff63", "#cbff68", "#d5ff6c", "#dfff70", "#e9ff74",
  "#f2f276", "#f3e671", "#f4da6d", "#f5cd69", "#f6c065",
  "#f7b361", "#f8a75d", "#f99b59", "#fa8f55", "#fb8351",
  "#fc774d", "#fd6b49", "#fe5f45", "#ff5341", "#ff473d",
  "#ff3b39", "#ff2f35", "#ff2331", "#ff172d", "#ff0b29",
  "#f8092a", "#f1082b", "#ea072d", "#e3062e", "#dc0530",
  "#d50431", "#ce0333", "#c70234", "#c00136", "#dc3545"
];

type NetworkProps = {
  values: {
    time: number,
    up: number,
    down: number,
  }[]
}

function formatTimestamp(timestamp: number) {
  const date = new Date(timestamp);
  const hours = date.getHours().toString().padStart(2, '0');
  const minutes = date.getMinutes().toString().padStart(2, '0');
  const seconds = date.getSeconds().toString().padStart(2, '0');

  return `${hours}:${minutes}:${seconds}`;
}

function formatTimestampDate(timestamp: number) {
  const date = new Date(timestamp);

  const year = date.getFullYear(); // 获取两位数的年份
  const month = (date.getMonth() + 1).toString().padStart(2, '0'); // 获取月份（从0开始，所以要加1）
  const day = date.getDate().toString().padStart(2, '0'); // 获取日期

  const hours = date.getHours().toString().padStart(2, '0'); // 获取小时
  const minutes = date.getMinutes().toString().padStart(2, '0'); // 获取分钟
  const seconds = date.getSeconds().toString().padStart(2, '0'); // 获取秒数

  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
}

const Networks: React.FC<NetworkProps> = (props) => {
  const config = {
    height: 230,
    data: {
      value: props.values,
      transform: [
        {
          type: 'fold',
          fields: ['up', 'down'],
          key: 'type',
          value: 'value',
        },
      ],
    },
    xField: (d: any) => formatTimestamp(d.time),
    yField: 'value',
    colorField: 'type',
    legend: {
      color: {
        itemLabelFill: 'rgba(229, 224, 216, 0.85)'
      }
    },
    axis: {
      x: {
        // labelAutoHide: 'greedy',
        labelFill: 'rgba(229, 224, 216, 0.85)'
      },
      y: {
        // labelAutoHide: 'greedy',
        labelFill: 'rgba(229, 224, 216, 0.85)'
      }
    },
  }

  return <Line {...config} />;
};


const Index: React.FC = () => {
  const [responsive, setResponsive] = useState(false);
  const [trafficData, setTrafficData] = useState([{time: 0, up: 0, down: 0}]);
  const dataRef = useRef([{time: 0, up: 0, down: 0}]);
  dataRef.current = trafficData;

  const now = new Date().getTime();
  useEffect(() => {
    const interval = setInterval(() => {
      const newData = {
        time: new Date().getTime(),
        up: Math.floor(Math.random() * 100),
        down: Math.floor(Math.random() * 100),
      };
      setTrafficData((prevData) => [...prevData.slice(-19), newData]); // 保持数据长度为20
    }, 2000); // 每秒更新一次数据

    return () => clearInterval(interval); // 组件卸载时清除定时器
  }, []);

  // useEffect(() => {
  //   // 创建XMLHttpRequest对象
  //   const xhr = new XMLHttpRequest();
  //
  //   // 配置请求
  //   xhr.open('GET', 'http://127.0.0.1:9090/traffic', true);
  //   xhr.setRequestHeader('Transfer-Encoding', 'chunked');
  //
  //   let partialData = '';
  //
  //   // 监听onreadystatechange事件
  //   xhr.onreadystatechange = () => {
  //     if (xhr.readyState === 3) { // readyState 3表示请求仍在处理中
  //       partialData += xhr.responseText;
  //
  //       // 处理响应数据，找到完整的JSON对象
  //       let boundary = partialData.indexOf('}');
  //       while (boundary !== -1) {
  //         const jsonString = partialData.slice(0, boundary + 1);
  //         partialData = partialData.slice(boundary + 1);
  //
  //         try {
  //           const data = JSON.parse(jsonString);
  //           const time = Date.now();
  //           const traffic = { time, ...data };
  //
  //           setTrafficData(prevData => {
  //             // 新数组包括新数据
  //             const updatedData = [...prevData, traffic];
  //
  //             // 如果长度超过15，删除最早的数据
  //             if (updatedData.length > 15) {
  //               updatedData.shift();
  //             }
  //             console.log("update data: ",updatedData)
  //             return updatedData;
  //           });
  //         } catch (e) {
  //           console.error('Error parsing JSON:', e);
  //         }
  //
  //         boundary = partialData.indexOf('}');
  //       }
  //     }
  //   };
  //
  //   // 发送请求
  //   xhr.send();
  //
  //   // 清理函数：组件卸载时中止请求
  //   return () => {
  //     xhr.abort();
  //   };
  // }, []);

  return (
    <PageContainer>
      <RcResizeObserver
        key="resize-observer"
        onResize={(offset) => {
          setResponsive(offset.width < 596);
        }}
      >
        <Badge.Ribbon color="cyan" text={formatTimestampDate(now)}>
          <ProCard
            title="资源统计"
            split={responsive ? 'horizontal' : 'vertical'}
            gutter={8}
          >
            <ProCard split="horizontal">
              <ProCard split={"vertical"} ghost>
                <StatisticCard
                  statistic={{
                  title: 'CPU',
                  value: '4',
                  suffix: '/ 8U'
                }} chart={
                  <Progress percent={50} steps={50} strokeColor={colors} size={"small"} format={(percent) => {
                    if (percent != undefined) {
                      return <div color={colors[percent / 2 - 1]}>{percent} %</div>
                    }
                    return <div>{percent} %</div>
                  }}/>
                }/>
                <StatisticCard statistic={{
                  title: 'Mem',
                  value: '1',
                  suffix: '/ 16G'
                }} chart={
                  <Progress percent={10} steps={50} strokeColor={colors} size={"small"} format={(percent) => {
                    if (percent != undefined) {
                      return <div color={colors[percent / 2 - 1]}>{percent} %</div>
                    }
                    return <div>{percent} %</div>
                  }}/>
                }/>
              </ProCard>
              <ProCard split={"vertical"}>
                <StatisticCard statistic={{
                  title: 'Swap',
                  value: '0',
                  suffix: '/ 8G'
                }} chart={
                  <Progress percent={0} steps={50} strokeColor={colors} size={"small"} format={(percent) => {
                    if (percent != undefined) {
                      return <div color={colors[percent / 2 - 1]}>{percent} %</div>
                    }
                    return <div>{percent} %</div>
                  }}/>
                }/>
                <StatisticCard statistic={{
                  title: 'Disk',
                  value: '100',
                  suffix: '/ 256G'
                }} chart={
                  <Progress percent={70} steps={50} strokeColor={colors} size={"small"} format={(percent) => {
                    if (percent != undefined) {
                      return <div color={colors[percent / 2 - 1]}>{percent} %</div>
                    }
                    return <div>{percent} %</div>
                  }}/>
                }/>
              </ProCard>
            </ProCard>
            <ProCard split="horizontal" title={"网络流量"}>
              <Networks values={trafficData}/>
            </ProCard>
          </ProCard>
        </Badge.Ribbon>
      </RcResizeObserver>
      <Divider type={responsive ? 'horizontal' : 'vertical'} />
      <StatisticCard.Group title={"账户统计"} direction={responsive ? 'column' : 'row'} >
        <StatisticCard
          statistic={{
            title: '云账户数',
            value: 3,
            icon: <IconFont type={"icon-cloud1"} style={{fontSize: "200%"}}/>,
          }}
        />
        <StatisticCard
          statistic={{
            title: '运行服务',
            value: 3,
            icon: <IconFont type={"icon-fuwu"} style={{fontSize: "200%"}}/>,
          }}
        />
        <StatisticCard
          statistic={{
            title: '部署函数',
            value: 900,
            icon: <IconFont type={"icon-fchanshujisuan"} style={{fontSize: "200%"}}/>,
          }}
        />
        <StatisticCard
          statistic={{
            title: '花费总计',
            value: 900,
            icon: <IconFont type={"icon-feiyong"} style={{fontSize: "200%"}}/>
          }}
        />
      </StatisticCard.Group>
      <Divider type={responsive ? 'horizontal' : 'vertical'}/>
      <ProCard split={responsive ? 'horizontal' : 'vertical'} style={{marginBlockStart: 8}}
               gutter={4} ghost >
        <ProCard
          colSpan="25%"
          className={styles.status}
          title={"SeaMoon"}
          extra={<Tag color="cyan">{"1.1.0"}</Tag>}
        >
          <Statistic className={styles.description} title={"当前状态"} valueRender={(v) => {
            return <Space><Badge status={"success"}/><span>{"success"}</span></Space>
          }}/>
          <Statistic className={styles.description} title="启动时长" value="30" suffix={"d"}/>
          <Statistic className={styles.description} title="资源占用" value=""/>
        </ProCard>
        <ProCard
          colSpan="25%"
          className={styles.status}
          title={"Xray"}
          extra={<Tag color="cyan">{"1.9.0"}</Tag>}
        >
          <Statistic className={styles.description} title={"当前状态"} valueRender={(v) => {
            return <Space><Badge status={"success"}/><span>{"success"}</span></Space>
          }}/>
          <Statistic className={styles.description} title="启动时长" value="30" suffix={"d"}/>
          <Statistic className={styles.description} title="资源占用" value=""/>
        </ProCard>
      </ProCard>

    </PageContainer>
  );
};

export default Index;
