# Alertmanager webhook for Telegram (GO Version)

Python Version (https://github.com/nopp/alertmanager-webhook-telegram-python) 

## INSTALL

* go get -d .

Running on docker
=================
    git clone https://github.com/nopp/alertmanager-webhook-telegram-go.git
    cd alertmanager-webhook-telegram-go/docker/
    docker build -t awt-go:0.1 .

    docker run -d --name telegram-bot \
    	-e "bottoken=telegramBotToken" \
    	-e "chatid=telegramChatID" \
    	-p 9229:9229 awt-go:0.1
