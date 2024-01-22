import React from "react";
import {Tag} from "antd";
import {PageContainer} from "@ant-design/pro-components";
import {ProList} from '@ant-design/pro-components';

const data = [
  {
    title: "HTTP",
    subTitle: <Tag color="#5BD8A6">语雀专栏</Tag>,
    actions: [<a key="detail">详情</a>, <a key="stop">停止</a>, <a key={"del"}>删除</a>],
    avatar:
      'https://gw.alipayobjects.com/zos/antfincdn/UCSiy1j6jx/xingzhuang.svg',
    content: "",
  },
];

const Proxy: React.FC = () => {
  return <PageContainer
    title={"代理"}
    content={`通过ProList展示代理信息，包括:
    `}
    extra={""}>
    <div
      style={{
        backgroundColor: '#eee',
        margin: -24,
        padding: 24,
      }}
    >
    - proxy <br />
      - proxy_list 列表 <br />
      - proxy_type 类别（http/s5/auto) <br />
      - proxy_node 节点 <br />
      - proxy_addr 地址 <br />
      - proxy_port 端口 <br />
      - proxy_status 状态展示 health_check + sdk get status <br />
      - proxy_speed 速率展示 <br />
      - proxy_conn 连接管理 <br />
      - proxy_auth 认证管理 <br />
      <ProList<any>
        ghost={true}
        itemCardProps={{
          ghost: true,
        }}
        pagination={{
          defaultPageSize: 8,
          showSizeChanger: false,
        }}
        showActions="hover"
        rowSelection={{}}
        grid={{gutter: 16, column: 3}}
        onItem={(record: any) => {
          return {
            onMouseEnter: () => {
              console.log(record);
            },
            onClick: () => {
              console.log(record);
            },
          };
        }}
        metas={{
          title: {},
          subTitle: {},
          type: {},
          avatar: {},
          content: {},
          actions: {
            cardActionProps: 'extra',
          },
        }}
        headerTitle="卡片列表展示"
        dataSource={data}
      />
    </div>
  </PageContainer>
}

export default Proxy

