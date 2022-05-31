# SMSFW SMPP Client

Features:

- Receive `DeliverSM` from SMSFW for blacklisted SMS
- Send notification to SMS originator, telling them they have been blocked.

TBD:

- Store the blocked MSISDN in a __cockroachDB__ for reference later
- Connect to 3rd party for sync blacklist numbers into DB.


## Deployment

- App folder: `/apps/smsfw-smpp-client`
- Application Configuration: `/apps/smsfw-smpp-client/config/application.yaml`
- Systemd services: `/etc/systemd/system/smsfw-smpp-client.service`
