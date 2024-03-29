import React from "react";
import {ProForm, ProFormText, ProFormSelect, ProFormSwitch} from "@ant-design/pro-components";
import {toNumber} from "lodash";
import {ProxyTypeValueEnum} from "@/enum/service";

export const ProxyForm: React.FC = (props) => {
  return <> <ProForm.Group
    title={"服务参数"}
  >
    <ProFormText
      name="name"
      label="代理名称"
      placeholder={""}
      colProps={{span: 8}}
      rules={[
        {
          required: true,
          message: "请输入代理服务名称!",
        }
      ]}
    />
    <ProFormText
      name="listen_address"
      label="监听地址"
      placeholder={""}
      colProps={{span: 8,offset:4}}
      rules={[
        {
          required: true,
          message: "请输入合法的监听地址!",
          pattern: RegExp(""),
        },
        {
          pattern: /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/,
          message: "请输入合理的IP地址",
        }
      ]}
    />
    <ProFormText
      name="listen_port"
      label="监听端口"
      placeholder={""}
      colProps={{span: 8}}
      rules={[
        {
          required: true,
          message: "请输入监听端口!",
        },
        {
          validator: (rule, value) => {
            const v = toNumber(value);
            if (v >= 65535 ||v <= 0 || isNaN(v) ) {
              return Promise.reject(new Error('请输入合法端口号(1-65535)'));
            }
            return Promise.resolve();
          },
        },
      ]}
    />
    <ProFormSelect
      name="type"
      colProps={{span: 8,offset:4}}
      label="监听协议"
      placeholder={""}
      valueEnum={ProxyTypeValueEnum}
      rules={[
        {
          required: true,
          message: "请选择监听的服务类型!",
        },
      ]}
    />
  </ProForm.Group>
    <ProForm.Group
      title={"高级选项"}
    >
      <ProFormSwitch
        name="tor"
        label={"开启 Tor 网桥"}
        tooltip={"选择开启 Tor 网桥, 创建或选择的函数必须也对应开启 Tor 标识"}
        checkedChildren={"开启"}
        unCheckedChildren={"关闭"}
        colProps={{
          span: 12,
        }}
      />
    </ProForm.Group>
  </>
}
