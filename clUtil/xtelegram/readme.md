### 使用教程
    1.将机器人添加到group群组
    2.任意一个账号发信息给机器人 :  @ttbc_bot /start 开始工作吧
    3.使用机器人的token，去读群组chatid:
        https://api.telegram.org/bot{telegram robot token}/getUpdates
        注意: 如果拿到的结果是: {"ok":true,"result":[]} 或则没有chatid，则在群组group删除机器人，重新添加 
    4.群组chatid: result[0]['my_chat_member']['chat']['id']
    {
        "ok": true,
        "result": [
            {
                "update_id": 362059980,
                "my_chat_member": {
                    "chat": {
                        "id": -4533149931,
                        "title": "\u6d4b\u8bd5BC\u673a\u5668\u4eba",
                        "type": "group",
                        "all_members_are_administrators": false
                    },
                    "from": {
                        "id": 6403436826,
                        "is_bot": false,
                        "first_name": "xi",
                        "last_name": "chen"
                    },
                    "date": 1727674132,
                    "old_chat_member": {
                        "user": {
                            "id": 7454279271,
                            "is_bot": true,
                            "first_name": "tcpBot",
                            "username": "ttbc_bot"
                        },
                        "status": "member"
                    },
                    "new_chat_member": {
                        "user": {
                            "id": 7454279271,
                            "is_bot": true,
                            "first_name": "tcpBot",
                            "username": "ttbc_bot"
                        },
                        "status": "left"
                    }
                }
            },
            {
                "update_id": 362059981,
                "message": {
                    "message_id": 4,
                    "from": {
                        "id": 6403436826,
                        "is_bot": false,
                        "first_name": "xi",
                        "last_name": "chen"
                    },
                    "chat": {
                        "id": -4533149931,
                        "title": "\u6d4b\u8bd5BC\u673a\u5668\u4eba",
                        "type": "group",
                        "all_members_are_administrators": false
                    },
                    "date": 1727674132,
                    "left_chat_participant": {
                        "id": 7454279271,
                        "is_bot": true,
                        "first_name": "tcpBot",
                        "username": "ttbc_bot"
                    },
                    "left_chat_member": {
                        "id": 7454279271,
                        "is_bot": true,
                        "first_name": "tcpBot",
                        "username": "ttbc_bot"
                    }
                }
            },
            {
                "update_id": 362059982,
                "my_chat_member": {
                    "chat": {
                        "id": -4533149931,
                        "title": "\u6d4b\u8bd5BC\u673a\u5668\u4eba",
                        "type": "group",
                        "all_members_are_administrators": true
                    },
                    "from": {
                        "id": 6403436826,
                        "is_bot": false,
                        "first_name": "xi",
                        "last_name": "chen"
                    },
                    "date": 1727674176,
                    "old_chat_member": {
                        "user": {
                            "id": 7454279271,
                            "is_bot": true,
                            "first_name": "tcpBot",
                            "username": "ttbc_bot"
                        },
                        "status": "left"
                    },
                    "new_chat_member": {
                        "user": {
                            "id": 7454279271,
                            "is_bot": true,
                            "first_name": "tcpBot",
                            "username": "ttbc_bot"
                        },
                        "status": "member"
                    }
                }
            },
            {
                "update_id": 362059983,
                "message": {
                    "message_id": 5,
                    "from": {
                        "id": 6403436826,
                        "is_bot": false,
                        "first_name": "xi",
                        "last_name": "chen"
                    },
                    "chat": {
                        "id": -4533149931,
                        "title": "\u6d4b\u8bd5BC\u673a\u5668\u4eba",
                        "type": "group",
                        "all_members_are_administrators": true
                    },
                    "date": 1727674176,
                    "new_chat_participant": {
                        "id": 7454279271,
                        "is_bot": true,
                        "first_name": "tcpBot",
                        "username": "ttbc_bot"
                    },
                    "new_chat_member": {
                        "id": 7454279271,
                        "is_bot": true,
                        "first_name": "tcpBot",
                        "username": "ttbc_bot"
                    },
                    "new_chat_members": [
                        {
                            "id": 7454279271,
                            "is_bot": true,
                            "first_name": "tcpBot",
                            "username": "ttbc_bot"
                        }
                    ]
                }
            }
        ]
    }
        
