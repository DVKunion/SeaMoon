import React, {useState} from "react";
import {ProFormSelect} from "@ant-design/pro-components";
import {getActiveProvider} from "@/services/cloud/api";
import {Space, Tag } from "antd";
import {CloudProvideTypeIcon, RegionEnum} from "@/enum/cloud";

export type ProviderProps = {
  onChange: (values: number) => void;
};

export const ProviderSelect: React.FC<ProviderProps> = (props: ProviderProps) => {
  const [cloud, setCloud] = useState<Partial<Cloud.Provider>>({});

  return <><ProFormSelect
      name="provider_id"
      label="选择关联云账户"
      width={"xl"}
      tooltip={"仅允许正常状态的账户"}
      showSearch={true}
      request={async () => {
        const res: { key: number; label: JSX.Element; value: number; obj: Cloud.Provider; }[] = [];
        const {data} = await getActiveProvider();
        data.forEach((item) => {
          res.push(
            {
              key: item.id,
              label: <Space>{CloudProvideTypeIcon[item.type]}{item.name}</Space>,
              value: item.id,
              obj: item
            }
          )
        })
        return res
      }}
      rules={[
        {
          required: true,
          message: "请选择关联云账户!",
        },
      ]}
      fieldProps={
        {
          onSelect: (value, option) => {
            setCloud(option["data-item"].obj);
            props.onChange(option["data-item"].obj.type)
          }
        }
      }
    />
    {cloud.id !== undefined && cloud.id !== 0 ?
      <>
        <Space size={120}>
        {cloud.info?.amount !== undefined ? <div><p>账户余额: </p><Tag
          color={cloud.info.amount > 0 ? "volcano" : "green"}>{"¥" + cloud.info?.amount}</Tag></div> : <></>}
        {cloud.count !== undefined && cloud.max_limit !== undefined ?
          <div><p>已部署函数限制: </p><Tag
            color={cloud.max_limit === 0 ? "volcano" : cloud.count <= cloud.max_limit ? "volcano" : "green"}>{cloud.count + " / " + (cloud.max_limit === 0 ? "∞" : cloud.max_limit)}</Tag></div> : <></>
        }
        </Space>
        <p style={{marginTop: "20px"}}>允许部署区域:</p>
      <Space>
        {cloud.regions?.map((region, index) => (
          <Tag key={index}>{RegionEnum[region]}</Tag>
        ))}
      </Space></> : <></>
    }
    </>
}
