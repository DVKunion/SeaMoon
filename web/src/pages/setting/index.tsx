import React, {useState, useRef} from "react";
import {PageContainer, ProCard, ProForm, ProFormText, ProFormInstance, ProFormSwitch} from "@ant-design/pro-components";
import {Button, message, Space, Tag} from "antd"
import {getSysConfig} from "@/services/setting/api";
import {GithubOutlined} from "@ant-design/icons";
import {updatePasswd} from "@/services/user/api";
import {handleUpdateSysConfig} from "@/pages/setting/hanlder";

const Setting: React.FC = () => {

  const formRef = useRef<ProFormInstance>();
  const [version, setVersion] = useState("")

  return <PageContainer
    title={<Space>系统配置 <Tag icon={<GithubOutlined/>} color="#76b39d">{version}</Tag></Space>}
  >
    <ProCard>
      <ProForm
        formRef={formRef}
        submitter={false}
        onFinish={async (values) => {
          if (values["admin_password"] !== "" && values["admin_password"] != undefined) {
            const res = await updatePasswd(values["admin_password"]);
            if (res.success) {
              message.success('修改密码成功', 1);
            } else {
              message.error(res.code + ":" + res.msg, 1)
            }
          }
          await handleUpdateSysConfig(values);
          values.auto_start = values.auto_start == "true";
          formRef?.current?.setFieldsValue(values);
        }}
        params={{}}
        request={async () => {
          const {data} = await getSysConfig();
          if (data.auto_start === "true") {
            // @ts-ignore
            data.auto_start = true;
          } else {
            // @ts-ignore
            data.auto_start = false;
          }
          setVersion(data.version);
          return data;
        }}
      >
        <ProForm.Group title={"HTTP 管理服务"}>
          <ProFormText
            name="control_addr"
            label="监听地址"
            tooltip={"如果你是通过 docker 启动的, 请不要修改此配置，否则可能会造成服务无法访问"}
            width={"md"}
            placeholder={"e.g.: 0.0.0.0"}
          />
          <ProFormText
            name="control_port"
            label="监听端口"
            tooltip={"如果你是通过 docker 启动的, 修改管理端口后，docker的端口映射也需要一起改变"}
            width={"md"}
            placeholder={"e.g.: 7777"}
          />
          <ProFormText
            name="control_log"
            label="服务日志"
            width={"lg"}
            placeholder={"e.g.: .seamoon.log"}
          />
        </ProForm.Group>
        <ProForm.Group title={"其他配置"}>
          <ProFormSwitch
            name="auto_start"
            label="自动运行"
            tooltip={"当服务重新启动时，自动启动所有运行状态下的代理, 否则重启将会自动停止所有代理服务"}
            width={"xs"}
            fieldProps={
              {
                checkedChildren: "开启",
                unCheckedChildren: "关闭",
              }
            }
            style={{display: 'flex', alignItems: 'center'}}
          />
        </ProForm.Group>
        <ProForm.Group title={"账户认证"}>
          <ProFormText.Password
            name="admin_password"
            label="管理密码"
            tooltip={"修改管理后台的登陆密码"}
            width={"md"}
            placeholder={""}
          />
        </ProForm.Group>
        <ProForm.Item>
          <Button type="primary" htmlType="submit" style={{float: "right", marginRight: "4em"}}>
            保存
          </Button>
        </ProForm.Item>
      </ProForm>
    </ProCard>
  </PageContainer>
}

export default Setting
