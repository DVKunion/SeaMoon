import {
  ProFormSelect,
  ProFormText, ProFormTextArea, ProFormDigit,
  StepsForm, ProForm,
} from '@ant-design/pro-components';
import {Modal} from 'antd';
import React, {useState} from 'react';
import {CloudProvideTypeValueEnum} from "@/enum/cloud";
import {CloudProviderAuthForm} from "@/pages/provider/components/AuthForm";
import {toNumber} from "lodash";

export type FormValueType = {
  access_id?: string,
  access_key?: string,
  access_secret?: string,
  token?: string,
  kube_config?: string,
} & Partial<Cloud.Provider>;

export type CreateFormProps = {
  onCancel: (flag?: boolean, formValue?: FormValueType) => void;
  onSubmit: (values: FormValueType) => Promise<void>;
  createModalVisible: boolean;
  values: Partial<Cloud.Provider>;
};

const CreateForm: React.FC<CreateFormProps> = (props) => {

  const [cloudType, setCloudType] = useState<number>(0);

  return (
    <StepsForm
      stepsProps={{
        size: 'small',
      }}
      stepsFormRender={(dom, submitter) => {
        return (
          <Modal
            width={640}
            bodyStyle={{padding: '32px 40px 48px'}}
            destroyOnClose
            title={"新增云账户"}
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
      onFinish={props.onSubmit}
    >
      <StepsForm.StepForm
        title={'基本信息'}
        onFinish={async (value) => {
          setCloudType(toNumber(value["type"]));
          return true;
        }}
      >
        <ProForm.Group>
          <ProFormText
            name="name"
            label="账户名称"
            placeholder={""}
            tooltip={"此处的账户名称仅用于平台展示，与实际云平台的账户名称无关"}
            width="xl"
            rules={[
              {
                required: true,
                message: "请输入账户名称!",
              },
            ]}
          />
          <ProFormSelect
            name="type"
            width="xl"
            label="账户类型"
            placeholder={""}
            tooltip={"选择云账户来源厂商"}
            valueEnum={CloudProvideTypeValueEnum}
            rules={[
              {
                required: true,
                message: "请选择账户类型!",
              },
            ]}
          />
          <ProFormDigit
            name="max_limit"
            width="xl"
            label="最大部署限制"
            min={0}
            tooltip={"设置该账户最大允许部署的函数数量, 0 表示无限制, 默认为 0"}
            placeholder={""}
          />
          <ProFormTextArea
            name="desc"
            label="账户描述"
            placeholder={""}
            tooltip={"用于备注账户,方便同一平台账户的区分"}
            width="xl"
          />
        </ProForm.Group>
      </StepsForm.StepForm>
      <StepsForm.StepForm
        title={"认证信息"}
      >
        <CloudProviderAuthForm type={cloudType}/>
      </StepsForm.StepForm>
    </StepsForm>
  );
};

export default CreateForm;
