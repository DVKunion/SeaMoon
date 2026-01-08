import React, {useRef, useState} from 'react';
import type {ProDescriptionsActionType} from '@ant-design/pro-components';
import {ProDescriptions} from '@ant-design/pro-components';
import {Button, Divider, Drawer, Modal, Popconfirm, Space, Tooltip} from "antd";
import {TunnelAuthFCTypeEnum, TunnelStatusTag, TunnelTypeValueEnum} from "@/enum/tunnel";
import {CloudProvideTypeValueEnum, RegionEnum} from "@/enum/cloud";
import {CheckCircleTwoTone, CloseCircleTwoTone, ExclamationCircleOutlined, PoweroffOutlined, SyncOutlined} from "@ant-design/icons";
// @ts-ignore
import {CopyToClipboard} from 'react-copy-to-clipboard';
import {FormValueType} from "./CreateForm";
import {getTunnelDependents} from "@/services/function/api";

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

  // 执行删除操作，带依赖检查提示
  const handleDelete = async () => {
    if (!props.values.id) {
      props.onDelete(props.values as FormValueType);
      return;
    }

    try {
      const result = await getTunnelDependents(props.values.id);
      if (result.success && result.data && result.data.length > 0) {
        // 有依赖的隧道，显示确认对话框
        Modal.confirm({
          title: '删除确认',
          icon: <ExclamationCircleOutlined />,
          content: (
            <div>
              <p>以下函数依赖该隧道作为级联代理，删除后这些函数也将被一并删除：</p>
              <ul>
                {result.data.map((dep: Serverless.Tunnel) => (
                  <li key={dep.id}>{dep.name}</li>
                ))}
              </ul>
              <p>确定要继续删除吗？</p>
            </div>
          ),
          okText: '确认删除',
          okType: 'danger',
          cancelText: '取消',
          onOk: () => {
            props.onDelete(props.values as FormValueType);
          },
        });
      } else {
        // 没有依赖，直接删除
        props.onDelete(props.values as FormValueType);
      }
    } catch (error) {
      // 查询失败，直接删除
      props.onDelete(props.values as FormValueType);
    }
  };

  return <Drawer
    title="函数详情"
    width={window.innerWidth < 768 ? "80%" : "39%"}
    onClose={props.onCancel}
    open={props.detailVisible}
    extra={
      <Space>
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
          onConfirm={handleDelete}
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
          title: '远程版本',
          key: 'version',
          editable: false,
          dataIndex: 'version',
          render: (dom, entry) => {
            if (entry.version) {
              // 版本格式: 2.1.0-dev2-commitHash(40字符)，从后截取commit，剩下的是版本号
              // 兼容低版本格式如 2.0.0- 没有完整commit的情况
              const fullVersion = entry.version;
              let versionNum = fullVersion;
              if (fullVersion.length > 41) {
                versionNum = fullVersion.slice(0, -41);
              } else if (fullVersion.endsWith('-')) {
                versionNum = fullVersion.slice(0, -1);
              }
              return <span style={{color: '#52c41a'}}>{versionNum}</span>
            }
            return <span style={{color: '#999'}}>-</span>
          }
        },
        {
          title: 'Commit',
          key: 'commit',
          editable: false,
          dataIndex: 'version',
          render: (dom, entry) => {
            // commit hash 是最后40个字符，需要确保有完整的 commit
            if (entry.version && entry.version.length > 41) {
              const commit = entry.version.slice(-40);
              // 验证是否是有效的 commit hash (只包含十六进制字符)
              if (/^[0-9a-f]{40}$/i.test(commit)) {
                return <span style={{fontSize: '12px', color: '#666'}}>{commit}</span>
              }
            }
            return <span style={{color: '#999'}}>-</span>
          }
        },
        {
          title: 'V2Ray版本',
          key: 'v2ray_version',
          editable: false,
          dataIndex: 'v2ray_version',
          render: (dom, entry) => {
            if (entry.v2ray_version && entry.v2ray_version.trim()) {
              // v2ray 版本格式: v2ray-core:-5.16.1，取 :- 后面的部分
              const v2rayVer = entry.v2ray_version.split(':-')[1] || entry.v2ray_version;
              return <span>{v2rayVer}</span>
            }
            return <span style={{color: '#999'}}>-</span>
          }
        },
        {
          title: '最后检查时间',
          key: 'last_check_time',
          editable: false,
          dataIndex: 'last_check_time',
          span: 2,
          render: (dom, entry) => {
            if (entry.last_check_time) {
              return <span>{entry.last_check_time}</span>
            }
            return <span style={{color: '#999'}}>-</span>
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
          key: 'tor',
          editable: false,
          dataIndex: "tunnel_config.tor",
          render: (dom, record) => {
            return record.tunnel_config?.tor ? <CheckCircleTwoTone style={{marginTop: "4px"}}  twoToneColor="#52c41a" /> :
              <CloseCircleTwoTone style={{marginTop: "4px"}} twoToneColor="#eb2f96"/>
          }
        },
        {
          title: '是否开启级联代理',
          key: 'cascade_proxy',
          editable: false,
          dataIndex: "tunnel_config.cascade_proxy",
          render: (dom, record) => {
            return record.tunnel_config?.cascade_proxy ? <CheckCircleTwoTone style={{marginTop: "4px"}}  twoToneColor="#52c41a" /> :
              <CloseCircleTwoTone style={{marginTop: "4px"}} twoToneColor="#eb2f96"/>
          }
        },
        {
          title: '级联代理地址',
          key: 'cascade_addr',
          editable: false,
          span: 2,
          dataIndex: "tunnel_config.cascade_addr",
          render: (dom, record) => {
            if (record.tunnel_config?.cascade_proxy && record.tunnel_config?.cascade_addr) {
              return <span>{record.tunnel_config.cascade_addr}</span>
            }
            return <span style={{color: '#999'}}>-</span>
          }
        }
      ]}
      dataSource={props.values}
    />
    <Divider/>
    {/*<ProDescriptions*/}
    {/*  title={"订阅信息"}*/}
    {/*  column={1}*/}
    {/*  actionRef={actionRef}*/}
    {/*  editable={{*/}
    {/*    onSave: async (keypath, newInfo, oriInfo) => {*/}
    {/*      props.values[keypath.toString()] = newInfo[keypath.toString()];*/}
    {/*      return true;*/}
    {/*    },*/}
    {/*  }}*/}
    {/*  columns={[*/}
    {/*    {*/}
    {/*      key: 'text',*/}
    {/*      editable: false,*/}
    {/*      render: (dom, entity) => {*/}
    {/*        return <a> 点击下载 clash 订阅配置 </a>*/}
    {/*      },*/}
    {/*    }*/}
    {/*  ]}*/}
    {/*  dataSource={props.values}*/}
    {/*/>*/}
  </Drawer>
}

export default DetailDrawer
