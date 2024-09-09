import {ModalForm, ProForm, ProFormText} from '@ant-design/pro-components';
import React from 'react';

export type FormValueType = {
  old_password: string,
  password: string,
  password_repeat: string,
}

export type UpdateFormProps = {
  onSubmit: (values: FormValueType) => Promise<boolean>;
};

const UpdateForm: React.FC<UpdateFormProps> = (props) => {

  return (
    <ModalForm<FormValueType>
      width={640}
      title={"修改密码"}
      trigger={
        <a type="primary" key="update">
          修改密码
        </a>
      }
      autoFocusFirstInput
      modalProps={{
        destroyOnClose: true,
      }}
      onFinish={props.onSubmit}
    >
      <ProForm.Group>
        <ProFormText.Password
          name="old_password"
          label="旧密码"
          key={"old_password"}
          placeholder={""}
          width="xl"
          rules={[
            {
              required: true,
              message: "请输入旧密码!",
            },
          ]}
        />
        <ProFormText.Password
          name="password"
          label="新密码"
          key={"password"}
          placeholder={""}
          width="xl"
          rules={[
            {
              required: true,
              message: "请输入修密码!",
            },
          ]}
        />
        <ProFormText.Password
          name="password_repeat"
          label="再次确认新密码"
          key={"password_repeat"}
          placeholder={""}
          width="xl"
          rules={[
            {
              required: true,
              message: "请输入修密码!",
            },
          ]}
        />
        </ProForm.Group>
    </ModalForm>
);
};

export default UpdateForm;
