### How to run the code

1. Clone the repository
2. Open the terminal and navigate to the project folder
3. Create the config file `config.cfg` in the project folder with the following content (you can change the values as needed):
```
[default]
BeginString=FIX.4.4
ConnectionType=initiator
StartDay=Friday
EndDay=Friday
StartTime=00:00:00
EndTime=00:00:00
HeartBtInt=30
ReconnectInterval=5
FileLogPath=./logs
FileStorePath=./data/session
EncryptMethod=0
ResetSeqNumFlag=Y
UseDataDictionary=N
ValidateFieldsOutOfOrder=N
ValidateUnorderedGroupFields=N
ValidateFieldsHaveValues=N
ValidateUserDefinedFields=N
ValidateIncomingMessage=N
AllowUnknownMsgFields=Y
ResetOnLogOut=Y
ResetOnDisconnect=Y
ResetSeqNumFlag=Y


[SESSION]
SocketConnectHost=127.0.0.1
SocketConnectPort=9822

TargetCompID=waanx
SenderCompID=999
```

4. Run the following command to install the required packages:
```
go mod tidy
```

5. Run the following command to build and run the project:
```
make dev-md
```