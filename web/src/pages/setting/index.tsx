import React, {useRef} from "react";
import type { ProFormInstance} from "@ant-design/pro-components";
import {PageContainer, ProCard, ProForm, ProFormSwitch, ProFormText} from "@ant-design/pro-components";
import {Button} from "antd"
import {getSysConfig} from "@/services/setting/api";
import {handleUpdateSysConfig} from "@/pages/setting/hanlder";

const Setting: React.FC = () => {

  const formRef = useRef<ProFormInstance>();

  return <PageContainer
    title="系统配置"
  >
    <ProCard>
      <ProForm
        grid={true}
        rowProps={{
          gutter: [16, 16],
        }}
        layout={"horizontal"}
        formRef={formRef}
        submitter={{
          render: (props, doms) => {
            return <Button type="primary" htmlType="submit" style={{float: "right", marginRight: "4em"}}>
              保存
            </Button>
          },
        }}
        onFinish={async (values) => {
          await handleUpdateSysConfig(values);
          values.auto_start = values.auto_start == "true";
          values.auto_sync = values.auto_sync == "true";
          formRef?.current?.setFieldsValue(values);
        }}
        params={{}}
        request={async () => {
          const {data} = await getSysConfig();
          data.auto_start = String(data.auto_start);
          data.auto_sync = String(data.auto_sync);
          return data;
        }}
      >
        <ProForm.Group
          title={"Backend API 管理服务"}>
          <ProFormText
            name="control_addr"
            label="监听地址"
            tooltip={"如果你是通过 docker 启动的, 请不要修改此配置，否则可能会造成服务无法访问"}
            colProps={{span: 8}}
            placeholder={"e.g.: 0.0.0.0"}
          />
          <ProFormText
            name="control_port"
            label="监听端口"
            tooltip={"如果你是通过 docker 启动的, 修改管理端口后，docker的端口映射也需要一起改变"}
            colProps={{span: 8, offset: 4}}
            placeholder={"e.g.: 7777"}
          />
          <ProFormText
            name="control_log"
            label="日志路径"
            tooltip={"系统日志存放的路径"}
            colProps={{span: 8}}
            placeholder={"e.g.: .seamoon.log"}
          />
        </ProForm.Group>
        <ProForm.Group
          title={"Daemon 服务"}>
          <ProFormText
            name="xray_port"
            label="xray service 监听端口"
            tooltip={"如果你是通过 docker 启动的, 请不要修改此配置，否则可能会造成服务无法访问"}
            colProps={{span: 8}}
            placeholder={"e.g.: 10085"}
          />
        </ProForm.Group>
        <ProForm.Group title={"其他配置"}>
          <ProFormSwitch
            name="auto_start"
            label="自动运行"
            tooltip={"当服务重新启动时，自动启动所有运行状态下的代理, 否则重启将会自动停止所有代理服务"}
            colProps={{span: 24}}
            fieldProps={
              {
                checkedChildren: "开启",
                unCheckedChildren: "关闭",
              }
            }
            style={{display: 'flex', alignItems: 'center'}}
          />
          <ProFormSwitch
            name="auto_sync"
            label="自动同步"
            tooltip={"当服务重新启动时，自动同步各账户远程信息来保证数据同步"}
            colProps={{span: 24}}
            fieldProps={
              {
                checkedChildren: "开启",
                unCheckedChildren: "关闭",
              }
            }
            style={{display: 'flex', alignItems: 'center'}}
          />
        </ProForm.Group>
      </ProForm>
    </ProCard>
  </PageContainer>
}

export default Setting
