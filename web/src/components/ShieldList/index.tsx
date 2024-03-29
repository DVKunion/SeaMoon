import React from 'react';
import {Space} from "antd";

const ShieldList: React.FC = ({children}) => (
  <Space>
    <img src="https://goreportcard.com/badge/github.com/DVKunion/SeaMoon" alt="go-report"/>
    <img src="https://img.shields.io/github/languages/top/DVKunion/SeaMoon.svg?&color=blueviolet"
         alt="languages"/>
    <img src="https://img.shields.io/badge/LICENSE-MIT-777777.svg" alt="license"/>
    <img src="https://img.shields.io/github/downloads/dvkunion/seamoon/total?color=orange" alt="downloads"/>
    <img src="https://img.shields.io/github/stars/DVKunion/SeaMoon.svg" alt="stars"/>
  </Space>
)

export default ShieldList;
