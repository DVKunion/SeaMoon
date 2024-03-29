import {
  CheckCard,
  StepsForm, ProForm,
} from '@ant-design/pro-components';
import {Modal, Avatar} from 'antd';
import React, {useState} from 'react';
import {toNumber} from "lodash";
import IconFont from "@/components/IconFont";
import {TunnelForm} from "@/components/StepForm/TunnelForm";
import {ProviderSelect} from "@/components/StepForm/ProviderSelect";
import {ProxyForm} from "@/components/StepForm/ProxyForm";
import {TunnelSelect} from "@/components/StepForm/TunnelSelect";

export type FormValueType = {
  cpu: number,
  port: string,
  memory: number,
  instance: number,
  tunnel_auth_type: number,
  tls: boolean,
  tor: boolean,
  tunnel_type: string,
  tunnel_name: string,
  tunnel_id: number,
  provider_id: number,
} & Partial<Service.Proxy>;

export type CreateFormProps = {
  onCancel: (flag?: boolean, formValue?: FormValueType) => void;
  onSubmit: (values: FormValueType) => Promise<void>;
  createModalVisible: boolean;
  values: Partial<Service.Proxy> & Partial<Serverless.Tunnel>;
};

const CreateForm: React.FC<CreateFormProps> = (props) => {

  const [current, setCurrent] = useState<number>(0);
  const [deploy, setDeploy] = useState<number>(1);
  const [type, setType] = useState<number>(0);
  const [tor, setTor] = useState<boolean>(false);

  return (
    <StepsForm
      stepsProps={{
        size: 'small',
      }}
      current={current}
      onCurrentChange={(c) => {
        setCurrent(c);
      }}
      stepsFormRender={(dom, submitter) => {
        return (
          <Modal
            width={740}
            bodyStyle={{padding: '32px 40px 60px'}}
            destroyOnClose
            title={"新增服务"}
            open={props.createModalVisible}
            footer={submitter}
            onCancel={() => {
              props.onCancel();
            }}
          >
            {dom}
          </Modal>
        );
      }}
      onFinish={(values) => {
        setCurrent(0);
        setDeploy(1);
        // @ts-ignore
        return props.onSubmit(values);
      }}
    >
      <StepsForm.StepForm
        title={'基本配置'}
        grid={true}
        rowProps={{
          gutter: [16, 16],
        }}
        onFinish={async (values) => {
          setTor(values.tor);
          return true
        }}
      >
        <ProxyForm/>
      </StepsForm.StepForm>
      <StepsForm.StepForm
        title={"选择关联类型"}
      >
        <CheckCard.Group style={{width: '100%'}}
                         defaultValue={deploy}
                         onChange={(value) => {
                           setDeploy(toNumber(value));
                         }}>
          <CheckCard
            title="选择已有的函数进行关联"
            avatar={
              <Avatar
                style={{backgroundColor: '#ffffff'}}
                icon={<IconFont type={"icon-saeServerlessyingyongyinqing1"}/>}
                size="large"
              />
            }
            description="从已有的函数实例中进行关联选择，建立隧道。"
            value={1}
          />
          <CheckCard
            title="选择账户并自动创建新实例"
            avatar={
              <Avatar
                style={{backgroundColor: '#ffffff'}}
                icon={<IconFont type={"icon-cloud1"}/>}
                size="large"
              />
            }
            description="通关关联账户自动创建对应协议的全新函数实例。"
            value={2}
          />
        </CheckCard.Group>
      </StepsForm.StepForm>
      <StepsForm.StepForm
        title={"实例选择"}
      >
        <ProForm.Group
        >
          {deploy === 1 ?
            <TunnelSelect values={props.values} tor={tor}/> : <ProviderSelect onChange={(values) => {
              setType(values);
            }}/>}
        </ProForm.Group>
      </StepsForm.StepForm>
      <StepsForm.StepForm
        title={"高级配置"}
        grid={true}
        rowProps={{
          gutter: [16, 16],
        }}
        initialValues={{
          cpu: '0.05',
          memory: '128',
          tunnel_auth_type: 1,
          instance: 5,
          port: 9000,
          tunnel_type: "websocket",
          tls: true,
          tor: false,
        }}
      >
        {deploy === 1 ? "关联函数无法进行高级配置, 请在对应函数实例页面进行修改。" :
          <TunnelForm type={type}/>
        }
      </StepsForm.StepForm>
    </StepsForm>
  );
};

export default CreateForm;
