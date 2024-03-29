import React from 'react';
import {ProFormText, ProFormTextArea, ProFormSelect, ProForm} from "@ant-design/pro-components";
import {ALiYunRegionEnum, TencentRegionEnum, SealosRegionEnum} from "@/enum/cloud";

export type CloudProviderType = {
  type: number
}

export const CloudProviderAuthForm: React.FC<CloudProviderType> = (props) => {
  switch (props.type) {
    case 1: // 阿里云 AID AK SK
      return <ProForm.Group>
        <ProFormSelect
          name="regions"
          width="xl"
          label="账户允许部署地区"
          placeholder={""}
          tooltip={"选择云账户允许部署的地区"}
          valueEnum={ALiYunRegionEnum}
          fieldProps={{
            mode: "multiple",
          }}
          rules={[
            {
              required: true,
              message: "请选择云账户所属的地区!",
            },
          ]}
        />
        <ProFormText
          name={"access_id"}
          label={"AccessId"}
          width="xl"
          placeholder={""}
          rules={[
            {
              required: true,
              message: "请填写主账户ID信息!",
            },
          ]}
        />
        <ProFormText
          name={"access_key"}
          label={"AccessKey"}
          width="xl"
          placeholder={""}
          rules={[
            {
              required: true,
              message: "请填写认证信息!",
            },
          ]}
        />
        <ProFormText
          name={"access_secret"}
          label={"SecretKey"}
          width="xl"
          placeholder={""}
          rules={[
            {
              required: true,
              message: "请填写认证信息!",
            },
          ]}
        />
      </ProForm.Group>;
    case 2:
      return <ProForm.Group>
        <ProFormSelect
          name="regions"
          width="xl"
          label="账户允许部署地区"
          showSearch={true}
          placeholder={""}
          tooltip={"选择云账户允许部署的地区"}
          valueEnum={TencentRegionEnum}
          fieldProps={{
            mode: "multiple",
          }}
          rules={[
            {
              required: true,
              message: "请选择云账户所属的地区!",
            },
          ]}
        />
        <ProFormText
          name={"access_key"}
          label={"SecretID"}
          width="xl"
          placeholder={""}
          rules={[
            {
              required: true,
              message: "请填写SecretID信息!",
            },
          ]}
        />
        <ProFormText
          name={"access_secret"}
          label={"SecretKey"}
          width="xl"
          placeholder={""}
          rules={[
            {
              required: true,
              message: "请填写SecretKey信息!",
            },
          ]}
        />
      </ProForm.Group>;
    case 3:
    case 4:
    case 5: // Sealos Kubeconfig
      return <ProForm.Group>
        <ProFormSelect
          name="regions"
          width="xl"
          label="账户允许部署地区"
          placeholder={""}
          tooltip={"选择云账户允许部署的地区, sealos账户不跨平台"}
          valueEnum={SealosRegionEnum}
          rules={[
            {
              required: true,
              message: "请选择云账户所属的地区!",
            },
          ]}
        />
        <ProFormTextArea
          name={"kube_config"}
          label={"KubeConfig 凭证"}
          placeholder={""}
          width="xl"
          rules={[
            {
              required: true,
              message: "请填写认证信息!",
            },
          ]}
        ></ProFormTextArea>
      </ProForm.Group>
    default:
      return <ProForm.Group/>
  }
};

export const CloudRegionSelector: React.FC<CloudProviderType> = (props) => {
  switch (props.type) {
    case 1:
      return <ProFormSelect
        name="regions"
        width="xl"
        placeholder={""}
        tooltip={"选择云账户允许部署的地区"}
        valueEnum={ALiYunRegionEnum}
        fieldProps={{
          mode: "multiple",
        }}
        rules={[
          {
            required: true,
            message: "请选择云账户所属的地区!",
          },
        ]}
      />
    case 2:
      return <ProFormSelect
        name="regions"
        width="xl"
        placeholder={""}
        tooltip={"选择云账户允许部署的地区"}
        fieldProps={{
          mode: "multiple",
        }}
        valueEnum={TencentRegionEnum}
        rules={[
          {
            required: true,
            message: "请选择云账户所属的地区!",
          },
        ]}
      />
    case 3:
    case 4:
    case 5:
      return <ProFormSelect
        name="regions"
        width="xl"
        placeholder={""}
        tooltip={"选择云账户允许部署的地区"}
        valueEnum={SealosRegionEnum}
        rules={[
          {
            required: true,
            message: "请选择云账户所属的地区!",
          },
        ]}
      />
  }
  return <></>
}

export const CloudRegionOneSelector: React.FC<CloudProviderType> = (props) => {
  switch (props.type) {
    case 0 :
      return <ProFormSelect
        name="region"
        label={"函数所在地区"}
        colProps={{span: 8, offset: 4}}
        tooltip={"选择函数部署的地区, 不要选择账户以外的区域"}
        rules={[
          {
            required: true,
            message: "请选择函数部署的地区!",
          },
        ]}
      />
    case 1:
      return <ProFormSelect
        name="region"
        label={"函数所在地区"}
        colProps={{span: 8, offset: 4}}
        tooltip={"选择函数部署的地区, 不要选择账户以外的区域"}
        valueEnum={ALiYunRegionEnum}
        rules={[
          {
            required: true,
            message: "请选择函数部署的地区!",
          },
        ]}
      />
    case 2:
      return <ProFormSelect
        name="region"
        label={"函数所在地区"}
        colProps={{span: 8, offset: 4}}
        tooltip={"选择函数部署的地区, 不要选择账户以外的区域"}
        valueEnum={TencentRegionEnum}
        rules={[
          {
            required: true,
            message: "请选择函数部署的地区!",
          },
        ]}
      />
    case 3:
    case 4:
    case 5:
      return <ProFormSelect
        name="region"
        label={"函数所在地区"}
        colProps={{span: 8, offset: 4}}
        tooltip={"选择函数部署的地区, 不要选择账户以外的区域"}
        valueEnum={SealosRegionEnum}
        rules={[
          {
            required: true,
            message: "请选择函数部署的地区!",
          },
        ]}
      />
  }
  return <></>
}


export const AuthColumns = [
  // 0
  [],
  // 1
  [
    {
      title: 'AccessId',
      key: 'access_id',
      dataIndex: 'access_id',
      render: () => {
        return "******************"
      }
    },
    {
      title: 'AccessKey',
      key: 'access_key',
      dataIndex: 'access_key',
      render: () => {
        return "******************"
      }
    },
    {
      title: 'AccessSecret',
      key: 'access_secret',
      dataIndex: 'access_secret',
      render: () => {
        return "******************"
      }
    }
  ],
  // 2
  [
    {
      title: 'SecretId',
      key: 'access_key',
      dataIndex: 'access_key',
      render: () => {
        return "******************"
      }
    },
    {
      title: 'SecretKey',
      key: 'access_secret',
      dataIndex: 'access_secret',
      render: () => {
        return "******************"
      }
    }
  ],
  // 3
  [],
  // 4
  [],
  // 5
  [
    {
      title: 'KubeConfig',
      key: 'kube_config',
      dataIndex: 'kube_config',
      valueType: 'textarea',
      render: () => {
        return "******************"
      }
    }
  ],
]
