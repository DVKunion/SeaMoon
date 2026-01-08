import React, {useRef, useState} from "react";
import {ActionType, PageContainer, ProList, StatisticCard, ProCard} from "@ant-design/pro-components";
import {Badge, Button, Space, Tag, Tooltip} from "antd";
import {PlusOutlined} from "@ant-design/icons";
import IconFont from "@/components/IconFont";
import CreateForm from "./components/CreateForm";
import DetailDrawer from "./components/DetailDrawer";
import {getServerlessTunnel} from "@/services/function/api";
import {CloudProvideTypeValueEnum, RegionEnum} from "@/enum/cloud";
import {TunnelStatusEnum, TunnelTypeValueEnum} from "@/enum/tunnel";
import styles from "./index.less";
import {handleCreateTunnel, handleDeleteTunnel, handleUpdateTunnel} from "@/pages/function/handle";
import {ClashDropDown, ShadowRocketDropDown, SSrDropDown} from "@/components/Subcirber";

const {Statistic} = StatisticCard;

const Tunnel: React.FC = () => {

  const actionRef = useRef<ActionType>();
  const [detailModalVisible, setDetailVisible] = useState<boolean>(false);
  const [createModalVisible, setCreateVisible] = useState<boolean>(false);
  const [currentRow, setCurrentRow] = useState<Serverless.Tunnel>();

  return <PageContainer
    title={"函数实例"}
    tabList={[
      {
        tab: '隧道实例',
        key: 'tunnel',
      },
    ]}
    extra={""}>
    <ProList<Serverless.Tunnel>
      className={styles.proList}
      actionRef={actionRef}
      pagination={{
        defaultPageSize: 12,
        showSizeChanger: true,
      }}
      request={async (params, sort, filter) => {
        return getServerlessTunnel(params.current === undefined ? 0 : params.current - 1, params.pageSize === undefined ? 10 : params.pageSize);
      }}
      rowKey={"ID"}
      showActions="hover"
      rowSelection={{}}
      grid={{gutter: 16, xs: 1, sm: 2, md: 2, lg: 2, xl: 3, xxl: 3}}
      onItem={(record: Serverless.Tunnel) => {
        return {
          onClick: () => {
            setCurrentRow(record);
            setDetailVisible(true);
          },
        };
      }}
      toolBarRender={() => {
        return [
          // <SSrDropDown/>,
          <ClashDropDown/>,
          <ShadowRocketDropDown/>,
          <Button key="button" icon={<PlusOutlined/>} type="primary"
                  style={{marginLeft: "10px"}}
                  onClick={() => {
                    setCreateVisible(true)
                  }}>
            新增
          </Button>
        ]
      }}
      metas={{
        title: {
          dataIndex: 'name',
        },
        subTitle: {},
        avatar: {
          dataIndex: 'type',
          formItemProps: {
            style: {
              fontSize: "150%",
            }
          },
          render: (_, record) => {
            return TunnelTypeValueEnum[record.type]
          }
        },
        content: {
          render: (dom, record) => {
            return <ProCard gutter={8} bordered={false} split={"horizontal"}>
              <ProCard bordered={false} split={"horizontal"}>
                <Statistic title="当前状态:" valueRender={() => <Badge style={{fontSize: "12px"}}
                                                                   status={TunnelStatusEnum[record.status]?.status}
                                                                   text={TunnelStatusEnum[record.status]?.text}/>}/>
                <Statistic title="账户类型:" valueRender={() => {
                  return <div>{CloudProvideTypeValueEnum[record.provider_type]} - {RegionEnum[record.tunnel_config.region]}</div>
                }}/>
                <Statistic title="隧道地址:"
                           value={record.address === undefined || record.address === null ? "-" : record.address}
                           valueRender={() => {
                             return record.address.length > 40 ? <Tooltip
                               title={record.address}>{record.address.substring(0, 39) + "..."}</Tooltip> : record.address
                           }}
                />
                <Statistic title="远程版本:"
                           valueRender={() => {
                             if (record.version) {
                               // 版本格式: 2.1.0-dev2-commitHash(40字符)，从后截取commit，剩下的是版本号
                               // 兼容低版本格式如 2.0.0- 没有完整commit的情况
                               const fullVersion = record.version;
                               let versionNum = fullVersion;
                               if (fullVersion.length > 41) {
                                 versionNum = fullVersion.slice(0, -41);
                               } else if (fullVersion.endsWith('-')) {
                                 versionNum = fullVersion.slice(0, -1);
                               }
                               // v2ray 版本格式: v2ray-core:-5.16.1，取 :- 后面的部分，兼容没有v2ray版本的情况
                               const v2rayVer = (record.v2ray_version?.trim() && record.v2ray_version.split(':-')[1]) || record.v2ray_version?.trim() || '-';
                               return <Tooltip title={`v2ray: ${v2rayVer}`}>
                                 <span style={{color: '#52c41a'}}>{versionNum}</span>
                               </Tooltip>
                             }
                             return <span style={{color: '#999'}}>-</span>
                           }}
                />
                <Statistic title="函数规格:"
                           valueRender={() => <Space><IconFont
                             type={"icon-cpu1"}/>{record.tunnel_config.cpu} M <IconFont
                             type={"icon-memory1"}/>{record.tunnel_config.memory} Mi </Space>}/>
              </ProCard>
            </ProCard>
          }
        },
        actions: {
          cardActionProps: 'extra',
          render: (_, record) => {
            return <div>
              {record.tunnel_config.tls ? <Tag color={"magenta"}>tls</Tag> : <></>}
              {record.tunnel_config.tor ? <Tag color={"blue"}>tor</Tag> : <></>}
            </div>
          }
        },
      }}
      headerTitle={<Space><IconFont type={"icon-tunnel_statistics_icon_tunnel"} style={{fontSize: "150%"}}/>隧道 -
        Tunnel</Space>}
    />
    <CreateForm
      onSubmit={async (value) => {
        const success = await handleCreateTunnel(value);
        setCreateVisible(false);
        if (actionRef.current) {
          actionRef.current.reload();
        }
        if (success) {
          setCurrentRow(undefined);
        }
      }}
      onCancel={() => {
        setCreateVisible(false);
        setCurrentRow(undefined);
      }}
      createModalVisible={createModalVisible}
      values={currentRow || {}}
    />
    <DetailDrawer
      onSubmit={async (value) => {
        const success = await handleUpdateTunnel(value)
        setDetailVisible(false);
        if (actionRef.current) {
          actionRef.current.reload();
        }
        if (success) {
          setCurrentRow(undefined);
        }
      }}
      onDelete={async (value) => {
        // 检查当前状态是否为停止，如果非停止，则禁止删除。
        const success = await handleDeleteTunnel(value.id)
        setDetailVisible(false);
        if (success) {
          setCurrentRow(undefined);
        }
        if (actionRef.current) {
          actionRef.current.reload();
        }
      }}
      onCancel={() => {
        setDetailVisible(false);
        setCurrentRow(undefined);
      }}
      detailVisible={detailModalVisible}
      values={currentRow || {}}
    />
  </PageContainer>
}

export default Tunnel
