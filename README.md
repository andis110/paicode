# Gamepai：游派区块链系统

## Introduction 简介

Gamepai is a APP for enjoying mobile gaming and social networking. It use a blockchain-base crypto-current called "pai" for in-game purchasing, exchanging and many other activities. [You can check more detail here](http://gamecenter.mobi)

游派是一个手机游戏和社交应用。它采用基于区块链技术的加密货币“PAI”实现游戏内购，道具交易等多种活动。更多详情请参考[这里](http://gamecenter.mobi)

## The blockchain 关于区块链

The blockchain is built with a forked from fabric 0.6, the hyperledger project. paicode is the project for the chaincode and client ("wallet") part

游派区块链是在超级账本（hyperledger）项目的fabric 0.6的改进版本上构建的。本项目（paicode）是其中的chaincode和客户端（即钱包）部分

## Installation 安装

You need go 1.7. To install the client (the so-called "thin wallet" for a crypto-current), just call:

项目开发使用go 1.7。安装项目的客户端（即“轻钱包”部分），使用下述命令：


```
go get -insecure gamecenter.mobi/paicode/client/gamepaicore
```

## Configuration and running 配置和运行

Configuration file is needed to run the client (gamepaicore). It should be core.windows.yaml under windows or core.yaml for other platforms, and be put under the working directory of client (gamepaicore) process

TLS is forced to access the gamepai services. CA certification for TLS is in misc/ca.crt and must be placed as specified in the peer.tls.rootcert field of the configuration file

客户端（gamepaicore）运行时需要在进程的工作目录下包含配置文件，对windows平台，配置文件是core.windows.yaml，其它平台是core.yaml。配置文件的模板可在misc目录中找到。

客户端访问当前已运行的游派节点要求使用TLS连接，要求将服务端的根证书（misc/ca.crt）放置到配置文件中的指定位置（peer.tls.rootcert）


