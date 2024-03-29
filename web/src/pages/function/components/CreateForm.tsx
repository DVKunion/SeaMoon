import React, {useState} from 'react';
import {Modal} from 'antd';
import {StepsForm} from '@ant-design/pro-components';
import {ProviderSelect} from "@/components/StepForm/ProviderSelect";
import {TunnelForm} from "@/components/StepForm/TunnelForm";

export type FormValueType = {
  cpu: number,
  region: string,
  memory: number,
  instance: number,
  tunnel_auth_type: number,
  tls: boolean,
  tor: boolean,
  tunnel_name: string,
  tunnel_type: string,
} & Partial<Serverless.Tunnel>;

export type CreateFormProps = {
  onCancel: (flag?: boolean, formValue?: FormValueType) => void;
  onSubmit: (values: FormValueType) => Promise<void>;
  createModalVisible: boolean;
  values: Partial<Serverless.Tunnel>;
};

const CreateForm: React.FC<CreateFormProps> = (props) => {

  const [current, setCurrent] = useState<number>(0);
  const [type, setType] = useState<number>(0);

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
            width={720}
            bodyStyle={{padding: '32px 40px 60px'}}
            destroyOnClose
            title={"新增函数"}
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
        // @ts-ignore
        return props.onSubmit(values);
      }}
    >
      <StepsForm.StepForm
        title={"实例选择"}
      >
        <ProviderSelect onChange={(values) => {
          setType(values);
        }}/>
      </StepsForm.StepForm>
      <StepsForm.StepForm
        grid={true}
        rowProps={{
          gutter: [16, 16],
        }}
        title={"函数配置"}
        initialValues={{
          cpu: type === 5 ? '0.1' : '0.05',
          memory: type === 5 ? '64' : '128',
          tunnel_auth_type: 1,
          instance: 1,
          port: 9000,
          tunnel_type: "websocket",
          tls: true,
          tor: false,
        }}
      >
        <TunnelForm type={type}/>
      </StepsForm.StepForm>
    </StepsForm>
  );
};

export default CreateForm;
