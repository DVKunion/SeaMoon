import {DefaultFooter} from '@ant-design/pro-components';

const Footer: React.FC = () => {

  const currentYear = new Date().getFullYear();

  return (
    <DefaultFooter
      copyright={`${currentYear} ${"DVK"}`}
    />
  );
};

export default Footer;
