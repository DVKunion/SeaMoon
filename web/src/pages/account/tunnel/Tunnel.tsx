import React, {useRef} from "react";
import {type ActionType, type ProColumns, ProTable} from "@ant-design/pro-components";

const TunnelTable: React.FC = () => {
  const actionRef = useRef<ActionType>();
  const columns: ProColumns<Account.Tunnel>[] = [
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
      copyable: true,
      ellipsis: true,
    },
  ]

  return <>
    <ProTable<Account.Tunnel>
      columns={columns}
      actionRef={actionRef}
      cardBordered
      rowKey="id"
      headerTitle={"隧道账户"}
      search={false}
      options={false}
      pagination={{
        defaultPageSize: 12,
        showSizeChanger: true,
      }}
      dateFormatter="string"
    /></>
}

export default TunnelTable
