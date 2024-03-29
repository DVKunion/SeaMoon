import React, {useRef} from 'react';
import type {ProDescriptionsActionType} from '@ant-design/pro-components';
import {ProDescriptions, ProFormSelect} from '@ant-design/pro-components';
import {FormValueType} from "@/pages/service/components/CreateForm";
import {Button, Divider, Drawer, Popconfirm, Space} from "antd";
import {PoweroffOutlined, SyncOutlined, ThunderboltTwoTone} from "@ant-design/icons";
import {ProxyDynamicTagList, ProxyTypeValueEnum} from "@/enum/service";
import {SpeedTransfer} from "@/components/SpeedTransfer";


export type DetailProps = {
  spin: boolean;
  onCancel: () => void;
  onDelete: (values: FormValueType) => Promise<void>;
  onSubmit: (values: FormValueType) => Promise<void>;
  onSpeed: (values: FormValueType) => Promise<void>;
  detailVisible: boolean;
  values: Partial<Service.Proxy>;
};

const DetailDrawer: React.FC<DetailProps> = (props) => {
  const actionRef = useRef<ProDescriptionsActionType>();

  return <Drawer
    title="服务详情"
    width={"39%"}
    onClose={props.onCancel}
    open={props.detailVisible}
    extra={
      <Space>
        <Button shape="round" icon={<ThunderboltTwoTone/>}
                onClick={() => {
                  props.onSpeed(props.values as FormValueType)
                }}
        > 测速</Button>
        {props.values.status === 3 ? <Button
          type={"primary"}
          shape={"round"}
          icon={<SyncOutlined spin={false}/>}
          onClick={() => {
            props.values.status = 2;
            props.onSubmit(props.values as FormValueType);
          }}
        >启动</Button> : <Button
          type={"primary"}
          shape={"round"}
          danger
          onClick={() => {
            props.values.status = 3;
            props.onSubmit(props.values as FormValueType);
          }}
          icon={<PoweroffOutlined/>}>停止</Button>}
      </Space>
    }
    footer={
      <Space style={{float: "right"}}><Button type={"primary"} onClick={() => {
        props.onSubmit(props.values as FormValueType);
      }}>更新</Button>
        <Popconfirm
          title="删除服务"
          onConfirm={() => {
            props.onDelete(props.values as FormValueType);
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
          title: '当前状态',
          key: 'status',
          editable: false,
          dataIndex: 'status',
          span: 2,
          render: (dom, entry) => {
            return <ProxyDynamicTagList status={entry.status ? entry.status : 0} spin={props.spin} msg={entry.status_message}/>
          }
        },
        {
          title: '代理名称',
          key: 'text',
          editable: false,
          dataIndex: 'name',
        },
        {
          title: '代理类型',
          dataIndex: 'type',
          render: (dom: any, entity: any, index: any, action: any) => {
            return ProxyTypeValueEnum[entity.type]
          },
          renderFormItem: () => {
            return <ProFormSelect
              name="type"
              width="md"
              placeholder={""}
              valueEnum={ProxyTypeValueEnum}
            />
          },
        },
        {
          title: '监听地址',
          key: 'text',
          dataIndex: 'listen_address',
        },
        {
          title: '监听端口',
          key: 'text',
          dataIndex: 'listen_port',
        },
        {
          title: '创建时间',
          key: 'date',
          valueType: "dateTime",
          editable: false,
          dataIndex: 'created_at',
        },
        {
          title: '修改时间',
          key: 'date',
          valueType: "dateTime",
          editable: false,
          dataIndex: 'updated_at',
        }
      ]}
      dataSource={props.values}
    />
    <Divider/>
    <ProDescriptions
      title={"关联信息 (todo)"}
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
        },
        {
          title: '关联函数实例信息',
          key: 'text',
          editable: false,
          render: (dom, entity) => {
            return <a> 函数实例详情 </a>
          },
        }
      ]}
      dataSource={props.values}
    />
    <Divider/>
    <ProDescriptions
      title={"网络数据"}
      column={2}
      actionRef={actionRef}
      columns={[
        {
          title: '当前连接数',
          key: 'number',
          editable: false,
          dataIndex: 'conn',
        },
        {
          title: '当前延迟',
          key: 'number',
          editable: false,
          dataIndex: 'lag',
          render: (dom, entity) => {
            return entity.lag + " ms"
          }
        },
        {
          title: '上行速率',
          key: 'number',
          editable: false,
          dataIndex: 'speed_up',
          render: (dom, entity) => {
            return entity.speed_up === 0 ? "未测速" : entity.speed_up + " Mbps"
          }
        },
        {
          title: '下行速率',
          key: 'number',
          editable: false,
          dataIndex: 'speed_down',
          render: (dom, entity) => {
            return entity.speed_down === 0 ? "未测速" : entity.speed_down + " Mbps"
          }
        },
        {
          title: '总流入',
          key: 'in_bound',
          editable: false,
          dataIndex: 'in_bound',
          render: (dom, entity) => {
            return SpeedTransfer({
              bytes: entity.in_bound || 0
            })
          }
        },
        {
          title: '总流出',
          key: 'out_bound',
          editable: false,
          dataIndex: 'out_bound',
          render: (dom, entity) => {
            return SpeedTransfer({
              bytes: entity.out_bound || 0
            })
          }
        },
      ]}
      dataSource={props.values}
    >
      <ProDescriptions.Item
        label="24小时流量数据"
        span={2}
      >
        todo
      </ProDescriptions.Item>
      <ProDescriptions.Item
        label="其他图形"
        span={2}
      >
        todo
      </ProDescriptions.Item>
    </ProDescriptions>
    <Divider/>
  </Drawer>
}

export default DetailDrawer
