# Changelog

## 1.3.1 - 2021-02-23

### Description

- Fixed error reporting log arguments
- Accounts max balances with null values now are handled as 0
- Environment name added to server index

### Changelog

ca2b96b4f4ad37029ac40542df4a821fd9a26cc0 - Version bump —version 1.3.1
072c7463ae06815d60f5147acae7e70aeb85ec1d - Environment name added to server index
80e5fa2fbfee769ce2ffd3f3707d695d9710f5b8 - Now null maxBalances are handled as value 0
fda2583ddb5e1a89984ead5c4d964cd1c95bc1ca - Fixed error reporting log arguments

## 1.3.0 - 2021-01-05

### Description

- gRPC accounts engine service added to get account information
- Major code refactor, including logs, type usages and filenames
- Sonarqube scan analysis added
- Initial API test setup and deposit validation endpoint tests added
- Error reporting to GCP fixed

### Changelog

1f74d83f01cf01abc5e362ec55eff545a2d79a6f - Version bump —version 1.3.0
3c80fb390771401eb661f75f6bd4e8b98c6e70f4 - Fixed the way external service error are handled
d5b62db84460327d15197e55449b6f723de7e715 - Removed message param from LogEntryWithId method
59b17422c623ccd322ce028a18b92bc7a3c81e3a - Add user id to error reporting in order to track report logs
ffbef4c5a99781429ae6abdb39debfc6527570e2 - Error client defer removed
a9a5401857d59eb97558957833e18cabbc72a583 - Error log referencing fcm token result query changed from error to warn severity level
b270d774d9ba7fc882bbbcbe9c2fcd8ce9f53540 - Debug and test logs removed from utils log
de9b59fa26fab25ba225b9a8288b8093d91ebf7a - Error reporting related logic changed
13e27ebae5e6ce86712499bd42197b26a98a14db - Added debug log
1984560bcbd1f9d009d4824436f2877842b768b3 - Log utils error reporting logic changed
043e533beaf586c4cc37b5f5c3602ebd5802dfea - Test script changed to avoid caching
6e4808a1daf975178bbd30f9f43e7de74d04c745 - Error reporting client and call logic added to log util
3965748bc1f9309317760df6007f9700998b5e63 - Project id value added to environment config
c0beda21860fbfb99f219e362a0fe810436b6a54 - Added more logs, fixed some memory pointer errores and test utils
a80d4e037877c309331dad2bb6fe97120a0d0c57 - Major code refactor, including logs, types usages and filenames
bba5f7daa490ec09aa337c975a29cbdd51f46774 - Fixed sonar properties and major file lint revision
dd39053398ee44c3c9093b0d4fc5a7d523fe766c - Sonar properties file added and other code smells fixed
c9e568d68188a186ea4a7bd184aebf4772b4af97 - Fixed some sonarqube code smells
067b23ab8e9bd2524cf9ff2d12105b386c339310 - cloud build file added for sonarqube scanning
4109e663e8e5d69c22d2a669d2272bfcdc5c119b - README, config/env.json and tests updated
674f641c6d8cbf6c65f79b4fa90689744513c4a5 - gitignore updated
cd9657ff1a3f50074291c997cac7845b56ec9b9f - Initial API test setup and deposit validation endpoint tests added
75b1d6716a05f672c292e545ea00f50406d2fcbf - Major services refactor in order to encapsulate logic and to be tested more easily
19ad39df8357b8d062922b76ae0d55770971fc5b - Accounts service refactored to prioritize gRPC protocol first, and in case of failure use http
ac7d7479329224dba6992829e94d8840b3c78db0 - Accounts engine service refactored
f313bb905c412493501224a45a16c3832ff07f5c - Http delete, get, options, post and put utils functions added
2febcb375bdda0a6fdb871f4037502eb6495285a - Debug util log fixed
bc6b749d282310b56ae87b0500e5d205cc21dacd - Debug logs removed
8644ef16a27ab4b6b0b31963ab8aa8db328dd0eb - Debug log for checking token expiry added again for testing
90b151e8b5393a852e7eb1edf89b2e25f571370f - Errormf utils function fixed
b0763487d46b4774231cb9ca55c53200c07f4c37 - Code reordered and debug logs removed
2d97b63ac1328ef573d8e2e54a0a6efd2324ba9e - Debug log added to check token expiry
ab618c930d9358a30ea5aa17d1ddeb01c34165ec - README updated
a02da722213ca7b9d06fef2c2b569adb57ecdb18 - k8s deploy yaml updated
b720ff363e0a6b02e19aabc9bcedd5e44a784039 - Fixed typo error on k8s deploy yaml
bdd94b404122357a68333b044bb456562840bd42 - Changed grpc utils to use env variable
a9169c7d4e161f44222f682acd2e50ed17560862 - gRPC client option credentials loading changed
7ca10f4928eb8e76b0574be50f45defef8d59bd1 - Debug log updated
000a23e24f65901608c538eb4818d0636e13b3d0 - k8s deploy yaml and grpc env variables renamed
7e2127ba0b09a83a54a40090a6b33e1533059dd6 - Debug log added to print grpc credential
9f49b254f20d3912f851352b7d378e70db6fb1f7 - Added missing return on deposit validation error handling
33514834afe25feae7792c4828a95ceb1fc2743f - k8s deploy yaml updated
70d6d682f064852ed98c881bf55ae9ce5dbbeefe - Deposit validations now uses accounts service instead of accounts engine service when getting account info
bbb1aca47fff6fa3f785c6efdd39c6b633f7d3a8 - Accounts service added
70be091844c265735b60ba2c5b559f9a8ff5be28 - gRPC utils added
8782390e1c400528d26edc4473c5ab46a6dc5b50 - Date util ParseStringISODate added
d4091da626c274ce7f00740298d1e101211205bb - Account model updated
41eaa7aa4593f4855d8b92188be97cf4bb860ede - Utils functions reordered and MapStringToMapInterface func added
a182998e72b935e1adadbd5deb2fa3ac0d1726f0 - gRPC accounts engine settings added
548290b6de55617b0f029d45a7c2e11ea8c025cd - fixed PR comments.
cfe4cc0b4c32ee6aee86ac5efd8735128ec3cc7b - Changed the name of KnownError struct to HandledError.
c5cc1c32b41f03e2769f5e63570f77d0ebcb5de8 - Changed the label of the handled errors in received transfers.
020d9217286cc87901a4cfab334fb4788cf066a0 - Changed deposit validations logic in order to check for transfer payments requests before finding a user.
cf631bac0e2b464b20203c10ade8202c1e61f535 - added go.mod and go.sum to version control.
26a9691cb09730d7c3abdbcf3e3df0f6abc72bca - Changed the status of transfer payments in deposit validations in order to avoid being used for another transfer.
f9e97cae8ba8d3e901c4a31b1966a5de509be6a5 - added missing field id to env.json
2af6893b9eaf371843854331f94c454c784e9cbd - Fixed indentation
713566b6e40af7699c9e1a30b571c0b50a16ebe9 - Accounts rules settings changed to user and commerce rules, and firebase service check refactored
177a3d56ac2d1c7b38f00e736b9d29ebaf1d6b31 - Postman collection updated with local, dev and staging environment and endpoints with payment transfer requests

