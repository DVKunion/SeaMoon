import React, {useRef, useState} from 'react';
import type {ProDescriptionsActionType} from '@ant-design/pro-components';
import {ProDescriptions} from '@ant-design/pro-components';
import {Button, Divider, Drawer, message, Popconfirm, Space, Tooltip} from "antd";
import {TunnelAuthFCTypeEnum, TunnelStatusTag, TunnelTypeValueEnum} from "@/enum/tunnel";
import {CloudProvideTypeValueEnum, RegionEnum} from "@/enum/cloud";
import {CheckCircleTwoTone, CloseCircleTwoTone, CopyOutlined, PoweroffOutlined, SyncOutlined} from "@ant-design/icons";
// @ts-ignore
import {CopyToClipboard} from 'react-copy-to-clipboard';
import {FormValueType} from "./CreateForm";

export type DetailProps = {
  onCancel: () => void;
  onDelete: (values: FormValueType) => Promise<void>;
  onSubmit: (values: FormValueType) => Promise<void>;
  detailVisible: boolean;
  values: Partial<Serverless.Tunnel>;
};

const DetailDrawer: React.FC<DetailProps> = (props) => {
  const actionRef = useRef<ProDescriptionsActionType>();
  const [spin, setSpin] = useState<boolean>(false);

  return <Drawer
    title="函数详情"
    width={"39%"}
    onClose={props.onCancel}
    open={props.detailVisible}
    extra={
      <Space>
        <CopyToClipboard
          text={props.values.address}
          onCopy={() => message.success("已复制函数地址")}>
          <Button shape={"round"} ghost icon={<CopyOutlined/>}></Button>
        </CopyToClipboard>
        {props.values.status === 3 ? <Button
          onMouseEnter={() => {
            setSpin(true);
          }}
          onMouseLeave={() => {
            setSpin(false);
          }}
          type={"primary"}
          shape={"round"}
          icon={<SyncOutlined spin={spin}/>}
          onClick={() => {
            props.values.status = 1;
            props.onSubmit(props.values as FormValueType);
          }}
        >启用</Button> : <Popconfirm
          title="停用仅会删除远端实例，不会删除此条记录，您可以再次通过启动来部署。"
          onConfirm={() => {
            props.values.status = 3;
            props.onSubmit(props.values as FormValueType);
          }}
          okText="确认"
          cancelText="取消"
        ><Button
          type={"primary"}
          shape={"round"}
          danger
          icon={<PoweroffOutlined/>}>停用</Button></Popconfirm>}
      </Space>
    }
    footer={
      <Space style={{float: "right"}}>
      {/*  <Button type={"primary"} onClick={() => {*/}
      {/*  props.onSubmit(props.values).then();*/}
      {/*}}>更新</Button>*/}
        <Popconfirm
          title="删除函数?"
          onConfirm={() => {
            props.onDelete(props.values as FormValueType).then();
          }}
          okText="确认"
          cancelText="取消"
        >
          <Button type={"primary"} danger>删除</Button>
        </Popconfirm>
      </Space>
    }
  >
    <ProDescriptions
      title={"基本信息"}
      column={2}
      actionRef={actionRef}
      editable={{
        onSave: async (keypath, newInfo, oriInfo) => {
          props.values[keypath.toString()] = newInfo[keypath.toString()];
          return true;
        },
      }}
      columns={[
        {
          title: 'UniqID',
          key: 'uniq_id',
          editable: false,
          dataIndex: 'uniq_id',
          span: 2
        },
        {
          title: '函数当前状态',
          key: 'status',
          editable: false,
          dataIndex: 'status',
          span: 2,
          render: (dom, entry) => {
            return <Tooltip title={entry.status_message}> {TunnelStatusTag[entry.status ? entry.status : 0]}</Tooltip>
          }
        },
        {
          title: '账户类型',
          dataIndex: 'provider_type',
          key: 'type',
          editable: false,
          valueEnum: CloudProvideTypeValueEnum,
        },
        {
          title: '隧道类型',
          dataIndex: 'type',
          key: 'type',
          editable: false,
          valueEnum: TunnelTypeValueEnum,
        },
        {
          title: '隧道名称',
          span: 2,
          key: 'name',
          copyable: true,
          editable: false,
          dataIndex: 'name',
        },
        {
          title: '隧道地址',
          copyable: true,
          editable: false,
          span: 2,
          key: 'address',
          dataIndex: 'address',
        },
        {
          title: '端口号',
          key: 'port',
          editable: false,
          dataIndex: 'port',
        },
        {
          title: '所在区域',
          key: 'region',
          editable: false,
          dataIndex: 'cloud_provider_region',
          render: (dom, record) => {
            return RegionEnum[record.tunnel_config?.region || ""]
          }
        },
        {
          title: '创建时间',
          key: 'created_at',
          valueType: "dateTime",
          editable: false,
          dataIndex: 'created_at',
        },
        {
          title: '修改时间',
          key: 'updated_at',
          valueType: "dateTime",
          editable: false,
          dataIndex: 'updated_at',
        }
      ]}
      dataSource={props.values}
    />
    <Divider/>
    <ProDescriptions
      title={"配置信息"}
      column={2}
      actionRef={actionRef}
      editable={{
        onSave: async (keypath, newInfo, oriInfo) => {
          props.values[keypath.toString()] = newInfo[keypath.toString()];
          return true;
        },
      }}
      columns={[
        {
          title: 'CPU规格',
          key: 'tunnel_config.cpu',
          editable: false,
          dataIndex: "tunnel_config",
          render: (dom, record) => {
            return record.tunnel_config?.cpu + " M"
          }
        },
        {
          title: '内存规格',
          key: 'tunnel_config.memory',
          editable: false,
          dataIndex: "tunnel_config",
          render: (dom, record) => {
            return record.tunnel_config?.memory  + " Mi"
          }
        },
        {
          title: '最大实例并发数',
          key: 'tunnel_config.instance',
          editable: false,
          dataIndex: "tunnel_config",
          render: (dom, record) => {
            return record.tunnel_config?.instance
          }
        },
        {
          title: '函数认证方式',
          key: 'tunnel_config.auth_type',
          editable: false,
          dataIndex: "tunnel_config.tunnel_auth_type",
          render: (dom, record) => {
            return TunnelAuthFCTypeEnum[record.tunnel_config?.tunnel_auth_type || 1]
          }
        },
        {
          title: '是否开启 tls 认证',
          key: 'tls',
          editable: false,
          dataIndex: "tunnel_config.tls",
          render: (dom, record) => {
            return record.tunnel_config?.tls ? <CheckCircleTwoTone style={{marginTop: "4px"}}  twoToneColor="#52c41a" /> :
              <CloseCircleTwoTone style={{marginTop: "4px"}} twoToneColor="#eb2f96"/>
          }
        },
        {
          title: '是否开启 tor 网桥',
          key: 'tls',
          editable: false,
          dataIndex: "tunnel_config.tor",
          render: (dom, record) => {
            return record.tunnel_config?.tor ? <CheckCircleTwoTone style={{marginTop: "4px"}}  twoToneColor="#52c41a" /> :
              <CloseCircleTwoTone style={{marginTop: "4px"}} twoToneColor="#eb2f96"/>
          }
        }
      ]}
      dataSource={props.values}
    />
    <Divider/>
    <ProDescriptions
      title={"关联信息 (todo) "}
      column={1}
      actionRef={actionRef}
      editable={{
        onSave: async (keypath, newInfo, oriInfo) => {
          props.values[keypath.toString()] = newInfo[keypath.toString()];
          return true;
        },
      }}
      columns={[
        {
          title: '关联云账户信息',
          key: 'text',
          editable: false,
          render: (dom, entity) => {
            return <a> 云账户详情 </a>
          },
        }
      ]}
      dataSource={props.values}
    />
  </Drawer>
}

export default DetailDrawer
