# Endpoints

## Access Tokens

`POST /access-tokens`

Generates an access token needed to consume the other endpoints.

Possible requests:

| Field | Type | Description | Optional | Example |
| ----- | ---- | ----------- | -------- | ------- |
| clientId | string | Client Id | No | `AZR1E7G...` |
| clientSecret | string | Client's secret key | No | `e508698...` |
| scope | string | Claims scope | Yes | `admin` |

| Field | Type | Description | Optional | Example |
| ----- | ---- | ----------- | -------- | ------- |
| clientId | string | Client Id | No | `AZR1E7G...` |
| privateKey | string | Client's private key | No | `ba632fe...` |
| scope | string | Claims scope | Yes | `admin` |

| Field | Type | Description | Optional | Example |
| ----- | ---- | ----------- | -------- | ------- |
| token | string | Auth token | No | `3fa1ef7...` |
| scope | string | Claims scope | Yes | `admin` |

Response: 

| Field | Type | Description | Example |
| ----- | ---- | ----------- | ------- |
| accessToken | string | Access token | `eyJhbGciOiJSUzI1NiIsIm...` |
| refreshToken | string | Refresh token | `AEu4IL2TDK_RJy7vFqoYVY...` |
| expiresIn | number | Access token's lifetime | `3600000` |

## Deposit Validations

`POST /deposit-validations`

Validates if the destination account can receive the incoming deposit which has not been processed yet.

Request: 

| Field | Type | Description | Optional | Example |
| ----- | ---- | ----------- | -------- | ------- |
| amount | number | Deposit amount multiplied by 100, $1.000 would be 100000 | No | `100000` |
| destinationAccountNumber | string | Destination account number, 999+NationalId without verifier digit | No | `0000000099918166392` |
| destinationAccountTypeId | string | Destination account type represented by id | No | `40` |
| destinationRut | string | National id of the user that will eventually receive the deposit | No | `000181663928` |
| message | string | Message associated with the deposit | Yes | `test msg` |
| originAccountNumber | string | Origin account number | No | `000000000019157907` |
| originAccountTypeId | string | Origin account type represented by id | No | `20` |
| originBankId | string | Origin bank represented by id | No | `0001` |
| originRut | string | National id of the user that made the deposit | No | `0000176988703` |
| receivedAt | string | Timestamp in which the deposit was received | No | `191003103534` |
| traceNumber | string | Trace number associated to the deposit | No | `000000000001` |

Response: 

| Field | Type | Description | Example |
| ----- | ---- | ----------- | ------- |
| code | string | Code with a summary of the validation result | `valid` |
| message | string | Message explaining the validation result | `Deposit is valid` |
| valid | boolean | Returns true indicating that the deposit is valid, otherwise false | `true` |

## Received Transfers

`POST /received-transfers`

Creates a receivedTransfer to be processed. 

Request: 

| Field | Type | Description | Optional | Example |
| ----- | ---- | ----------- | ------- |  ------- |
| amount | number | Deposit amount multiplied by 100, $1.000 would be 100000 | No | `100000` |
| destinationAccountNumber | string | Destination account number, 999+NationalId without verifier digit | No | `0000000099918166392` |
| destinationAccountTypeId | string | Destination account type represented by id | No | `40` |
| destinationRut | string | National id of the user that will eventually receive the deposit | No | `000181663928` |
| message | string | Message associated with the deposit | Yes | `test msg` |
| originAccountNumber | string | Origin account number | No | `000000000019157907` |
| originAccountTypeId | string | Origin account type represented by id | No | `20` |
| originBankId | string | Origin bank represented by id | No | `0001` |
| originRut | string | National id of the user that made the deposit | No | `0000176988703` |
| receivedAt | string | Timestamp in which the deposit was received | No | `191003103534` |
| traceNumber | string | Trace number associated to the deposit | No | `000000000001` |

Response: 

| Field | Type | Description | Example |
| ----- | ---- | ----------- | ------- |
| code | string | Code with a summary of the received transfer creation result | `ok` |
| message | string | Message explaining the received transfer creation result | `Received transfer was processed successfully.` |
| id | boolean | Received transfer id | `191003103535*000000000001*0001` |

## Reversed Transfers

`POST /reversed-transfers`

Begins a reverse procedure of a received transfer.

Request:

| Field | Type | Description | Optional | Example |
| ----- | ---- | ----------- | -------- | ------- |
| originBankId | string | Origin bank represented by id | No | `0001` |
| receivedAt | string | Timestamp in which the deposit was received | No | `191003103534` |
| traceNumber | string | Trace number associated to the deposit | No | `000000000001` |

Response: 

| Field | Type | Description | Example |
| ----- | ---- | ----------- | ------- |
| code | string | Code with a summary of the reverse procedure creation result | `success` |
| message | string | Message explaining the reverse procedure creation result | `Transfer was processed successfully` |
| id | boolean | Received transfer id | `191003103535*000000000001*0001` |

# Error Codes

| Code | Description |
| ---- | ----------- |
| 01 | Error decoding request parameters |
| 02 | Error checking request parameters |
| 03 | Error getting client claims |
| 04 | Error creating custom token with claims |
| 05 | Error verifying custom token |
| 06 | Error getting user snapshot |
| 07 | Error getting user data |
| 08 | Error getting account |
| 09 | Error sending over max balance push notification |
| 10 | Error checking max daily bank transfer deposit amount rule |
| 11 | Error sending over max daily bank transfer deposit amount push notification |
| 12 | Error checking max monthly bank transfer deposit amount rule |
| 13 | Error sending over max monthly bank transfer deposit amount push notification |
| 14 | Error decoding auth token |
| 15 | Error verifying authentication |
| 16 | Error starting transaction |
| 17 | Error reading in transaction |
| 18 | Error writing in transaction |