## 1.2.0 - 2020-10-22

### Description

- Deposit validations supports transfer payment requests
- Received and reversed transfers endpoint changed in order to use new transfers collection
- Added more logs to track important events and changed some to include metadata as a searchable property on stackdriver

### Changelog

656b206b5473bf573daefe7f2908c86c3cf2456c - Version bump -version 1.2.0
d0b2f9382dc2e14bd7d84e458743d60abd051b9e - Small fixes on logging
17012d50ffef80d1c330d4174da0451b08d753da - Received and reversed transfers endpoint refactored to use new transfers collection
e91681d29a782b8824be02eed9a3764ac98bf2ae - Removed order by and sort is now applied from code logic
f5c0462abeab88ac4d7547ef76f1d5452e254ba1 - Transfer payment requests get documents list call updated
8c3829feb7b22dad1f87b42048c1d3baa9367596 - Database service query options implemented to apply limit, order by, parent id and/or transaction to getSnapshots and getDocumentsList
a19106f79410b7a5bb4b928166bbbfaaadc736bd - Logging improved on endpoints to make filtering easier
37417834f19c88b52abad4d5acaf66662bc5dc13 - Step error codes updated
ca091b1ac027d610e0e373a1f1cd4916964328e0 - Fixed indentation
dec6d042b9957f5e2cdb74f2df0ec7b2f927ae93 - Deposit validations endpoint changed to handle deposits of transfer payment requests
e58a2275a02e8a67541392cb87390dc3405e2404 - Small log refactored applied
5a0f4fbad760f5d9b229404f3800b36a927b228c - Database, users, transfer payment request service and model added
3b9a9288e43bbcad650e1449da2eabfac2cfea5b - Account engine service updated

