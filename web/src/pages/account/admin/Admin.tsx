import React, {useRef} from "react";
import type {ActionType, ProColumns} from "@ant-design/pro-components";
import {ProTable} from "@ant-design/pro-components";
import UpdateForm from "@/pages/account/admin/UpdateForm";
import {message} from "antd";

const AdminTable: React.FC = () => {
  const actionRef = useRef<ActionType>();

  const columns: ProColumns<Account.Admin>[] = [
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
    },
    {
      title: '最后一次登录时间',
      key: 'last_time',
      dataIndex: 'last_time',
      valueType: 'dateTime',
      responsive: ['md'],
      hideInSearch: true,
      disable: true,
    },
    {
      title: '最后一次登录地址',
      key: 'last_addr',
      dataIndex: 'last_addr',
    },
    {
      title: '其他操作',
      key: 'option',
      valueType: 'option',
      render: () => [
        <a key="forbidden" onClick={() => {
        }}>
          禁用
        </a>,
        <UpdateForm key={"update"} onSubmit={async (values) => {
          // check validator first
          if (values.password === values.old_password) {
            message.error("新密码与旧密码相同");
            return true;
          }
          if (values.password !== values.password_repeat) {
            message.error("两次输入的密码不一致");
            return true;
          }
          // do update
          return true;
        }}/>
      ],
    },
  ]

  return <>
    <ProTable<Account.Admin>
      columns={columns}
      actionRef={actionRef}
      cardBordered
      rowKey="id"
      headerTitle={"客户端管理账户"}
      search={false}
      options={false}
      pagination={{
        defaultPageSize: 12,
        showSizeChanger: true,
      }}
      dateFormatter="string"
      dataSource={[{
        name: "admin",
        type: "admin",
        last_addr: "127.0.0.1",
        last_time: "2024-08-06",
      }]}
    />
  </>
}

export default AdminTable
