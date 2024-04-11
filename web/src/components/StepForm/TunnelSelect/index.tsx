import React, {useState} from "react";
import {ProFormSelect} from "@ant-design/pro-components";
import {Space} from "antd";
import {getServerlessTunnel} from "@/services/function/api";
import {TunnelTypeValueEnum} from "@/enum/tunnel";
import {CloudProvideTypeValueEnum} from "@/enum/cloud";


export type TunnelProps = {
  tor: boolean,
  values: Partial<Serverless.Tunnel>
};

export const TunnelSelect: React.FC<TunnelProps> = (props: TunnelProps) => {
  const [tunnel, setTunnel] = useState<Partial<Serverless.Tunnel>>({});

  return <><ProFormSelect
    dependencies={[props.tor]}
    name="tunnel_id"
    label="选择关联函数实例"
    width={"xl"}
    tooltip={"仅允许正常状态的实例，如果开启了 tor, 则会自动筛选 tor 标签的实例"}
    showSearch={true}
    placeholder={""}
    request={async () => {
      const res: { key: number; label: JSX.Element; value: number; obj: Serverless.Tunnel; }[] = [];
      const {data} = await getServerlessTunnel(0, 999999);
      data.forEach((item) => {
        if(!props.tor || (props.tor && item.tunnel_config.tor)) {
          res.push({
            key: item.id,
            label: <Space>{CloudProvideTypeValueEnum[item.provider_type || 0]} - {TunnelTypeValueEnum[item.type]} - {item.name}</Space>,
            value: item.id,
            obj: item
          });
        }
      })
      return res
    }}
    rules={[
      {
        required: true,
        message: "请选择关联函数实例!",
      },
    ]}
    fieldProps={
      {
        onSelect: (value, option) => {
          setTunnel(option["data-item"].obj);
        }
      }
    }
  />
    {tunnel !== undefined && tunnel.id !== 0 ?
      <Space>
        {/*{cloud.amount !== undefined ? <>账户余额: <Tag*/}
        {/*  color={cloud.amount > 0 ? "volcano" : "green"}>{"¥" + cloud.amount}</Tag></> : <></>}*/}
        {/*{cloud.count !== undefined && cloud.max_limit !== undefined ?*/}
        {/*  <>已部署函数限制: <Tag*/}
        {/*    color={cloud.max_limit === 0 ? "volcano" : cloud.count <= cloud.max_limit ? "volcano" : "green"}>{cloud.count + " / " + (cloud.max_limit === 0 ? "∞" : cloud.max_limit)}</Tag></> : <></>*/}
        {/*}*/}
      </Space> : <></>
    }
  </>
}