## 1.1.4 - 2020-10-01

### Description

- GKE 1.16+ upgrade related api versions upgraded

### Changelog

edb0aed234fa73918ac27b109b8dd467724ee8c2 - Version bump —version 1.1.4
f5c81d5fe90989ac91065966e3cedecf6ae7003a - removed restrictions for banco ripley
74ebe5292909b71551e0ee25ba4445b369521d14 - blocked banco ripley again
1f9e6c3b70ffea3b2211952bf9269ae2ed666919 - removed code that was blocking cash in from banco ripley
0fb47f45b5785337cdb4f4bc1ec2e8ba1a3256e0 - blocked cash ins from bank 0053
cfc12e6ac04d93051bb7ce73aeb88d4478290169 - GKE 1.16+ upgrade related api versions upgraded

## 1.1.3 - 2020-06-08

### Description

- Account max balance and category related transactional limit check now correctly handles cases when either comes with null value

### Changelog

e8f6d105c5a577e7ba9f06581b0224de52527f72 - Version bump —version 1.1.3
6f20e437087318c873c1fece17da88b00b699102 - Indentation fixed
ada30e7f18dba1e642f103d96d39ea3a6bd436ef - Small log refactor and max balance check now handles accounts without value
0929440a579e3e4cf80c94313566bfe5e0b65f06 - GetAccount now initializes max balance and category values to check if they come with value or not
e24684ab07e3986a123d72b1b0948371b9beeb93 - checkMaxDepositAmountInTimePeriod function now correctly handles accounts without category
636c0d5db1bf246d17282aaf5f58bd9163847289 - Category property added to Account model

## 1.1.2 - 2020-04-30

### Description

- Requester ip address added to logging metadata

### Changelog

a2a4e5d6ed771a4c984a066320621875afd6002f - Version bump -version 1.1.2
50db4f08861d5b363e72c54b264f755e31440a39 - GetIP return string format applied
9d1257b60a08b5dc116cc8387d101518d4b8e4aa - Test logging both request related ip
6f63d4ef2c3942be6437e005e6504df58ef9b260 - Request ip added to logging metadata

## 1.1.1 - 2020-04-29

### Description

- Added camelcase routing for endpoints 
- Hotfix cash in service 

### Changelog

1b3f153992e43040b304a6ad124680809f882716 - Version bump -version 1.1.1
cc1dc9796183c38e8b130ef4e640cbb281f71fec - Added camelcase routing for endpoints
70ded27422892cf32b4c7353d786942bfdca4602 - CHANGELOG commit id fixed
7c81afc15460a7ee179f77c43c766d6f15ed41b2 - Cash in service and ingress config fixed to apply cloud armor security policies on port 8

## 1.1.0 - 2020-04-28

### Description

- Logging refactored
- Push message changed when deposit validation checks over max balance limit
- Deposit validations now checks if origin rut is on banned origin rut list
- K8s configuration updated for new cluster configuration

### Changelog

