import React from "react";
import {ProForm, ProFormSwitch, ProFormText, ProFormSelect} from "@ant-design/pro-components";
import {toNumber} from "lodash";
import {TunnelAuthFCTypeEnum, TunnelTypeValueEnum} from "@/enum/tunnel";
import {CloudRegionOneSelector} from "@/pages/provider/components/AuthForm";

export type TunnelFormProps = {
  type: number
}

export const TunnelForm: React.FC<TunnelFormProps> = (props) => {
  return <> <ProForm.Group
    title={"基础信息"}
  >
    <ProFormText
      name="tunnel_name"
      label={"函数名称"}
      colProps={{span: 8}}
      rules={[
        {
          required: true,
          message: "请填写函数名称!",
        },
        {
          pattern: /^[a-zA-Z0-9_]+$/,
          message: "代理服务只能为英文/数字/下划线!",
        },
        {
          max: 24,
          message: "代理名称不要超过60个字符",
        }
      ]
      }
    />
    <CloudRegionOneSelector type={props.type} />
  </ProForm.Group>
    <ProForm.Group
      title={"函数规格"}
    >
      <ProFormText
        name="cpu"
        label={"CPU限制"}
        colProps={{span: 8}}
        tooltip={"必须为 0.05 的倍数"}
        placeholder={""}
        rules={[
          {
            required: true,
            message: "请填写cpu资源限制!",
          },
          {
            validator: (rule, value) => {
              if (toNumber(value) < 0.05) {
                return Promise.reject(new Error('最低 cpu 数值为 0.05'));
              }
              if (toNumber(value) * 100 % 5 !== 0) {
                return Promise.reject(new Error('必须是 0.05 的倍数'));
              }
              return Promise.resolve();
            }
          }
        ]}
      />
      <ProFormText
        name="memory"
        label={"内存限制"}
        colProps={{span: 8, offset: 4}}
        placeholder={""}
        rules={[
          {
            required: true,
            message: "请填写内存资源限制!",
          },
          {
            validator: (rule, value) => {
              if (toNumber(value) < 64) {
                return Promise.reject(new Error('最低内存数值为64'));
              }
              return Promise.resolve();
            }
          }
        ]}
      />
      <ProFormText
        name="instance"
        label={"最大处理数"}
        tooltip={"表示一个实例最大能够并发处理的请求数"}
        colProps={{span: 8}}
        placeholder={""}
        rules={[
          {
            required: true,
            message: "请填写最大实例处理数!",
          },
          {
            validator: (rule, value) => {
              if (toNumber(value) < 1) {
                return Promise.reject(new Error('处理数最低为1'));
              }
              return Promise.resolve();
            }
          }
        ]}
      />
      <ProFormText
        name="port"
        label={"端口号配置"}
        tooltip={"自定义配置实例服务端口号"}
        colProps={{span: 8, offset: 4}}
        width={"md"}
        placeholder={""}
        rules={[
          {
            required: true,
            message: "请填写正确的端口号!",
          },
          {
            validator: (rule, value) => {
              if (toNumber(value) >= 65535 || toNumber(value) <= 0) {
                return Promise.reject(new Error('请输入合法端口'));
              }
              return Promise.resolve();
            },
          }
        ]}
      />
      <ProFormSelect
        name="tunnel_auth_type"
        label={"函数认证方式"}
        tooltip={"云函数自身提供认证方式，配置该项可防止被刷"}
        colProps={{span: 8}}
        placeholder={""}
        valueEnum={TunnelAuthFCTypeEnum}
        showSearch={true}
        rules={[
          {
            required: true,
            message: "请选择函数的认证方式!",
          },
        ]}
      />
      <ProFormSelect
        name="tunnel_type"
        label={"隧道协议类型"}
        colProps={{span: 8, offset: 4}}
        placeholder={""}
        valueEnum={TunnelTypeValueEnum}
        rules={[
          {
            required: true,
            message: "请选择隧道协议类型!",
          },
        ]}
      />
    </ProForm.Group>
    <ProForm.Group title={"高级选项"} grid={true} rowProps={{
      gutter: [16, 16],
    }}>
      <ProFormSwitch
        name="tls"
        label={"开启 TLS"}
        checkedChildren={"开启"}
        unCheckedChildren={"关闭"}
        colProps={{
          span: 12,
        }}
      />
      <ProFormSwitch
        name="tor"
        label={"开启 Tor 网桥"}
        tooltip={"开启 Tor 模式会导致内存资源使用增多"}
        checkedChildren={"开启"}
        unCheckedChildren={"关闭"}
        colProps={{
          span: 12,
        }}
      />
    </ProForm.Group>
  </>
}
