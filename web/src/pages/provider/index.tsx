import React, {useRef, useState} from "react";
import {PageContainer, ProTable, ActionType, ProColumns} from "@ant-design/pro-components";
import {Button, Popconfirm} from "antd";
import {PlusOutlined} from "@ant-design/icons";
import {getCloudProvider} from "@/services/cloud/api";
import {handleCreateCloud, handleDeleteCloud, handleSyncCloud, handleUpdateCloud} from './handle';
import {CloudProvideTypeEnum, CloudProviderStatusEnum} from "@/enum/cloud";
import CreateForm from './components/CreateForm';
import DetailDrawer from "./components/DetailDrawer";
import {Badge, Tooltip} from "_antd@4.24.15@antd";

const Provider: React.FC = () => {

  const actionRef = useRef<ActionType>();
  const [detailModalVisible, setDetailVisible] = useState<boolean>(false);
  const [createModalVisible, setCreateVisible] = useState<boolean>(false);
  const [currentRow, setCurrentRow] = useState<Cloud.Provider>();


  const columns: ProColumns<Cloud.Provider>[] = [
    {
      dataIndex: 'id',
      key: 'id',
      valueType: 'indexBorder',
      width: 48,
    },
    {
      title: '账户名称',
      key: 'name',
      dataIndex: 'name',
      copyable: true,
      ellipsis: true,
    },
    {
      title: '账户类型',
      key: 'type',
      dataIndex: 'type',
      filters: true,
      onFilter: true,
      ellipsis: true,
      valueType: 'select',
      valueEnum: CloudProvideTypeEnum,
    },
    {
      title: '已部署',
      key: "count",
      dataIndex: 'count',
      sorter: true,
    },
    {
      title: '账户余额',
      key: "amount",
      dataIndex: 'amount',
      valueType: 'money',
      render: (_, record) => {
        return "¥ " + record.info.amount
      }
    },
    {
      title: '账户状态',
      key: 'status',
      dataIndex: 'status',
      filters: true,
      onFilter: true,
      ellipsis: true,
      valueType: 'select',
      render: (dom, entry) => {
        return <Tooltip title={entry.status_message}>
          <Badge style={{fontSize: "12px"}}
                 status={CloudProviderStatusEnum[entry.status]?.status}
                 text={CloudProviderStatusEnum[entry.status]?.text}
          /></Tooltip>
      },
      valueEnum: CloudProviderStatusEnum,
    },
    {
      title: '创建时间',
      key: 'showTime',
      dataIndex: 'created_at',
      valueType: 'dateTime',
      sorter: true,
      hideInSearch: true,
      disable: true,
    },
    {
      title: '操作',
      key: 'option',
      valueType: 'option',
      render: (text, record, _, action) => [
        <a key="detail" onClick={() => {
          setCurrentRow(record);
          setDetailVisible(true);
        }}>
          详情
        </a>,
        <Popconfirm
          title={"删除账户将会删除对应的函数实例、同时关联的服务也会被停止。确认删除?"}
          onConfirm={async () => {
            await handleDeleteCloud(record.id);
            if (actionRef.current) {
              actionRef.current.reload();
            }
          }}
          okText="确认"
          cancelText="取消"
        >
          <a key="delete"> 删除</a>
        </Popconfirm>
      ],
    },
  ];

  return <PageContainer
    title={"云账户"}>
    <ProTable<Cloud.Provider>
      columns={columns}
      actionRef={actionRef}
      cardBordered
      // @ts-ignore
      request={async (params, sort, filter) => {
        return getCloudProvider(params.current === undefined ? 0 : params.current - 1, params.pageSize === undefined ? 10 : params.pageSize)
      }}
      rowKey="id"
      headerTitle={"账户信息"}
      search={false}
      options={false}
      form={{
        // 由于配置了 transform，提交的参与与定义的不同这里需要转化一下
        syncToUrl: (values, type) => {
          if (type === 'get') {
            return {
              ...values,
              created_at: [values.startTime, values.endTime],
            };
          }
          return values;
        },
      }}
      pagination={{
        defaultPageSize: 12,
        showSizeChanger: true,
      }}
      dateFormatter="string"
      toolBarRender={() => [
        <Button key="button" icon={<PlusOutlined/>} type="primary" onClick={() => {
          setCreateVisible(true)
        }}>
          新增
        </Button>,
      ]}
    />
    <CreateForm
      onSubmit={async (value) => {
        const success = await handleCreateCloud(value);
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
        const success = await handleUpdateCloud(value)
        setDetailVisible(false);
        if (actionRef.current) {
          actionRef.current.reload();
        }
        if (success) {
          setCurrentRow(undefined);
        }
      }}
      onSync={async (value) => {
        const success = await handleSyncCloud(value.id)
        setDetailVisible(false);
        if (actionRef.current) {
          actionRef.current.reload();
        }
        if (success) {
          setCurrentRow(undefined);
        }
      }}
      onDelete={async (value) => {
        await handleDeleteCloud(value.id)
        setDetailVisible(false);
        setCurrentRow(undefined);
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

export default Provider