ec47943f1c9af3004fcea2499da1ef3f3ac4b130 - Version bump -version 1.1.0
bd009bf6ef3ffade74ab3d8e9f1c7ddb20f766e0 - gitignore updated
40c5346f8da29716800458492769a601a48d2387 - project.banned_origin_rut property added to env file
273923eca17f3ae93393cb62daee5ab653791ab0 - Env variable file transactional limits configuration values updated
04025c4fbacac3ed2ed09770288c39f5e90ff510 - README updated
0d486593dfc5218e219e25274e24c980ba9f674c - Project environment variable and static project names removed and k8s configuration files added and updated for new cluster configuration
e1c5ab2f944f5ceefbd5bed0bc9984a13f173e49 - Dockerfile changed, go version upgraded
376b8cec2ed8f5d3212646d02f0ea2123c47f059 - Deposit validations now checks if origin rut is on banned origin rut list
7f14e6c239215b60d2a8aa2f92aa1e748eef31c8 - Safe ips configuration repurposed and renamed for banned origin ruts
983e565a82a2bef16414f2fe5a2366b9eb562963 - ArrayIncludes function added to utils
aec8f6937931da554cb48c4f7ed33bffa211381e - gitignore and related config files renamed to a more generic ones
397399e3a388e5c25ad24dd9b2247c155eed7c8c - Push message changed when deposit validation checks over max balance limit
b44acb7ee0c0e1e2e6d481bad9e8ab58ae1efb11 - GetLoggerID util added
9ab92982f02b0375d6b62e53a8e85baf654ae5ea - GetCurrentStepCode and GetCurrentStepLabel util added
ec761de4e4411f12caf3679b2220bc5231de9d91 - Rest endpoints logging refactored and transaction implementation added when needed
aff32baf70a6f6d1451d0d665207f1ac13402229 - Http respond, truncate file and functio name utils moved, log utils refactored and http utils added
cb2e9426cee3816994bd9c4bc8f3eeffb37b6876 - Rest middleware logging refactored
5b5747248e317fae61f4afb9d430ed5f3e1e6206 - Current step and start time context keys constants added and requestId renamed

## 1.0.2 - 2020-01-27

### Description

- Fixed bank transfer transactional limit rule check related to date query

### Changelog

07b1d9af7535038738bc1750380d81dfb83dd091 - Version bump —version 1.0.2
0d689cef44d8bde67d22c5f5d9839793bdf156a1 - checkMaxDepositAmountInTimePeriod date query fixed


## 1.0.1 - 2019-12-26

### Description

- Fixed account engine refresh token logic
- Added proxy files

### Changelog

e921f6333811f20d3c3d6541510f69f85818aa16 - Version bump —version 1.0.1
3557a136cccd01b1942522b1619c69a49ad8e8df - Fixed refreshToken next refresh timestamp calculation logic
cbbdb83cbf9c4bc0f1439ab981ecc5e06cb1654a - added proxy files


## 1.0.0 - 2019-12-24

### Description

- REST API Endponts:
    - `POST /access-tokens`: Generates an access token needed to consume the other endpoints
    - `POST /deposit-validations`: Validates if the destination account can receive the incoming deposit which has not been processed yet
    - `POST /received-transfers`: Creates a receivedTransfer to be processed
    - `POST /reversed-transfers`: Begins a reverse procedure of a received transfer
- LOCAL, STAGING and PROD environments
- Base project

### Changelog

76209e2d7477f03e86db81b798e0a1dcabb7fbbd - Version bump —version 1.0.0
a55f55b0819264e8c99f5b4fa27426ca00004ab6 - Remove unused code and added comment to explain account engine service RefreshWindow usage
1b35ea89dacb68430b7a4b62c0af06373420ebd1 - Changed context usage for rest endpoints
bbc78c766fa3ab631c6721e06a1c2276c8f19bf0 - Format and code optimized
60212c7b8ae889e7ce50d542a4f9f158d182c0b5 - Ingress deployment files updated
ef08d76269e4671de341dd8b2c614b0099101e89 - docker-compose-yml deleted and ingress deploy updated
07bd78f31b7b4d06621429a64cdd7b7e66efbc50 - Reversed transfers success response message changed
600e44336533dc0385995ee32fa0b390870fa97b - Internal server error messages changed to be more generic and account engine refresh token logic fixed
08fbab71e26c5d337b92ca40151bc369505ab55e - README, gitignore and cloud build config updated
68f881d51811a2dac7732b30be91aa42b1dafa3c - Small tweaks in settings, endpoints and service snapshots validations
77affc739f07d3742e447822f61556cf742f8db6 - k8s deployment files added
55dba67489c5166ff3dc583cf79e5efba1938728 - Postman endpoints collection added
ed55b596452aec8306e35e3000b11e5020f95a27 - Small tweaks in ENDPOINTS_README
6bb1f9695d72585f806cc3488a303057e25be6d2 - README update and ENDPOINTS_README added
3a359ec4dc0d294c873a7871c5f02e1225282869 - Small fixes to access tokens, reversed transfers endpoints and firebase service
6cbc44c02abc27c682b25a763535a9740b2db4dc - Added initial docker config, readme and changelog
52fb1b8a3dc13d2f3b4b8f965c30c9e9e0d08602 - Added constants package
c3ca41b0834aa410b9100e9714b6356b18417423 - Reversed transfers endpoint implemented
4c093a4d0c05278c0b47e7dc3e8713e50cda9b08 - Changed logs usage for endpoints
af0600dbb5fc13504237580782c32916b603cbe9 - Services refactored and logs usages changed
3de61900d986b2dcf6dfcc7867ace1dd151edcea - Models package changed
a02b58a51802bd6b68842698b1827be61c1a9d0d - Settings package refactored
40bce93613d601a7f818c61201c3a8a0a93f27a5 - Utils package refactored
1452623f4bcbdea932cb1b659e40f61a9b85630e - Received transfers endpoint implemented
cd52f3a14c1be28714ffa7406d9b1c1bee942169 - General refactored with models and log changes
b9f2ea1bfaefda365f7a3f5688b84a872c44d234 - Firebase auth token print and decode logic added to firebase service
f91b386f8573e569546f5c477449a340d222a3b9 - SafeIPs configuration added to settings
e1f5ba4292efd90031883d0c467caae621f58eba - GetClientClaims function moved to claims file
b170de721d70adcbba67959ddab9b828dd19b86a - Utils packages refactored
7136c8210f31244f8b11e60b6b0a7c6e8d73a461 - Models package refactored
287d12bce31838ac1a65a350c7595ab7228e6a46 - Deposit validations endpoint implemented
2717b11d1975dcaffa8938850510369c09d9ac5b - Access tokens endpoint updated
ed99dc03bd80775650c6a7e8d4b284acbee2dd87 - Account Engine service added
4862877be4ff16b269346ec218538f289784be87 - Firebase service refactored
261cca4dea155936c07db728f81001e329bc00f0 - Identity toolkit service refactored
edc843f1cd0b0a81057b16c4f625d7a584fe7140 - Notification service package added
fc44ae531a7c8196f0d2c4f1783b984d667bed12 - AccountEngineConfig and AccountCategoryConfig settings added
8a4a267d11b32137b7e095f6833dc047b5b2ad02 - Account and User models added
1d9c858a193b52b659e885cfb9067d4e749fcd66 - Format utils package added
7046803381d847a8ed39662773adc601ba626ed3 - Transactional rules util package added
9c1cf972f929ed67cb2dbbcfa12cf8c03a3e0bd0 - Main utils package updated
eaac0c697189fdd42e0c898e9fbe8901323d6c2a - test flow removed and Log util changed
a46e1eaeef9575961d19516aff9328e2cf031c9e - Main application and gitignore added
676eaa792398b98ecbcb453c41027e1dfd365b4a - Main endpoints handlers added
40dd7f3a35111e87353e5d0b43fc4073edb28f27 - Relyingparty package added
a80d319bf13e687a669d04f8368395266c791b4e - Google Identity Toolkit package added
c9a05f1c5edb181c7c2b6fdfdeb812efbd0caa0f - Firebase package added
ada85b77f374f5ce1d4d1868b5fe77a78dde3bd7 - Access tokens endpoint added
9968a9d333b2eb52c97ebc60fe94067eaa66c49e - Home index endpoint added
15dfebbc30d890030f0f34883cbd58b1fdecb872 - Deposit validations base endpoint added
0374cc20b5d5af870bb8f54a3d0d5b745cc17614 - Received transfers base endpoint added
90e8e3d59672c1e9429ae085c979d7844f162cde - Reversed transfers base endpoint added
dccd217aa7d8e261ed0e7c8b7acf5296ba04dcaa - Added models package with access-token related entities
8a7888994740b8a6cfbf2fbc22c828e85a339cf7 - Added settings package
89db6baa74d00906895905d12f8fc62bf12b5e40 - Added utils package
977a5524ad6a636d0b03d665148deeb97c80ddd5 - Initial commit